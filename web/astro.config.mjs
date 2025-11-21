// @ts-check
import { defineConfig } from 'astro/config';
import { env } from 'process';
import svelte from '@astrojs/svelte';

const port = typeof process !== 'undefined' && env.PORT 
  ? parseInt(env.PORT) 
  : 4320;
const isPreview = typeof process !== 'undefined' && env.NODE_ENV === 'preview';

// https://astro.build/config
/** @type {import('astro').AstroUserConfig} */
export default defineConfig({
  integrations: [svelte()],
  
  output: 'static',
  
  server: {
    port,
    host: true
  },

  vite: {
    server: {
      watch: {
        ignored: ['**/node_modules/**', '**/dist/**'],
      },
      hmr: {
        protocol: 'ws',
        clientPort: isPreview ? port : undefined,
        host: isPreview ? '0.0.0.0' : undefined
      }
    }
  }
});
