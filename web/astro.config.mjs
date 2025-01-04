// @ts-check
import { defineConfig } from 'astro/config';
import { env } from 'process';

import svelte from '@astrojs/svelte';

import node from '@astrojs/node';

const isProd = env.NODE_ENV === 'production';

// https://astro.build/config
export default defineConfig({
  integrations: [svelte({ extensions: ['.svelte'] })],

  adapter: node({
    mode: 'standalone',
  }),

  server: {
    port: isProd ? 4321 : 4320,
    host: true
  },
  output: 'static',
  build: {
    serverEntry: 'entry.mjs',
  },
  vite: {
    server: {
      watch: {
        ignored: ['**/node_modules/**', '**/dist/**'],
      },
      port: isProd ? 4321 : 4320,
      strictPort: true,
    },
  },
});
