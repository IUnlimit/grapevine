import 'vuetify/styles'
import '@mdi/font/css/materialdesignicons.css'
import { createVuetify } from 'vuetify'
import * as components from 'vuetify/components'
import * as directives from 'vuetify/directives'

const grapePurple = {
  dark: false,
  colors: {
    background: '#FAFAFA',
    surface: '#FFFFFF',
    primary: '#6A1B9A',
    'primary-darken-1': '#4A148C',
    secondary: '#AB47BC',
    'secondary-darken-1': '#7B1FA2',
    error: '#E53935',
    info: '#1E88E5',
    success: '#43A047',
    warning: '#FB8C00',
    'on-background': '#1C1B1F',
    'on-surface': '#1C1B1F',
  },
}

const grapeDark = {
  dark: true,
  colors: {
    background: '#1A1A2E',
    surface: '#16213E',
    primary: '#CE93D8',
    'primary-darken-1': '#AB47BC',
    secondary: '#BA68C8',
    'secondary-darken-1': '#9C27B0',
    error: '#EF5350',
    info: '#42A5F5',
    success: '#66BB6A',
    warning: '#FFA726',
    'on-background': '#E6E1E5',
    'on-surface': '#E6E1E5',
  },
}

export default createVuetify({
  components,
  directives,
  theme: {
    defaultTheme: 'grapePurple',
    themes: {
      grapePurple,
      grapeDark,
    },
  },
  defaults: {
    VCard: { elevation: 2, rounded: 'lg' },
    VBtn: { rounded: 'lg' },
    VTextField: { variant: 'outlined', density: 'comfortable' },
    VSelect: { variant: 'outlined', density: 'comfortable' },
  },
})
