import eslint from '@eslint/js';
import typescript from '@typescript-eslint/eslint-plugin';
import tsParser from '@typescript-eslint/parser';
import svelte from 'eslint-plugin-svelte';
import astro from 'eslint-plugin-astro';
import { parseForESLint as astroParseForESLint } from 'astro-eslint-parser';
import globals from 'globals';

export default [
  {
    ignores: ['dist/**', '.astro/**', 'node_modules/**'],
  },
  eslint.configs.recommended,
  ...astro.configs['flat/recommended'],
  ...svelte.configs['flat/recommended'],
  {
    // TypeScript files configuration
    files: ['**/*.{ts,tsx}'],
    plugins: { '@typescript-eslint': typescript },
    languageOptions: {
      parser: tsParser,
      globals: { ...globals.browser },
      parserOptions: {
        project: './tsconfig.json',
        ecmaVersion: 2022,
        sourceType: 'module',
      },
    },
    rules: {
      ...typescript.configs.recommended.rules,
      // TypeScript's own compiler checks this; ESLint flags browser DOM globals.
      'no-undef': 'off',
    },
  },
  {
    // Astro files: wire up the parser with TypeScript support
    files: ['**/*.astro'],
    languageOptions: {
      parser: { parseForESLint: astroParseForESLint },
      parserOptions: {
        parser: tsParser,
        extraFileExtensions: ['.astro'],
        sourceType: 'module',
      },
    },
  },
  {
    // Svelte files: browser globals + TypeScript-aware rules
    files: ['**/*.svelte'],
    plugins: { '@typescript-eslint': typescript },
    languageOptions: {
      globals: { ...globals.browser },
      parserOptions: {
        parser: tsParser,
        typescript: true,
      },
    },
    rules: {
      ...typescript.configs.recommended.rules,
      'no-undef': 'off',
      'no-unused-vars': 'off',
    },
  },
];
