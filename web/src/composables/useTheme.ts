import { ref, watch } from 'vue'

type Theme = 'light' | 'dark'

const theme = ref<Theme>('light')

export function useTheme() {
  function init() {
    const saved = localStorage.getItem('gogeoip_theme') as Theme | null
    if (saved) {
      theme.value = saved
    } else if (window.matchMedia('(prefers-color-scheme: dark)').matches) {
      theme.value = 'dark'
    }
    applyTheme()
  }

  function applyTheme() {
    const root = document.documentElement
    if (theme.value === 'dark') {
      root.classList.add('dark')
    } else {
      root.classList.remove('dark')
    }
  }

  function toggle() {
    theme.value = theme.value === 'light' ? 'dark' : 'light'
    localStorage.setItem('gogeoip_theme', theme.value)
    applyTheme()
  }

  function set(t: Theme) {
    theme.value = t
    localStorage.setItem('gogeoip_theme', t)
    applyTheme()
  }

  watch(theme, applyTheme)

  return { theme, init, toggle, set }
}
