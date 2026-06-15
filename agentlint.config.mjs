import {
  CacheStrategy,
  defineConfig,
  EngineeringRuleId,
  recommendedPreset,
  RuleSeverity,
} from 'agentlint'

const sharedStrictRules = {
  [EngineeringRuleId.NoDuplicateCode]: RuleSeverity.Error,
  [EngineeringRuleId.NoDuplicateFunction]: RuleSeverity.Error,
}

const goSourceFiles = [
  'main.go',
  'internal/**/*.go',
  'pkg/**/*.go',
  'plugins/**/*.go',
]

const goTestFiles = [
  'internal/**/*_test.go',
  'pkg/**/*_test.go',
  'plugins/**/*_test.go',
]

const goLanguagePlugins = [
  '@agentlint/language-go',
]

export default defineConfig([
  {
    ignore: {
      inheritGitignore: true,
    },
    cache: {
      enabled: true,
      strategy: CacheStrategy.Content,
    },
  },
  {
    extends: [recommendedPreset],
    files: goSourceFiles,
    ignores: [
      '**/*_test.go',
    ],
    languagePlugins: goLanguagePlugins,
    rules: {
      ...sharedStrictRules,
      [EngineeringRuleId.NoHardcoding]: RuleSeverity.Error,
    },
  },
  {
    extends: [recommendedPreset],
    files: goTestFiles,
    languagePlugins: goLanguagePlugins,
    rules: {
      ...sharedStrictRules,
      [EngineeringRuleId.NoHardcoding]: RuleSeverity.Off,
    },
  },
])
