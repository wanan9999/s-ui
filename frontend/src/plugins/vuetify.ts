/**
 * plugins/vuetify.ts
 *
 * Framework documentation: https://vuetifyjs.com`
 */

// Styles
import '@mdi/font/css/materialdesignicons.css'
import 'vuetify/styles'

import { fa, en, vi, zhHans, zhHant, ru } from 'vuetify/locale'

// Composables
import { createVuetify } from 'vuetify'

// https://vuetifyjs.com/en/introduction/why-vuetify/#feature-guides
export default createVuetify({
  defaults: {
    VRow: { density: 'comfortable' },
    VCard: { elevation: 0, rounded: 'lg', border: true },
    VBtn: { elevation: 0, rounded: 'lg', variant: 'flat' },
    VDialog: { scrollable: true, transition: 'fade-transition' },
    VMenu: {
      location: 'bottom end', offset: 8,
      maxHeight: 'min(420px, calc(100dvh - 32px))',
      maxWidth: 'calc(100vw - 32px)',
      transition: 'fade-transition',
      VList: { density: 'compact', VSwitch: { density: 'compact' } },
    },
    VTooltip: { location: 'top', openDelay: 300, scrollStrategy: 'close' },
    VTabs: { showArrows: true, color: 'primary' },
    VField: { color: 'primary', rounded: 'lg' },
    VExpansionPanels: { elevation: 0, variant: 'accordion' },
    VSwitch: { density: 'comfortable', color: 'primary', hideDetails: 'auto', inset: true },
    VCheckbox: { density: 'comfortable', color: 'primary', hideDetails: 'auto' },
    VRadioGroup: { density: 'comfortable', color: 'primary', hideDetails: 'auto' },
    VFileInput: {
      variant: 'outlined', density: 'comfortable', hideDetails: 'auto',
    },
    VTextField: {
      variant: 'outlined', density: 'comfortable', hideDetails: 'auto',
    },
    VSelect: {
      variant: 'outlined', density: 'comfortable', hideDetails: 'auto',
    },
    VCombobox: {
      variant: 'outlined', density: 'comfortable', hideDetails: 'auto',
    },
    VTextarea: {
      variant: 'outlined', density: 'comfortable', hideDetails: 'auto',
    },
  },
  theme: {
    defaultTheme: localStorage.getItem('theme') ?? 'system',
    themes: {
      light: {
        dark: false,
        colors: {
          primary: '#5364D9', secondary: '#64748B',
          background: '#F4F6FA', surface: '#FFFFFF',
          'surface-variant': '#E5EAF3', 'on-surface-variant': '#475569',
          'on-background': '#243247', 'on-surface': '#243247',
          error: '#D9475C', success: '#16866C', warning: '#B87518', info: '#3173BA',
        },
        variables: { 'border-color': '#64748B', 'border-opacity': 0.15, 'high-emphasis-opacity': 1, 'medium-emphasis-opacity': 0.75 },
      },
      dark: {
        dark: true,
        colors: {
          primary: '#A2AEFF', secondary: '#9CAEC8',
          background: '#141C2A', surface: '#1D283A',
          'surface-variant': '#344158', 'on-surface-variant': '#CDD7E6',
          'on-background': '#E1E8F2', 'on-surface': '#E1E8F2',
          error: '#FF8194', success: '#5BC8A7', warning: '#F0BD6E', info: '#7CB5F1',
        },
        variables: { 'border-color': '#A5B4CC', 'border-opacity': 0.17, 'high-emphasis-opacity': 1, 'medium-emphasis-opacity': 0.75 },
      },
    },
  },
  locale: {
    locale: localStorage.getItem("locale") ?? 'zhHans',
    fallback: 'en',
    messages: { en, fa, vi, zhHans, zhHant, ru },
  },
})
