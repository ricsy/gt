//go:build ignore

package main

import (
	"encoding/json"
	"flag"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"unicode"
)

type metadata struct {
	Endpoints      []endpoint               `json:"endpoints"`
	RequestStructs map[string]requestStruct `json:"requestStructs"`
	DirectCalls    []directCall             `json:"directCalls"`
}

type config struct {
	Go goConfig `json:"go"`
}

type goConfig struct {
	EndpointFile          string   `json:"endpointFile"`
	EndpointGroupType     string   `json:"endpointGroupType"`
	ResponseGlobs         []string `json:"responseGlobs"`
	APIGlobs              []string `json:"apiGlobs"`
	RequestStructSuffixes []string `json:"requestStructSuffixes"`
	ParameterTags         []string `json:"parameterTags"`
	ParameterNameCase     string   `json:"parameterNameCase"`
	DirectCallMethod      string   `json:"directCallMethod"`
	IgnoreDirectCallFiles []string `json:"ignoreDirectCallFiles"`
}

type endpoint struct {
	Group  string `json:"group"`
	Action string `json:"action"`
	Method string `json:"method"`
	Path   string `json:"path"`
	File   string `json:"file"`
	Line   int    `json:"line"`
}

type requestStruct struct {
	File   string  `json:"file"`
	Line   int     `json:"line"`
	Fields []field `json:"fields"`
}

type field struct {
	Name      string   `json:"name"`
	ParamName string   `json:"paramName"`
	Tags      []string `json:"tags,omitempty"`
}

type directCall struct {
	File   string `json:"file"`
	Line   int    `json:"line"`
	Method string `json:"method,omitempty"`
	Path   string `json:"path,omitempty"`
}

func main() {
	root := flag.String("root", ".", "repository root")
	configPath := flag.String("config", "scripts/api-doc-diff.config.json", "comparison config path")
	flag.Parse()

	cfg, err := readConfig(*configPath)
	if err != nil {
		exitErr(err)
	}

	fset := token.NewFileSet()
	result := metadata{
		RequestStructs: map[string]requestStruct{},
	}

	endpoints, err := extractEndpoints(fset, *root, cfg)
	if err != nil {
		exitErr(err)
	}
	result.Endpoints = endpoints

	requestStructs, err := extractRequestStructs(fset, *root, cfg)
	if err != nil {
		exitErr(err)
	}
	result.RequestStructs = requestStructs

	directCalls, err := extractDirectCalls(fset, *root, cfg)
	if err != nil {
		exitErr(err)
	}
	result.DirectCalls = directCalls

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(result); err != nil {
		exitErr(err)
	}
}

func readConfig(path string) (config, error) {
	var cfg config
	file, err := os.Open(path)
	if err != nil {
		return cfg, err
	}
	defer file.Close()
	err = json.NewDecoder(file).Decode(&cfg)
	return cfg, err
}

func extractEndpoints(fset *token.FileSet, root string, cfg config) ([]endpoint, error) {
	path := filepath.Join(root, filepath.FromSlash(cfg.Go.EndpointFile))
	file, err := parser.ParseFile(fset, path, nil, 0)
	if err != nil {
		return nil, err
	}

	var endpoints []endpoint
	for _, decl := range file.Decls {
		gen, ok := decl.(*ast.GenDecl)
		if !ok || gen.Tok != token.VAR {
			continue
		}
		for _, spec := range gen.Specs {
			valueSpec, ok := spec.(*ast.ValueSpec)
			if !ok {
				continue
			}
			for i, value := range valueSpec.Values {
				if i >= len(valueSpec.Names) {
					continue
				}
				lit, ok := value.(*ast.CompositeLit)
				if !ok || !isEndpointGroup(lit.Type, cfg.Go.EndpointGroupType) {
					continue
				}
				group := valueSpec.Names[i].Name
				for _, elt := range lit.Elts {
					kv, ok := elt.(*ast.KeyValueExpr)
					if !ok {
						continue
					}
					actionIdent, ok := kv.Key.(*ast.Ident)
					if !ok {
						continue
					}
					method, path, ok := endpointLiteral(kv.Value)
					if !ok {
						continue
					}
					pos := fset.Position(kv.Pos())
					endpoints = append(endpoints, endpoint{
						Group:  group,
						Action: actionIdent.Name,
						Method: method,
						Path:   path,
						File:   rel(root, pos.Filename),
						Line:   pos.Line,
					})
				}
			}
		}
	}
	return endpoints, nil
}

func extractRequestStructs(fset *token.FileSet, root string, cfg config) (map[string]requestStruct, error) {
	files, err := expandGlobs(root, cfg.Go.ResponseGlobs)
	if err != nil {
		return nil, err
	}

	result := map[string]requestStruct{}
	for _, path := range files {
		file, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
		if err != nil {
			return nil, err
		}
		for _, decl := range file.Decls {
			gen, ok := decl.(*ast.GenDecl)
			if !ok || gen.Tok != token.TYPE {
				continue
			}
			for _, spec := range gen.Specs {
				typeSpec, ok := spec.(*ast.TypeSpec)
				if !ok || !isRequestType(typeSpec.Name.Name, cfg.Go.RequestStructSuffixes) {
					continue
				}
				structType, ok := typeSpec.Type.(*ast.StructType)
				if !ok {
					continue
				}
				rs := requestStruct{
					File: rel(root, path),
					Line: fset.Position(typeSpec.Pos()).Line,
				}
				for _, astField := range structType.Fields.List {
					for _, name := range astField.Names {
						paramName, tags := parameterName(name.Name, astField.Tag, cfg.Go.ParameterTags, cfg.Go.ParameterNameCase)
						rs.Fields = append(rs.Fields, field{
							Name:      name.Name,
							ParamName: paramName,
							Tags:      tags,
						})
					}
				}
				result[typeSpec.Name.Name] = rs
			}
		}
	}

	return result, nil
}

func extractDirectCalls(fset *token.FileSet, root string, cfg config) ([]directCall, error) {
	files, err := expandGlobs(root, cfg.Go.APIGlobs)
	if err != nil {
		return nil, err
	}

	var calls []directCall
	ignored := ignoredFiles(cfg.Go.IgnoreDirectCallFiles)
	for _, path := range files {
		relative := rel(root, path)
		if strings.HasSuffix(path, "_test.go") || ignored[relative] {
			continue
		}
		file, err := parser.ParseFile(fset, path, nil, 0)
		if err != nil {
			return nil, err
		}
		ast.Inspect(file, func(node ast.Node) bool {
			call, ok := node.(*ast.CallExpr)
			if !ok {
				return true
			}
			selector, ok := call.Fun.(*ast.SelectorExpr)
			if !ok || selector.Sel.Name != cfg.Go.DirectCallMethod {
				return true
			}
			pos := fset.Position(call.Pos())
			dc := directCall{File: rel(root, pos.Filename), Line: pos.Line}
			if len(call.Args) > 0 {
				dc.Method = exprLabel(call.Args[0])
			}
			if len(call.Args) > 1 {
				dc.Path = exprLabel(call.Args[1])
			}
			calls = append(calls, dc)
			return true
		})
	}

	return calls, nil
}

func expandGlobs(root string, patterns []string) ([]string, error) {
	seen := map[string]bool{}
	var files []string
	for _, pattern := range patterns {
		matches, err := filepath.Glob(filepath.Join(root, filepath.FromSlash(pattern)))
		if err != nil {
			return nil, err
		}
		for _, match := range matches {
			if seen[match] {
				continue
			}
			seen[match] = true
			files = append(files, match)
		}
	}
	return files, nil
}

func ignoredFiles(paths []string) map[string]bool {
	result := map[string]bool{}
	for _, path := range paths {
		result[filepath.ToSlash(path)] = true
	}
	return result
}

func isEndpointGroup(expr ast.Expr, typeName string) bool {
	ident, ok := expr.(*ast.Ident)
	return ok && ident.Name == typeName
}

func endpointLiteral(expr ast.Expr) (string, string, bool) {
	lit, ok := expr.(*ast.CompositeLit)
	if !ok {
		return "", "", false
	}
	var method, path string
	if len(lit.Elts) >= 2 {
		method = exprLabel(lit.Elts[0])
		path = exprLabel(lit.Elts[1])
	}
	for _, elt := range lit.Elts {
		kv, ok := elt.(*ast.KeyValueExpr)
		if !ok {
			continue
		}
		key, ok := kv.Key.(*ast.Ident)
		if !ok {
			continue
		}
		switch key.Name {
		case "Method":
			method = exprLabel(kv.Value)
		case "Path":
			path = exprLabel(kv.Value)
		}
	}
	return method, path, method != "" && path != ""
}

func exprLabel(expr ast.Expr) string {
	switch value := expr.(type) {
	case *ast.BasicLit:
		if value.Kind == token.STRING {
			unquoted, err := strconv.Unquote(value.Value)
			if err == nil {
				return unquoted
			}
		}
		return value.Value
	case *ast.Ident:
		return value.Name
	case *ast.SelectorExpr:
		return exprLabel(value.X) + "." + value.Sel.Name
	case *ast.BinaryExpr:
		return exprLabel(value.X) + value.Op.String() + exprLabel(value.Y)
	default:
		return ""
	}
}

func isRequestType(name string, suffixes []string) bool {
	for _, suffix := range suffixes {
		if strings.HasSuffix(name, suffix) {
			return true
		}
	}
	return false
}

func parameterName(fieldName string, tag *ast.BasicLit, tagKeys []string, nameCase string) (string, []string) {
	if tag == nil {
		return convertName(fieldName, nameCase), nil
	}
	raw, err := strconv.Unquote(tag.Value)
	if err != nil {
		return convertName(fieldName, nameCase), nil
	}
	var tags []string
	for _, key := range tagKeys {
		value := lookupStructTag(raw, key)
		if value == "" || value == "-" {
			continue
		}
		name := strings.Split(value, ",")[0]
		if name == "" || name == "-" {
			continue
		}
		tags = append(tags, key+":"+name)
		return name, tags
	}
	return convertName(fieldName, nameCase), tags
}

func lookupStructTag(tag, key string) string {
	prefix := key + ":"
	for tag != "" {
		tag = strings.TrimLeft(tag, " ")
		if tag == "" {
			break
		}
		i := strings.Index(tag, ":")
		if i <= 0 {
			break
		}
		name := tag[:i]
		tag = tag[i+1:]
		if !strings.HasPrefix(tag, "\"") {
			break
		}
		value, rest, ok := readQuoted(tag)
		if !ok {
			break
		}
		if name == key || name+":" == prefix {
			return value
		}
		tag = rest
	}
	return ""
}

func readQuoted(input string) (string, string, bool) {
	for i := 1; i < len(input); i++ {
		if input[i] == '"' && input[i-1] != '\\' {
			value, err := strconv.Unquote(input[:i+1])
			if err != nil {
				return "", "", false
			}
			return value, input[i+1:], true
		}
	}
	return "", "", false
}

func lowerCamel(value string) string {
	if value == "" {
		return ""
	}
	runes := []rune(value)
	for i := 0; i < len(runes); i++ {
		if i+1 < len(runes) && unicode.IsUpper(runes[i+1]) && i > 0 {
			break
		}
		runes[i] = unicode.ToLower(runes[i])
		if i+1 == len(runes) || (i+1 < len(runes) && unicode.IsLower(runes[i+1])) {
			break
		}
	}
	return string(runes)
}

func convertName(value, nameCase string) string {
	switch nameCase {
	case "snake_case":
		return snakeCase(value)
	case "lowerCamel", "lower_camel":
		return lowerCamel(value)
	default:
		return lowerCamel(value)
	}
}

func snakeCase(value string) string {
	if value == "" {
		return ""
	}
	var result []rune
	runes := []rune(value)
	for i, r := range runes {
		if unicode.IsUpper(r) {
			if i > 0 {
				prev := runes[i-1]
				nextLower := i+1 < len(runes) && unicode.IsLower(runes[i+1])
				if prev != '_' && (unicode.IsLower(prev) || unicode.IsDigit(prev) || nextLower) {
					result = append(result, '_')
				}
			}
			result = append(result, unicode.ToLower(r))
			continue
		}
		result = append(result, r)
	}
	return string(result)
}

func rel(root, path string) string {
	relative, err := filepath.Rel(root, path)
	if err != nil {
		return path
	}
	return filepath.ToSlash(relative)
}

func exitErr(err error) {
	_ = json.NewEncoder(os.Stderr).Encode(map[string]string{"error": err.Error()})
	os.Exit(1)
}
