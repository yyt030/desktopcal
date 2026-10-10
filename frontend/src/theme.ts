// 主题定义：颜色变量是实现主题切换的唯一来源，
// App.vue 应用主题、SettingsPopup 展示可选项都从这里读取，避免两处各写一份。
import type { Settings } from './types'

export interface ThemeOption {
  /** 与 settings.theme 对应的标识 */
  value: string
  /** 设置页展示的名称 */
  label: string
  /** 设置页的圆形色块预览色 */
  preview: string
  /** 写入 :root 的 CSS 变量 */
  vars: Record<string, string>
}

/** 默认主题标识，theme 非法时回退到它。 */
export const DEFAULT_THEME = 'dark'

export const THEMES: ThemeOption[] = [
  {
    value: 'dark',
    label: '深色',
    preview: '#1e1e2e',
    vars: {
      '--bg-primary': '#1e1e2e',
      '--bg-secondary': '#2a2a3e',
      '--bg-hover': '#33334d',
      '--bg-active': 'rgba(74,144,217,0.28)',
      '--today-bg': 'rgba(74,144,217,0.18)',
      '--today-border': 'rgba(74,144,217,0.6)',
      '--weekend': 'rgba(232,120,120,0.85)',
      '--weekend-bg': 'rgba(0,0,0,0.12)',
      '--text-primary': '#ffffff',
      '--text-secondary': 'rgba(255,255,255,0.6)',
      '--text-muted': 'rgba(255,255,255,0.3)',
      '--border': 'rgba(255,255,255,0.1)',
      '--accent': '#4A90D9',
      '--accent-hover': '#5BA0E9',
    },
  },
  {
    value: 'light',
    label: '浅色',
    preview: '#f5f5f5',
    vars: {
      '--bg-primary': '#f5f5f5',
      '--bg-secondary': '#e8e8e8',
      '--bg-hover': '#ddd',
      '--bg-active': 'rgba(37,99,235,0.16)',
      '--today-bg': 'rgba(37,99,235,0.12)',
      '--today-border': 'rgba(37,99,235,0.5)',
      '--weekend': 'rgba(200,70,70,0.9)',
      '--weekend-bg': 'rgba(0,0,0,0.04)',
      '--text-primary': '#1a1a1a',
      '--text-secondary': 'rgba(0,0,0,0.55)',
      '--text-muted': 'rgba(0,0,0,0.35)',
      '--border': 'rgba(0,0,0,0.12)',
      '--accent': '#2563eb',
      '--accent-hover': '#1d4ed8',
    },
  },
  {
    value: 'blue',
    label: '蓝色',
    preview: '#1a3a5c',
    vars: {
      '--bg-primary': '#1a3a5c',
      '--bg-secondary': '#1e4470',
      '--bg-hover': '#254e80',
      '--bg-active': 'rgba(91,163,245,0.3)',
      '--today-bg': 'rgba(91,163,245,0.18)',
      '--today-border': 'rgba(91,163,245,0.6)',
      '--weekend': 'rgba(232,140,140,0.85)',
      '--weekend-bg': 'rgba(0,0,0,0.12)',
      '--text-primary': '#e8f0ff',
      '--text-secondary': 'rgba(232,240,255,0.6)',
      '--text-muted': 'rgba(232,240,255,0.35)',
      '--border': 'rgba(255,255,255,0.12)',
      '--accent': '#5ba3f5',
      '--accent-hover': '#7bb8f8',
    },
  },
  {
    value: 'green',
    label: '绿色',
    preview: '#1a3c2e',
    vars: {
      '--bg-primary': '#1a3c2e',
      '--bg-secondary': '#1e4a36',
      '--bg-hover': '#256040',
      '--bg-active': 'rgba(76,175,80,0.3)',
      '--today-bg': 'rgba(76,175,80,0.18)',
      '--today-border': 'rgba(76,175,80,0.6)',
      '--weekend': 'rgba(232,140,140,0.85)',
      '--weekend-bg': 'rgba(0,0,0,0.12)',
      '--text-primary': '#e8ffe8',
      '--text-secondary': 'rgba(232,255,232,0.6)',
      '--text-muted': 'rgba(232,255,232,0.35)',
      '--border': 'rgba(255,255,255,0.12)',
      '--accent': '#4caf50',
      '--accent-hover': '#66bb6a',
    },
  },
]

/** 按标识查找主题，未命中时回退到默认主题。 */
export function resolveTheme(theme: string): ThemeOption {
  return THEMES.find(t => t.value === theme) ?? THEMES.find(t => t.value === DEFAULT_THEME)!
}

/** 把主题色、透明度、编辑框字号写入 :root，供各组件通过 CSS 变量消费。 */
export function applyAppearance(settings: Settings): void {
  const root = document.documentElement
  for (const [key, value] of Object.entries(resolveTheme(settings.theme).vars)) {
    root.style.setProperty(key, value)
  }
  root.style.setProperty('--app-opacity', String(settings.opacity / 100))
  root.style.setProperty('--editor-font-size', `${settings.fontSize}px`)
}
