import 'vuetify/styles'
import '@mdi/font/css/materialdesignicons.css'
import { createVuetify } from 'vuetify'
import * as components from 'vuetify/components'
import * as directives from 'vuetify/directives'

const grapeLight = {
  dark: false,
  colors: {
    background: '#F4F2F7',
    surface: '#FFFFFF',
    'surface-bright': '#FFFFFF',
    'surface-variant': '#EDE9F3',
    primary: '#7C3AED',
    'primary-darken-1': '#6D28D9',
    secondary: '#10B981',
    'secondary-darken-1': '#059669',
    error: '#EF4444',
    info: '#3B82F6',
    success: '#10B981',
    warning: '#F59E0B',
    'on-background': '#1E1B2E',
    'on-surface': '#1E1B2E',
    'on-primary': '#FFFFFF',
  },
}

const grapeDark = {
  dark: true,
  colors: {
    background: '#0F0B1A',
    surface: '#1A1528',
    'surface-bright': '#241E35',
    'surface-variant': '#2D2640',
    primary: '#A78BFA',
    'primary-darken-1': '#8B5CF6',
    secondary: '#34D399',
    'secondary-darken-1': '#10B981',
    error: '#F87171',
    info: '#60A5FA',
    success: '#34D399',
    warning: '#FBBF24',
    'on-background': '#E8E0F0',
    'on-surface': '#E8E0F0',
    'on-primary': '#0F0B1A',
  },
}

export default createVuetify({
  components,
  directives,
  theme: {
    defaultTheme: 'grapeDark',
    themes: {
      grapeLight,
      grapeDark,
    },
  },
  defaults: {
    VCard: {
      elevation: 0,
      rounded: 'xl',
    },
    VBtn: {
      rounded: 'lg',
    },
    VTextField: {
      variant: 'outlined',
      density: 'comfortable',
      rounded: 'lg',
    },
    VSelect: {
      variant: 'outlined',
      density: 'comfortable',
      rounded: 'lg',
    },
    VChip: {
      rounded: 'lg',
    },
  },
})
