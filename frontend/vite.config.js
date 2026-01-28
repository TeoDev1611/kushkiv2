import {defineConfig} from 'vite'
import {svelte} from '@sveltejs/vite-plugin-svelte'
import sveltePreprocess from 'svelte-preprocess'
import path from 'path'

// https://vitejs.dev/config/
export default defineConfig({
  resolve: {
    alias: {
      'wailsjs': path.resolve(__dirname, './wailsjs'),
      '$lib': path.resolve(__dirname, './src/lib')
    }
  },
  plugins: [svelte({
    preprocess: sveltePreprocess({
      typescript: {
        tsconfigFile: './tsconfig.json'
      }
    })
  })]
})
