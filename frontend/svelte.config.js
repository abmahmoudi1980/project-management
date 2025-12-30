import { vitePreprocess } from '@sveltejs/vite-plugin-svelte';

export default {
  // Enable runes mode (Svelte 5)
  compilerOptions: {
    runes: true
  },
  preprocess: vitePreprocess()
};
