import js from '@eslint/js'
import tsPlugin from '@typescript-eslint/eslint-plugin'
import tsParser from '@typescript-eslint/parser'
import pluginVue from 'eslint-plugin-vue'
import configPrettier from 'eslint-config-prettier'

const tsRules = {
  ...tsPlugin.configs.recommended.rules,
  '@typescript-eslint/no-explicit-any': 'warn',
  '@typescript-eslint/no-unused-vars': ['warn', { argsIgnorePattern: '^_' }],
  // empty catch blocks are common in Wails IPC / event handlers
  'no-empty': ['error', { allowEmptyCatch: true }],
}

export default [
  {
    // declaration files are typically generated or framework boilerplate
    ignores: ['src/**/*.d.ts'],
  },

  js.configs.recommended,

  // TypeScript source files
  {
    files: ['src/**/*.ts'],
    languageOptions: {
      parser: tsParser,
      parserOptions: { ecmaVersion: 'latest', sourceType: 'module' },
    },
    plugins: { '@typescript-eslint': tsPlugin },
    rules: tsRules,
  },

  // Vue single-file components (vue-eslint-parser delegates <script lang="ts"> to tsParser)
  ...pluginVue.configs['flat/recommended'],
  {
    files: ['src/**/*.vue'],
    languageOptions: {
      parserOptions: {
        parser: tsParser,
        ecmaVersion: 'latest',
        sourceType: 'module',
      },
    },
    plugins: { '@typescript-eslint': tsPlugin },
    rules: {
      ...tsRules,
      'vue/multi-word-component-names': 'off',
    },
  },

  // Disable formatting rules that conflict with Prettier
  configPrettier,
]
