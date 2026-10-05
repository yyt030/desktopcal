<template>
  <header class="header" @mousedown="onDragStart">
    <div class="header-left">
      <h1 class="title">桌面日历</h1>
      <span class="date-display">{{ year }}年{{ month }}月</span>
    </div>
    <div class="header-right">
      <button class="btn btn-today" @click="$emit('today')">今天</button>
      <div class="data-group">
        <button class="btn btn-data" @click="$emit('export')" title="导出数据">
          <svg width="14" height="14" viewBox="0 0 14 14" fill="currentColor">
            <path d="M7 1a.5.5 0 0 1 .5.5v7.793l2.146-2.147a.5.5 0 0 1 .708.708l-3 3a.5.5 0 0 1-.708 0l-3-3a.5.5 0 1 1 .708-.708L6.5 9.293V1.5A.5.5 0 0 1 7 1zM2 11.5a.5.5 0 0 1 .5.5v1h9v-1a.5.5 0 0 1 1 0v1.5a.5.5 0 0 1-.5.5h-10a.5.5 0 0 1-.5-.5V12a.5.5 0 0 1 .5-.5z"/>
          </svg>
        </button>
        <button class="btn btn-data" @click="$emit('import')" title="导入数据">
          <svg width="14" height="14" viewBox="0 0 14 14" fill="currentColor">
            <path d="M7 10a.5.5 0 0 1-.5-.5V1.707L4.354 3.854a.5.5 0 1 1-.708-.708l3-3a.5.5 0 0 1 .708 0l3 3a.5.5 0 0 1-.708.708L7.5 1.707V9.5a.5.5 0 0 1-.5.5zM2 11.5a.5.5 0 0 1 .5.5v1h9v-1a.5.5 0 0 1 1 0v1.5a.5.5 0 0 1-.5.5h-10a.5.5 0 0 1-.5-.5V12a.5.5 0 0 1 .5-.5z"/>
          </svg>
        </button>
      </div>
      <button class="btn btn-settings" @click="$emit('settings')" title="设置">
        <svg width="16" height="16" viewBox="0 0 16 16" fill="currentColor">
          <path d="M8 4.754a3.246 3.246 0 1 0 0 6.492 3.246 3.246 0 0 0 0-6.492zM5.754 8a2.246 2.246 0 1 1 4.492 0 2.246 2.246 0 0 1-4.492 0z"/>
          <path d="M9.796 1.343c-.527-1.79-3.065-1.79-3.592 0l-.094.319a.873.873 0 0 1-1.255.52l-.292-.16c-1.64-.892-3.433.902-2.54 2.541l.159.292a.873.873 0 0 1-.52 1.255l-.319.094c-1.79.527-1.79 3.065 0 3.592l.319.094a.873.873 0 0 1 .52 1.255l-.16.292c-.892 1.64.901 3.434 2.541 2.54l.292-.159a.873.873 0 0 1 1.255.52l.094.319c.527 1.79 3.065 1.79 3.592 0l.094-.319a.873.873 0 0 1 1.255-.52l.292.16c1.64.893 3.434-.902 2.54-2.541l-.159-.292a.873.873 0 0 1 .52-1.255l.319-.094c1.79-.527 1.79-3.065 0-3.592l-.319-.094a.873.873 0 0 1-.52-1.255l.16-.292c.893-1.64-.902-3.433-2.541-2.54l-.292.159a.873.873 0 0 1-1.255-.52l-.094-.319zm-2.633.283c.246-.835 1.428-.835 1.674 0l.094.319a1.873 1.873 0 0 0 2.693 1.115l.291-.16c.764-.415 1.6.42 1.184 1.185l-.159.292a1.873 1.873 0 0 0 1.116 2.692l.318.094c.835.246.835 1.428 0 1.674l-.319.094a1.873 1.873 0 0 0-1.115 2.693l.16.291c.415.764-.42 1.6-1.185 1.184l-.291-.159a1.873 1.873 0 0 0-2.693 1.116l-.094.318c-.246.835-1.428.835-1.674 0l-.094-.319a1.873 1.873 0 0 0-2.692-1.115l-.292.16c-.764.415-1.6-.42-1.184-1.185l.159-.291A1.873 1.873 0 0 0 1.945 8.93l-.319-.094c-.835-.246-.835-1.428 0-1.674l.319-.094A1.873 1.873 0 0 0 3.06 4.377l-.16-.292c-.415-.764.42-1.6 1.185-1.184l.292.159a1.873 1.873 0 0 0 2.692-1.115l.094-.319z"/>
        </svg>
      </button>
      <div class="nav-group">
        <button class="btn btn-nav" @click="$emit('prev')">
          <svg width="16" height="16" viewBox="0 0 16 16" fill="currentColor">
            <path d="M10.354 3.354a.5.5 0 0 0-.708-.708l-5 5a.5.5 0 0 0 0 .708l5 5a.5.5 0 0 0 .708-.708L5.707 8l4.647-4.646z"/>
          </svg>
        </button>
        <button class="btn btn-nav" @click="$emit('next')">
          <svg width="16" height="16" viewBox="0 0 16 16" fill="currentColor">
            <path d="M5.646 3.354a.5.5 0 0 1 .708-.708l5 5a.5.5 0 0 1 0 .708l-5 5a.5.5 0 0 1-.708-.708L10.293 8 5.646 3.354z"/>
          </svg>
        </button>
      </div>
      <div class="window-controls">
        <button class="btn btn-window" @click="onMinimise" title="最小化">
          <svg width="12" height="12" viewBox="0 0 12 12" fill="currentColor">
            <path d="M2 6h8v1H2z"/>
          </svg>
        </button>
        <button class="btn btn-window" @click="onMaximise" title="最大化">
          <svg width="12" height="12" viewBox="0 0 12 12" fill="currentColor">
            <path d="M3 3h6v6H3V3zm1 1v4h4V4H4z"/>
          </svg>
        </button>
        <button class="btn btn-window btn-close" @click="onClose" title="关闭">
          <svg width="12" height="12" viewBox="0 0 12 12" fill="currentColor">
            <path d="M3.5 3.5l5 5m0-5l-5 5" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/>
          </svg>
        </button>
      </div>
    </div>
  </header>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { minimiseWindow, maximiseWindow, closeWindow, getWindowPosition, setWindowPosition } from '@/api'

defineProps<{
  year: number
  month: number
}>()

defineEmits<{
  prev: []
  next: []
  today: []
  export: []
  import: []
  settings: []
}>()

function onMinimise() {
  minimiseWindow()
}

function onMaximise() {
  maximiseWindow()
}

function onClose() {
  closeWindow()
}

// 窗口拖动
let isDragging = false
let startX = 0
let startY = 0
let startWinX = 0
let startWinY = 0

async function onDragStart(e: MouseEvent) {
  // 只在左键且不在右侧按钮区时开始拖动
  if (e.button !== 0) return
  const target = e.target as HTMLElement
  if (target.closest('.header-right')) return
  
  isDragging = true
  startX = e.screenX
  startY = e.screenY
  const pos = await getWindowPosition()
  startWinX = pos.x
  startWinY = pos.y
  
  document.addEventListener('mousemove', onDragMove)
  document.addEventListener('mouseup', onDragEnd)
}

function onDragMove(e: MouseEvent) {
  if (!isDragging) return
  const dx = e.screenX - startX
  const dy = e.screenY - startY
  setWindowPosition(startWinX + dx, startWinY + dy)
}

function onDragEnd() {
  isDragging = false
  document.removeEventListener('mousemove', onDragMove)
  document.removeEventListener('mouseup', onDragEnd)
}
</script>

<style scoped>
.header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 8px 4px;
  flex-shrink: 0;
}

.header-left {
  display: flex;
  align-items: baseline;
  gap: 16px;
  cursor: default;
  user-select: none;
}

.title {
  font-size: 20px;
  font-weight: 600;
  letter-spacing: 0.5px;
  background: linear-gradient(135deg, var(--text-primary), var(--text-secondary));
  -webkit-background-clip: text;
  background-clip: text;
  -webkit-text-fill-color: transparent;
}

.date-display {
  font-size: 15px;
  color: var(--text-secondary);
  font-weight: 400;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 12px;
}

.btn {
  padding: 6px 14px;
  border-radius: var(--radius-sm);
  font-size: 13px;
  transition: all 0.2s;
}

.btn-today {
  background: var(--bg-secondary);
  border: 1px solid var(--border);
}

.btn-today:hover {
  background: var(--bg-hover);
  border-color: var(--accent);
  color: var(--accent);
}

.data-group {
  display: flex;
  gap: 2px;
  background: var(--bg-secondary);
  border-radius: var(--radius-sm);
  border: 1px solid var(--border);
  overflow: hidden;
}

.btn-data {
  padding: 6px 8px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.btn-data:hover {
  background: var(--bg-hover);
  color: var(--accent);
}

.btn-settings {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 6px 8px;
  background: var(--bg-secondary);
  border: 1px solid var(--border);
  border-radius: var(--radius-sm);
}

.btn-settings:hover {
  background: var(--bg-hover);
  border-color: var(--accent);
  color: var(--accent);
}

.nav-group {
  display: flex;
  gap: 2px;
  background: var(--bg-secondary);
  border-radius: var(--radius-sm);
  border: 1px solid var(--border);
  overflow: hidden;
}

.btn-nav {
  padding: 6px 10px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.btn-nav:hover {
  background: var(--bg-hover);
}

.window-controls {
  display: flex;
  gap: 2px;
  margin-left: 8px;
  padding-left: 8px;
  border-left: 1px solid var(--border);
}

.btn-window {
  padding: 6px 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: transparent;
  border: none;
  border-radius: var(--radius-sm);
  color: var(--text-secondary);
}

.btn-window:hover {
  background: var(--bg-hover);
  color: var(--text-primary);
}

.btn-close:hover {
  background: #e81123;
  color: white;
}
</style>
