<template>
  <Teleport to="body">
    <Transition name="popup">
      <div v-if="visible" class="popup-overlay" @click.self="emit('close')">
        <div class="popup settings-popup">
          <div class="popup-header">
            <h2 class="settings-title">设置</h2>
            <button class="btn-close" @click="emit('close')" title="关闭">
              <svg width="14" height="14" viewBox="0 0 14 14" fill="currentColor">
                <path d="M3.7 3.7a.5.5 0 0 1 .7 0L7 6.3l2.6-2.6a.5.5 0 0 1 .7.7L7.7 7l2.6 2.6a.5.5 0 0 1-.7.7L7 7.7l-2.6 2.6a.5.5 0 0 1-.7-.7L6.3 7 3.7 4.4a.5.5 0 0 1 0-.7z"/>
              </svg>
            </button>
          </div>

          <div class="popup-body">
            <nav class="settings-nav">
              <button
                v-for="tab in tabs"
                :key="tab.value"
                class="nav-item"
                :class="{ active: activeTab === tab.value }"
                @click="activeTab = tab.value"
              >{{ tab.label }}</button>
            </nav>

            <div class="settings-content">
              <template v-if="activeTab === 'layout'">
                <div class="setting-group">
                  <label class="setting-label">显示模式</label>
                  <div class="setting-control">
                    <div class="option-btns">
                      <button
                        class="option-btn"
                        :class="{ active: localSettings.layoutMode === 'loose' }"
                        @click="localSettings.layoutMode = 'loose'"
                      >宽松</button>
                      <button
                        class="option-btn"
                        :class="{ active: localSettings.layoutMode === 'compact' }"
                        @click="localSettings.layoutMode = 'compact'"
                      >紧凑</button>
                    </div>
                    <span class="mode-hint">{{ localSettings.layoutMode === 'loose' ? '每月单独一页，不混入其它月' : '月历连续排布，交接处无空格' }}</span>
                  </div>
                </div>
                <div class="setting-group">
                  <label class="setting-label">每周第一天</label>
                  <div class="setting-control">
                    <div class="option-btns">
                      <button
                        class="option-btn"
                        :class="{ active: localSettings.firstDayOfWeek === 1 }"
                        @click="localSettings.firstDayOfWeek = 1"
                      >周一</button>
                      <button
                        class="option-btn"
                        :class="{ active: localSettings.firstDayOfWeek === 0 }"
                        @click="localSettings.firstDayOfWeek = 0"
                      >周日</button>
                    </div>
                  </div>
                </div>
                <div class="setting-group">
                  <label class="setting-label">显示农历</label>
                  <div class="setting-control">
                    <button
                      class="toggle-btn"
                      :class="{ active: localSettings.showLunar }"
                      @click="localSettings.showLunar = !localSettings.showLunar"
                    >{{ localSettings.showLunar ? '开启' : '关闭' }}</button>
                  </div>
                </div>
                <div class="setting-group">
                  <label class="setting-label">显示节气</label>
                  <div class="setting-control">
                    <button
                      class="toggle-btn"
                      :class="{ active: localSettings.showSolarTerm }"
                      @click="localSettings.showSolarTerm = !localSettings.showSolarTerm"
                    >{{ localSettings.showSolarTerm ? '开启' : '关闭' }}</button>
                  </div>
                </div>
              </template>

              <template v-else-if="activeTab === 'appearance'">
                <div class="setting-group">
                  <label class="setting-label">颜色主题</label>
                  <div class="setting-control">
                    <div class="option-btns">
                      <button
                        v-for="t in themes"
                        :key="t.value"
                        class="option-btn"
                        :class="{ active: localSettings.theme === t.value }"
                        @click="localSettings.theme = t.value"
                      >
                        <span class="theme-preview" :style="{ background: t.preview }"></span>
                        <span>{{ t.label }}</span>
                      </button>
                    </div>
                  </div>
                </div>
                <div class="setting-group">
                  <label class="setting-label">窗口透明度</label>
                  <div class="setting-control">
                    <input
                      type="range"
                      min="30"
                      max="100"
                      v-model.number="localSettings.opacity"
                      class="range-input"
                    />
                    <span class="setting-value">{{ localSettings.opacity }}%</span>
                  </div>
                </div>
                <div class="setting-group">
                  <label class="setting-label">编辑框字体大小</label>
                  <div class="setting-control">
                    <input
                      type="range"
                      min="10"
                      max="24"
                      v-model.number="localSettings.fontSize"
                      class="range-input"
                    />
                    <span class="setting-value">{{ localSettings.fontSize }}px</span>
                  </div>
                </div>
              </template>

              <template v-else-if="activeTab === 'window'">
                <div class="setting-group">
                  <label class="setting-label">窗口宽度</label>
                  <div class="setting-control">
                    <input
                      type="range"
                      min="800"
                      max="1600"
                      step="50"
                      v-model.number="localSettings.width"
                      class="range-input"
                    />
                    <span class="setting-value">{{ localSettings.width }}px</span>
                  </div>
                </div>
                <div class="setting-group">
                  <label class="setting-label">窗口高度</label>
                  <div class="setting-control">
                    <input
                      type="range"
                      min="500"
                      max="1000"
                      step="50"
                      v-model.number="localSettings.height"
                      class="range-input"
                    />
                    <span class="setting-value">{{ localSettings.height }}px</span>
                  </div>
                </div>
                <div class="setting-group">
                  <label class="setting-label">窗口置顶</label>
                  <div class="setting-control">
                    <button
                      class="toggle-btn"
                      :class="{ active: localSettings.alwaysTop }"
                      @click="localSettings.alwaysTop = !localSettings.alwaysTop"
                    >{{ localSettings.alwaysTop ? '开启' : '关闭' }}</button>
                  </div>
                </div>
                <div class="setting-group">
                  <label class="setting-label">关闭时最小化到托盘</label>
                  <div class="setting-control">
                    <button
                      class="toggle-btn"
                      :class="{ active: localSettings.closeToTray }"
                      @click="localSettings.closeToTray = !localSettings.closeToTray"
                    >{{ localSettings.closeToTray ? '开启' : '关闭' }}</button>
                    <span class="setting-hint">{{ localSettings.closeToTray ? '关闭窗口后保留 in 系统托盘' : '关闭窗口后退出程序' }}</span>
                  </div>
                </div>
              </template>

              <template v-else-if="activeTab === 'data'">
                <div class="setting-group">
                  <label class="setting-label">数据存储路径</label>
                  <div class="setting-control">
                    <div class="path-display">
                      <span class="path-text" :title="dbPath">{{ dbPath || '加载中...' }}</span>
                    </div>
                  </div>
                </div>
                <div class="setting-group">
                  <label class="setting-label">导出目录</label>
                  <div class="setting-control">
                    <div class="path-display">
                      <span class="path-text" :title="localSettings.exportPath || '未设置'">
                        {{ localSettings.exportPath || '未设置（默认使用主目录）' }}
                      </span>
                    </div>
                    <button class="btn-pick-path" @click="pickExportPath">选择</button>
                  </div>
                </div>
                <div class="setting-group">
                  <label class="setting-label">导出范围</label>
                  <div class="setting-control">
                    <div class="option-btns">
                      <button
                        class="option-btn"
                        :class="{ active: exportRange === 'month' }"
                        @click="exportRange = 'month'"
                      >单月</button>
                      <button
                        class="option-btn"
                        :class="{ active: exportRange === 'half' }"
                        @click="exportRange = 'half'"
                      >半年</button>
                    </div>
                  </div>
                </div>
                <div class="setting-group">
                  <label class="setting-label">导出操作</label>
                  <div class="setting-control">
                    <button class="btn-export" @click="doExport">立即导出</button>
                    <span class="setting-hint">{{ exportHint }}</span>
                  </div>
                </div>
              </template>

              <template v-else-if="activeTab === 'hotkey'">
                <div class="setting-group">
                  <label class="setting-label">显示/隐藏窗口</label>
                  <div class="setting-control">
                    <div class="hotkey-input-wrapper">
                      <input
                        type="text"
                        class="hotkey-input"
                        :value="formatHotkey(localSettings.hotkey)"
                        @keydown.prevent="captureHotkey"
                        @focus="isCapturing = true"
                        @blur="isCapturing = false"
                        placeholder="按下快捷键..."
                        readonly
                      />
                      <button class="btn-reset" @click="resetHotkey" title="重置为默认">
                        <svg width="14" height="14" viewBox="0 0 14 14" fill="currentColor">
                          <path d="M7 1a6 6 0 1 0 6 6h-1.5A4.5 4.5 0 1 1 7 2.5V1z"/>
                          <path d="M5.5 1H10v1.5H5.5z"/>
                        </svg>
                      </button>
                    </div>
                    <span class="mode-hint">{{ isCapturing ? '请按下新的快捷键...' : '点击输入框后按下快捷键' }}</span>
                  </div>
                </div>
              </template>

              <template v-else-if="activeTab === 'about'">
                <div class="about-section">
                  <div class="about-app-name">桌面日历</div>
                  <div class="about-author">作者：yueyt</div>
                  <div class="about-version">{{ versionInfo }}</div>
                </div>
              </template>
            </div>
          </div>

          <div class="popup-footer">
            <span class="footer-tab-name">{{ currentTabLabel }}</span>
            <div v-if="activeTab !== 'about'" class="footer-actions">
              <button class="btn-form btn-cancel" @click="emit('close')">取消</button>
              <button class="btn-form btn-save" @click="save">保存</button>
            </div>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { THEMES } from '../theme'
import { getVersion, pickDirectory, exportData, exportDataRange, getDatabasePath } from '../api'
import type { Settings } from '../types'

const props = defineProps<{
  visible: boolean
  settings: Settings
}>()

const emit = defineEmits<{
  close: []
  save: [settings: Settings]
  toast: [message: string]
}>()

const tabs = [
  { value: 'layout', label: '日历布局' },
  { value: 'appearance', label: '外观' },
  { value: 'window', label: '窗口' },
  { value: 'data', label: '数据' },
  { value: 'hotkey', label: '快捷键' },
  { value: 'about', label: '关于' },
]

// 主题可选项直接复用 theme.ts 的定义，保证设置页与实际生效的主题始终一致
const themes = THEMES

const activeTab = ref('layout')

const currentTabLabel = computed(
  () => tabs.find(t => t.value === activeTab.value)?.label ?? ''
)

const localSettings = ref<Settings>({ ...props.settings })
const versionInfo = ref('')
const dbPath = ref('')
const isCapturing = ref(false)
const exportRange = ref<'month' | 'half'>('month')
const exportMessage = ref<{ type: 'success' | 'error'; text: string } | null>(null)

const exportHint = computed(() => {
  const now = new Date()
  const year = now.getFullYear()
  const month = now.getMonth() + 1
  if (exportRange.value === 'month') {
    return `导出当前月：${year}年${month}月`
  } else {
    const startMonth = month <= 6 ? 1 : 7
    const endMonth = month <= 6 ? 6 : 12
    return `导出半年：${year}年${startMonth}月 - ${month <= 6 ? 6 : 12}月`
  }
})

getVersion().then(v => { versionInfo.value = v })
getDatabasePath().then(p => { dbPath.value = p })

function formatHotkey(hotkey: string): string {
  if (!hotkey) return ''
  return hotkey.split('+').map(part => {
    const map: Record<string, string> = {
      cmd: '⌘', ctrl: 'Ctrl', alt: 'Alt', shift: 'Shift'
    }
    return map[part.toLowerCase()] || part.toUpperCase()
  }).join('+')
}

function captureHotkey(e: KeyboardEvent) {
  const parts: string[] = []
  if (e.metaKey) parts.push('cmd')
  if (e.ctrlKey) parts.push('ctrl')
  if (e.altKey) parts.push('alt')
  if (e.shiftKey) parts.push('shift')
  
  const key = e.key.toLowerCase()
  if (!['meta', 'control', 'alt', 'shift'].includes(key)) {
    parts.push(key)
    localSettings.value.hotkey = parts.join('+')
  }
}

function resetHotkey() {
  const isMac = navigator.platform.toUpperCase().indexOf('MAC') >= 0
  localSettings.value.hotkey = isMac ? 'cmd+shift+c' : 'ctrl+shift+c'
}

async function pickExportPath() {
  try {
    const dir = await pickDirectory()
    if (dir) {
      localSettings.value.exportPath = dir
    }
  } catch {
    // 用户取消或出错，忽略
  }
}

async function doExport() {
  try {
    const now = new Date()
    const year = now.getFullYear()
    const month = now.getMonth() + 1

    if (exportRange.value === 'month') {
      await exportData(year, month)
    } else {
      const startMonth = month <= 6 ? 1 : 7
      const endMonth = month <= 6 ? 6 : 12
      await exportDataRange(year, startMonth, year, endMonth)
    }
    emit('toast', '数据已导出')
  } catch (err) {
    emit('toast', '导出失败：' + (err instanceof Error ? err.message : String(err)))
  }
}

watch(() => props.visible, (val) => {
  if (val) {
    localSettings.value = { ...props.settings }
    activeTab.value = 'layout'
  }
})

function save() {
  emit('save', { ...localSettings.value })
  emit('close')
}
</script>

<style scoped>
.popup-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  backdrop-filter: blur(4px);
}

.settings-popup {
  width: 560px;
  max-height: 80vh;
  display: flex;
  flex-direction: column;
  background: var(--bg-secondary, #1e1e2e);
  border-radius: 12px;
  border: 1px solid var(--border, rgba(255, 255, 255, 0.1));
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.5);
  overflow: hidden;
  color: var(--text-primary, #fff);
}

.popup-header {
  flex-shrink: 0;
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 20px;
  border-bottom: 1px solid var(--border, rgba(255, 255, 255, 0.08));
}

.settings-title {
  font-size: 16px;
  font-weight: 600;
  color: var(--text-primary, #fff);
  margin: 0;
}

.btn-close {
  width: 28px;
  height: 28px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 6px;
  color: var(--text-secondary, rgba(255, 255, 255, 0.6));
  background: var(--bg-hover, rgba(255, 255, 255, 0.06));
  transition: all 0.2s;
}

.btn-close:hover {
  color: var(--text-primary, #fff);
}

.popup-body {
  flex: 1;
  min-height: 0;
  display: flex;
}

.settings-nav {
  flex-shrink: 0;
  width: 130px;
  display: flex;
  flex-direction: column;
  gap: 4px;
  padding: 16px 10px;
  border-right: 1px solid var(--border, rgba(255, 255, 255, 0.08));
  background: var(--bg-primary, rgba(255, 255, 255, 0.02));
}

.nav-item {
  padding: 10px 12px;
  border-radius: 8px;
  font-size: 13px;
  text-align: left;
  color: var(--text-secondary, rgba(255, 255, 255, 0.6));
  background: transparent;
  cursor: pointer;
  transition: all 0.15s;
}

.nav-item:hover {
  background: var(--bg-hover, rgba(255, 255, 255, 0.06));
  color: var(--text-primary, #fff);
}

.nav-item.active {
  background: var(--bg-hover, rgba(74, 144, 217, 0.18));
  color: var(--accent, #4A90D9);
  font-weight: 600;
}

.settings-content {
  flex: 1;
  min-width: 0;
  overflow-y: auto;
  padding: 20px;
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.settings-content::-webkit-scrollbar {
  width: 6px;
}

.settings-content::-webkit-scrollbar-thumb {
  background: var(--border, rgba(255, 255, 255, 0.15));
  border-radius: 3px;
}

.setting-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.setting-label {
  font-size: 13px;
  color: var(--text-primary, #fff);
  font-weight: 600;
}

.setting-control {
  display: flex;
  align-items: center;
  gap: 12px;
  flex-wrap: wrap;
}

.mode-hint {
  font-size: 11px;
  color: var(--text-secondary, rgba(255, 255, 255, 0.45));
}

.setting-hint {
  font-size: 11px;
  color: var(--text-secondary, rgba(255, 255, 255, 0.45));
  margin-left: 8px;
}

.option-btns {
  display: flex;
  gap: 8px;
}

.option-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 14px;
  border-radius: 8px;
  background: var(--bg-primary, rgba(255, 255, 255, 0.04));
  border: 1px solid var(--border, transparent);
  color: var(--text-secondary, rgba(255, 255, 255, 0.7));
  font-size: 12px;
  cursor: pointer;
  transition: all 0.15s;
}

.option-btn:hover {
  background: var(--bg-hover, rgba(255, 255, 255, 0.08));
  color: var(--text-primary, #fff);
}

.option-btn.active {
  border-color: var(--accent, #4A90D9);
  background: var(--bg-hover, rgba(255, 255, 255, 0.08));
  color: var(--text-primary, #fff);
  font-weight: 600;
}

.theme-preview {
  width: 18px;
  height: 18px;
  border-radius: 50%;
  border: 2px solid var(--border, rgba(255, 255, 255, 0.2));
}

.range-input {
  flex: 1;
  height: 4px;
  -webkit-appearance: none;
  appearance: none;
  background: var(--border, rgba(255, 255, 255, 0.1));
  border-radius: 2px;
  outline: none;
}

.range-input::-webkit-slider-thumb {
  -webkit-appearance: none;
  appearance: none;
  width: 14px;
  height: 14px;
  border-radius: 50%;
  background: var(--accent, #4A90D9);
  cursor: pointer;
}

.setting-value {
  font-size: 12px;
  color: var(--text-primary, rgba(255, 255, 255, 0.5));
  min-width: 50px;
  text-align: right;
}

.toggle-btn {
  padding: 6px 16px;
  border-radius: 6px;
  font-size: 12px;
  background: var(--bg-primary, rgba(255, 255, 255, 0.06));
  border: 1px solid var(--border, transparent);
  color: var(--text-secondary, rgba(255, 255, 255, 0.6));
  transition: all 0.15s;
}

.toggle-btn.active {
  background: var(--accent, #4A90D9);
  border-color: var(--accent, #4A90D9);
  color: #fff;
  font-weight: 600;
}

.toggle-btn:hover {
  background: var(--bg-hover, rgba(255, 255, 255, 0.1));
}

.about-section {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  padding: 30px 0;
}

.about-app-name {
  font-size: 18px;
  font-weight: 700;
  color: var(--text-primary, #fff);
}

.about-author {
  font-size: 13px;
  color: var(--text-secondary, rgba(255, 255, 255, 0.6));
}

.about-version {
  font-size: 12px;
  color: var(--text-secondary, rgba(255, 255, 255, 0.45));
  font-family: monospace;
}

.popup-footer {
  flex-shrink: 0;
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 8px;
  padding: 14px 20px;
  border-top: 1px solid var(--border, rgba(255, 255, 255, 0.08));
}

.footer-tab-name {
  font-size: 12px;
  color: var(--text-secondary, rgba(255, 255, 255, 0.4));
}

.footer-actions {
  display: flex;
  gap: 8px;
}

.btn-form {
  padding: 8px 16px;
  border-radius: 6px;
  font-size: 13px;
  font-weight: 500;
}

.btn-cancel {
  color: var(--text-secondary, rgba(255, 255, 255, 0.6));
  background: transparent;
}

.btn-cancel:hover {
  color: var(--text-primary, #fff);
}

.btn-save {
  background: var(--accent, #4A90D9);
  color: #fff;
  font-weight: 600;
}

.popup-enter-active, .popup-leave-active {
  transition: all 0.25s ease;
}

.popup-enter-from, .popup-leave-to {
  opacity: 0;
}

.popup-enter-from .settings-popup, .popup-leave-to .settings-popup {
  transform: scale(0.95) translateY(10px);
}

.hotkey-input-wrapper {
  display: flex;
  align-items: center;
  gap: 8px;
}

.hotkey-input {
  flex: 1;
  padding: 8px 12px;
  border-radius: 6px;
  background: var(--bg-primary, rgba(255, 255, 255, 0.04));
  border: 1px solid var(--border, rgba(255, 255, 255, 0.1));
  color: var(--text-primary, #fff);
  font-size: 13px;
  font-family: monospace;
  outline: none;
  transition: border-color 0.15s;
}

.hotkey-input:focus {
  border-color: var(--accent, #4A90D9);
}

.btn-reset {
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 6px;
  background: var(--bg-primary, rgba(255, 255, 255, 0.04));
  border: 1px solid var(--border, rgba(255, 255, 255, 0.1));
  color: var(--text-secondary, rgba(255, 255, 255, 0.6));
  cursor: pointer;
  transition: all 0.15s;
}

.btn-reset:hover {
  background: var(--bg-hover, rgba(255, 255, 255, 0.08));
  color: var(--text-primary, #fff);
}

.path-display {
  flex: 1;
  min-width: 0;
  padding: 8px 12px;
  border-radius: 6px;
  background: var(--bg-primary, rgba(255, 255, 255, 0.04));
  border: 1px solid var(--border, rgba(255, 255, 255, 0.1));
}

.path-text {
  display: block;
  font-size: 12px;
  color: var(--text-secondary, rgba(255, 255, 255, 0.7));
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.btn-pick-path {
  padding: 8px 14px;
  border-radius: 6px;
  font-size: 12px;
  background: var(--bg-primary, rgba(255, 255, 255, 0.06));
  border: 1px solid var(--border, rgba(255, 255, 255, 0.1));
  color: var(--text-secondary, rgba(255, 255, 255, 0.7));
  cursor: pointer;
  transition: all 0.15s;
}

.btn-pick-path:hover {
  background: var(--bg-hover, rgba(255, 255, 255, 0.1));
  color: var(--text-primary, #fff);
}

.btn-export {
  padding: 8px 16px;
  border-radius: 6px;
  font-size: 13px;
  font-weight: 500;
  background: var(--accent, #4A90D9);
  color: #fff;
  border: none;
  cursor: pointer;
  transition: all 0.15s;
}

.btn-export:hover {
  background: var(--accent-hover, #3a7bc8);
}
</style>
