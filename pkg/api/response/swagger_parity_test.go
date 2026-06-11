package response

import (
	"reflect"
	"testing"
)

func TestProjectMatchesReviewedSwaggerFields(t *testing.T) {
	assertHasJSONField[Project](t, "path")
	assertHasJSONField[Project](t, "public")
	assertHasJSONField[Project](t, "internal")
	assertHasJSONField[Project](t, "permission")
	assertHasJSONField[Project](t, "project_labels")
	assertMissingJSONField[Project](t, "clone_url")
}

func TestPullRequestMatchesReviewedSwaggerFields(t *testing.T) {
	assertHasJSONField[PullRequest](t, "draft")
	assertHasJSONField[PullRequest](t, "_links")
	assertHasJSONField[PullRequest](t, "ref_pull_requests")
	assertMissingJSONField[PullRequest](t, "comments")
	assertMissingJSONField[PullRequest](t, "files_changed")
	assertMissingJSONField[PullRequest](t, "merged")
}

func TestReleaseMatchesReviewedSwaggerFields(t *testing.T) {
	assertHasJSONField[Release](t, "prerelease")
	assertMissingJSONField[Release](t, "published_at")
	assertMissingJSONField[Release](t, "html_url")
}

func TestSSHKeyBasicMatchesReviewedSwaggerFields(t *testing.T) {
	assertHasJSONField[SSHKeyBasic](t, "key")
	assertMissingJSONField[SSHKeyBasic](t, "title")
	assertMissingJSONField[SSHKeyBasic](t, "url")
}

func TestCreateCheckRunOptionsMatchReviewedSwaggerFields(t *testing.T) {
	assertHasJSONField[CreateCheckRunOptions](t, "pull_request_id")
	assertHasJSONField[CreateCheckRunOptions](t, "started_at")
	assertHasJSONField[CreateCheckRunOptions](t, "output[text]")
	assertHasJSONField[CreateCheckRunOptions](t, "output[annotations][path]")
	assertHasJSONField[CreateCheckRunOptions](t, "output[images][image_url]")
	assertHasJSONField[CreateCheckRunOptions](t, "actions[identifier]")
}

func TestUpdateCheckRunOptionsMatchReviewedSwaggerFields(t *testing.T) {
	assertHasJSONField[UpdateCheckRunOptions](t, "pull_request_id")
	assertHasJSONField[UpdateCheckRunOptions](t, "name")
	assertHasJSONField[UpdateCheckRunOptions](t, "output[text]")
	assertHasJSONField[UpdateCheckRunOptions](t, "output[annotations][annotation_level]")
	assertHasJSONField[UpdateCheckRunOptions](t, "actions[label]")
}

func assertHasJSONField[T any](t *testing.T, jsonName string) {
	t.Helper()

	if _, ok := jsonFieldMap[T]()[jsonName]; !ok {
		t.Fatalf("expected json field %q to exist", jsonName)
	}
}

func assertMissingJSONField[T any](t *testing.T, jsonName string) {
	t.Helper()

	if _, ok := jsonFieldMap[T]()[jsonName]; ok {
		t.Fatalf("expected json field %q to be absent", jsonName)
	}
}

func jsonFieldMap[T any]() map[string]reflect.StructField {
	typ := reflect.TypeOf((*T)(nil)).Elem()
	fields := make(map[string]reflect.StructField, typ.NumField())
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		tag := field.Tag.Get("json")
		if tag == "" || tag == "-" {
			continue
		}
		name := tag
		for j := 0; j < len(name); j++ {
			if name[j] == ',' {
				name = name[:j]
				break
			}
		}
		fields[name] = field
	}
	return fields
}
