/// <reference types="vite/client" />

declare module 'vuetify/styles' {
  const styles: string
  export default styles
}

declare module 'vuetify' {
  export * from 'vuetify/lib/framework.mjs'
}

declare module 'vuetify/components' {
  const components: Record<string, any>
  export = components
}

declare module 'vuetify/directives' {
  const directives: Record<string, any>
  export = directives
}

declare module '@mdi/font/css/materialdesignicons.css' {
  const css: string
  export default css
}

declare module 'vue-chartjs' {
  import { DefineComponent } from 'vue'
  export const Doughnut: DefineComponent<any, any, any>
  export const Bar: DefineComponent<any, any, any>
  export const Line: DefineComponent<any, any, any>
  export const Pie: DefineComponent<any, any, any>
}
