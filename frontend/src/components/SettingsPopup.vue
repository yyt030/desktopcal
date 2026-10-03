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
              <label class="setting-label">颜色主题</label>
              <div class="setting-control">
                <div class="theme-options">
                  <button
                    v-for="t in themes"
                    :key="t.value"
                    class="theme-btn"
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
              <label class="setting-label">窗口置顶</label>
              <div class="setting-control">
                <button
                  class="toggle-btn"
                  :class="{ active: localSettings.alwaysTop }"
                  @click="localSettings.alwaysTop = !localSettings.alwaysTop"
                >
                  {{ localSettings.alwaysTop ? '开启' : '关闭' }}
                </button>
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
          </div>

          <div class="popup-footer">
            <button class="btn-form btn-cancel" @click="emit('close')">取消</button>
            <button class="btn-form btn-save" @click="save">保存</button>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import type { Settings } from '../types'

const props = defineProps<{
  visible: boolean
  settings: Settings
}>()

const emit = defineEmits<{
  close: []
  save: [settings: Settings]
}>()

const themes = [
  { value: 'dark', label: '深色', preview: '#1e1e2e' },
  { value: 'light', label: '浅色', preview: '#f5f5f5' },
  { value: 'blue', label: '蓝色', preview: '#1a3a5c' },
  { value: 'green', label: '绿色', preview: '#1a3c2e' },
]

const localSettings = ref<Settings>({ ...props.settings })

watch(() => props.visible, (val) => {
  if (val) {
    localSettings.value = { ...props.settings }
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
  width: 400px;
  background: #1e1e2e;
  border-radius: 12px;
  border: 1px solid rgba(255, 255, 255, 0.1);
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.5);
  overflow: hidden;
}

.popup-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 20px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.08);
}

.settings-title {
  font-size: 16px;
  font-weight: 600;
  color: #fff;
  margin: 0;
}

.btn-close {
  width: 28px;
  height: 28px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 6px;
  color: rgba(255, 255, 255, 0.6);
  background: rgba(255, 255, 255, 0.06);
  transition: all 0.2s;
}

.btn-close:hover {
  background: rgba(255, 255, 255, 0.12);
  color: #fff;
}

.popup-body {
  padding: 20px;
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.setting-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.setting-label {
  font-size: 13px;
  color: rgba(255, 255, 255, 0.8);
  font-weight: 500;
}

.setting-control {
  display: flex;
  align-items: center;
  gap: 12px;
}

.range-input {
  flex: 1;
  height: 4px;
  -webkit-appearance: none;
  appearance: none;
  background: rgba(255, 255, 255, 0.1);
  border-radius: 2px;
  outline: none;
}

.range-input::-webkit-slider-thumb {
  -webkit-appearance: none;
  appearance: none;
  width: 14px;
  height: 14px;
  border-radius: 50%;
  background: #4A90D9;
  cursor: pointer;
}

.setting-value {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.5);
  min-width: 50px;
  text-align: right;
}

.theme-options {
  display: flex;
  gap: 8px;
}

.theme-btn {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
  padding: 8px 12px;
  border-radius: 8px;
  background: rgba(255, 255, 255, 0.04);
  border: 1px solid transparent;
  color: rgba(255, 255, 255, 0.7);
  font-size: 11px;
  cursor: pointer;
  transition: all 0.15s;
}

.theme-btn:hover {
  background: rgba(255, 255, 255, 0.08);
}

.theme-btn.active {
  border-color: #4A90D9;
  color: #fff;
}

.theme-preview {
  width: 24px;
  height: 24px;
  border-radius: 50%;
  border: 2px solid rgba(255, 255, 255, 0.2);
}

.toggle-btn {
  padding: 6px 16px;
  border-radius: 6px;
  font-size: 12px;
  background: rgba(255, 255, 255, 0.06);
  color: rgba(255, 255, 255, 0.6);
  transition: all 0.15s;
}

.toggle-btn.active {
  background: #4A90D9;
  color: #fff;
}

.toggle-btn:hover {
  background: rgba(255, 255, 255, 0.1);
}

.toggle-btn.active:hover {
  background: #3a7bc8;
}

.popup-footer {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
  padding: 16px 20px;
  border-top: 1px solid rgba(255, 255, 255, 0.08);
}

.btn-form {
  padding: 8px 16px;
  border-radius: 6px;
  font-size: 13px;
  font-weight: 500;
}

.btn-cancel {
  color: rgba(255, 255, 255, 0.6);
  background: transparent;
}

.btn-cancel:hover {
  color: #fff;
}

.btn-save {
  background: #4A90D9;
  color: #fff;
}

.btn-save:hover {
  background: #3a7bc8;
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
</style>
