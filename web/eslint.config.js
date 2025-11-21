import eslint from '@eslint/js';
import typescript from '@typescript-eslint/eslint-plugin';
import tsParser from '@typescript-eslint/parser';
import svelteParser from 'svelte-eslint-parser';
import astroParser from 'astro-eslint-parser';
import svelte from 'eslint-plugin-svelte';
import astro from 'eslint-plugin-astro';

export default [
  {
    // Base configuration for all files
    ...eslint.configs.recommended,
    ignores: ['**/dist/**', '**/node_modules/**'],
  },
  {
    // TypeScript files configuration
    files: ['**/*.{ts,tsx}'],
    plugins: { '@typescript-eslint': typescript },
    languageOptions: {
      parser: tsParser,
      parserOptions: {
        project: './tsconfig.json',
        ecmaVersion: 2022,
        sourceType: 'module',
      },
    },
    rules: {
      ...typescript.configs.recommended.rules,
    },
  },
  {
    // Svelte files configuration
    files: ['**/*.svelte'],
    plugins: { svelte },
    languageOptions: {
      parser: svelteParser,
      parserOptions: {
        parser: tsParser,
        typescript: true,
      },
    },
    rules: {
      ...svelte.configs.recommended.rules,
    },
  },
  {
    // Astro files configuration
    files: ['**/*.astro'],
    plugins: { astro },
    languageOptions: {
      parser: astroParser,
      parserOptions: {
        parser: tsParser,
        extraFileExtensions: ['.astro'],
        sourceType: 'module',
      },
    },
    rules: {
      ...astro.configs['flat/recommended'].rules,
    },
  },
];
