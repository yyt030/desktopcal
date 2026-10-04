<template>
  <Teleport to="body">
    <Transition name="popup">
      <div v-if="visible" class="popup-overlay" @click.self="handleClose">
        <div class="popup">
          <div class="popup-header">
            <div class="popup-header-left">
              <div class="popup-date">
                <span class="popup-day">{{ displayDay }}</span>
                <span class="popup-weekday">{{ weekdayName }}</span>
              </div>
              <div v-if="lunar && showLunar" class="popup-lunar">
                <span>{{ lunar.monthName }}{{ lunar.dayName }}</span>
                <span class="lunar-extra">{{ lunar.ganZhi }}年 {{ lunar.animal }}</span>
              </div>
              <div v-if="solarTerm && showSolarTerm" class="popup-solar-term">
                {{ solarTerm }}
              </div>
            </div>
            <div class="popup-header-right">
              <button class="btn-close" @click="handleClose" title="关闭">
                <svg width="14" height="14" viewBox="0 0 14 14" fill="currentColor">
                  <path d="M3.7 3.7a.5.5 0 0 1 .7 0L7 6.3l2.6-2.6a.5.5 0 0 1 .7.7L7.7 7l2.6 2.6a.5.5 0 0 1-.7.7L7 7.7l-2.6 2.6a.5.5 0 0 1-.7-.7L6.3 7 3.7 4.4a.5.5 0 0 1 0-.7z"/>
                </svg>
              </button>
            </div>
          </div>

          <div class="popup-body">
            <div class="line-numbers" ref="lineNumbersRef">
              <div v-for="n in lineCount" :key="n" class="line-num">{{ n }}.</div>
            </div>
            <textarea
              ref="textareaRef"
              v-model="text"
              class="events-textarea"
              placeholder="输入日程，每行一条..."
              :rows="Math.max(lineCount, 5)"
              @scroll="syncScroll"
              @keydown="handleKeydown"
              @contextmenu="handleContextMenu"
            ></textarea>
          </div>

          <div class="popup-footer">
            <span class="footer-hint">右键行首加 ✓ 标记完成</span>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import { ref, computed, watch, nextTick } from 'vue'
import type { Event, LunarDay } from '../types'

const props = defineProps<{
  visible: boolean
  date: string
  events: Event[]
  lunar: LunarDay | undefined
  solarTerm: string
  quickAddDate: string
  showLunar: boolean
  showSolarTerm: boolean
}>()

const emit = defineEmits<{
  close: []
  saveLines: [lines: { title: string; completed: boolean }[]]
}>()

const weekdays = ['周日', '周一', '周二', '周三', '周四', '周五', '周六']

const text = ref('')
const textareaRef = ref<HTMLTextAreaElement>()
const lineNumbersRef = ref<HTMLDivElement>()
const dirty = ref(false)

const displayDay = computed(() => {
  if (!props.date) return ''
  const d = new Date(props.date)
  return d.getDate()
})

const weekdayName = computed(() => {
  if (!props.date) return ''
  const d = new Date(props.date)
  return weekdays[d.getDay()]
})

const lineCount = computed(() => {
  if (!text.value) return 1
  const count = text.value.split('\n').length
  return Math.max(count, 1)
})

function eventsToText(events: Event[]): string {
  if (events.length === 0) return ''
  return events.map(e => {
    const prefix = e.completed ? '✓ ' : ''
    return prefix + e.title
  }).join('\n')
}

function parseLines(txt: string): { title: string; completed: boolean }[] {
  return txt.split('\n')
    .map(line => {
      const trimmed = line.trim()
      if (!trimmed) return null
      if (trimmed.startsWith('✓ ') || trimmed.startsWith('✓')) {
        return { title: trimmed.replace(/^✓\s*/, '').trim(), completed: true }
      }
      return { title: trimmed, completed: false }
    })
    .filter((l): l is { title: string; completed: boolean } => l !== null && l.title !== '')
}

function syncScroll() {
  if (lineNumbersRef.value && textareaRef.value) {
    lineNumbersRef.value.scrollTop = textareaRef.value.scrollTop
  }
}

function handleKeydown(e: KeyboardEvent) {
  if (e.key === 'Escape') {
    e.preventDefault()
    saveAndClose()
    return
  }
  if (e.key === 'Enter' && (e.ctrlKey || e.metaKey)) {
    e.preventDefault()
    saveAndClose()
  }
}

function handleContextMenu(e: MouseEvent) {
  e.preventDefault()
  if (!textareaRef.value) return
  const ta = textareaRef.value
  const pos = ta.selectionStart
  const lines = text.value.split('\n')
  let charCount = 0
  let lineIndex = 0
  for (let i = 0; i < lines.length; i++) {
    if (charCount + lines[i].length >= pos) {
      lineIndex = i
      break
    }
    charCount += lines[i].length + 1
    lineIndex = i
  }
  const line = lines[lineIndex]
  if (line.startsWith('✓ ')) {
    lines[lineIndex] = line.slice(2)
  } else if (line.startsWith('✓')) {
    lines[lineIndex] = line.slice(1).trimStart()
  } else {
    lines[lineIndex] = '✓ ' + line
  }
  text.value = lines.join('\n')
}

function saveAndClose() {
  if (dirty.value) {
    const lines = parseLines(text.value)
    emit('saveLines', lines)
    dirty.value = false
  }
  emit('close')
}

function handleClose() {
  saveAndClose()
}

watch(() => props.date, () => {
  text.value = eventsToText(props.events)
  dirty.value = false
})

watch(() => props.visible, (val) => {
  if (val) {
    nextTick(() => {
      textareaRef.value?.focus()
    })
  }
})

watch(() => props.events, (newEvents) => {
  if (!dirty.value) {
    text.value = eventsToText(newEvents)
  }
}, { deep: true })

watch(() => props.quickAddDate, (val) => {
  if (val) {
    nextTick(() => {
      textareaRef.value?.focus()
      const len = text.value.length
      if (len > 0 && !text.value.endsWith('\n')) {
        text.value += '\n'
      }
      nextTick(() => {
        if (textareaRef.value) {
          textareaRef.value.selectionStart = textareaRef.value.value.length
          textareaRef.value.selectionEnd = textareaRef.value.value.length
        }
      })
    })
  }
})

watch(text, () => {
  dirty.value = true
})
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

.popup {
  width: 420px;
  max-height: 70vh;
  display: flex;
  flex-direction: column;
  background: var(--bg-secondary);
  color: var(--text-primary);
  border-radius: 12px;
  border: 1px solid var(--border);
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.5), 0 0 0 1px rgba(255, 255, 255, 0.05);
  overflow: hidden;
}

.popup-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  padding: 20px 16px 16px;
  border-bottom: 1px solid var(--border);
}

.popup-header-left {
  flex: 1;
}

.popup-header-right {
  display: flex;
  gap: 6px;
  flex-shrink: 0;
}

.popup-date {
  display: flex;
  align-items: baseline;
  gap: 10px;
}

.popup-day {
  font-size: 42px;
  font-weight: 700;
  line-height: 1;
  background: linear-gradient(135deg, var(--text-primary), var(--text-secondary));
  -webkit-background-clip: text;
  background-clip: text;
  -webkit-text-fill-color: transparent;
}

.popup-weekday {
  font-size: 15px;
  color: var(--text-secondary);
}

.popup-lunar {
  margin-top: 8px;
  font-size: 13px;
  color: var(--text-secondary);
  display: flex;
  gap: 8px;
  align-items: center;
}

.lunar-extra {
  font-size: 11px;
  color: var(--text-muted);
}

.popup-solar-term {
  margin-top: 6px;
  font-size: 12px;
  color: #e8a0a0;
  font-weight: 500;
}

.btn-close {
  width: 28px;
  height: 28px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 6px;
  color: var(--text-secondary);
  background: var(--bg-hover);
  transition: all 0.2s;
}

.btn-close:hover {
  background: var(--bg-active);
  color: var(--text-primary);
}

.popup-body {
  flex: 1;
  display: flex;
  overflow-y: auto;
  min-height: 0;
  padding: 8px 0;
}

.line-numbers {
  flex-shrink: 0;
  padding: 0 4px 0 12px;
  overflow: hidden;
  user-select: none;
}

.line-num {
  font-size: var(--editor-font-size, 13px);
  line-height: 1.6;
  height: 1.6em;
  color: var(--text-muted);
  text-align: right;
  min-width: 20px;
  font-variant-numeric: tabular-nums;
}

.events-textarea {
  flex: 1;
  background: transparent;
  border: none;
  color: var(--text-primary);
  font-size: var(--editor-font-size, 13px);
  line-height: 1.6;
  resize: none;
  outline: none;
  padding: 0 12px 0 4px;
  font-family: inherit;
  min-height: 120px;
}

.events-textarea::placeholder {
  color: var(--text-muted);
}

.popup-footer {
  padding: 6px 16px;
  border-top: 1px solid var(--border);
}

.footer-hint {
  font-size: 11px;
  color: var(--text-muted);
}

.popup-enter-active, .popup-leave-active {
  transition: all 0.25s ease;
}

.popup-enter-from, .popup-leave-to {
  opacity: 0;
}

.popup-enter-from .popup, .popup-leave-to .popup {
  transform: scale(0.95) translateY(10px);
}
</style>
