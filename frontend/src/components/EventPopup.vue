<template>
  <Teleport to="body">
    <Transition name="popup">
      <div v-if="visible" class="popup-overlay">
        <div
          ref="popupRef"
          class="popup"
          :style="{ width: localWidth + 'px', height: localHeight + 'px' }"
        >
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
            <div v-if="lineCount > 1" class="line-numbers" ref="lineNumbersRef">
              <div v-for="n in lineCount" :key="n" class="line-num" :ref="el => { if (el) lineNumRefs[n - 1] = el as HTMLElement }">{{ n }}.</div>
            </div>
            <div
              ref="editorRef"
              class="events-editor"
              contenteditable="true"
              spellcheck="false"
              @scroll="syncScroll"
              @keydown="handleEditorKeydown"
              @contextmenu="handleContextMenu"
              @input="onEditorInput"
              v-editor-content="{ text: text, version: lineSyncVersion }"
            ></div>
          </div>

          <div class="popup-footer">
            <div class="bg-row">
              <span class="bg-label">背景</span>
              <button
                v-for="c in bgColors"
                :key="c"
                class="bg-swatch"
                :class="{ active: bgColor === c }"
                :style="{ background: c }"
                :title="c"
                @click="emit('setBg', c)"
              ></button>
              <button
                class="bg-clear"
                :class="{ active: !bgColor }"
                @click="emit('setBg', '')"
              >清除</button>
            </div>
            <span class="footer-hint">右键切换完成状态</span>
          </div>
          <div
            class="resize-handle"
            @mousedown.prevent="startResize"
          ></div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import { ref, computed, watch, nextTick, type ObjectDirective } from 'vue'
import type { Event, LunarDay } from '../types'
import { parseLocalDate } from '../utils/date'

const props = defineProps<{
  visible: boolean
  date: string
  events: Event[]
  lunar?: LunarDay | null
  solarTerm: string
  /** 每次双击日期自增，作为「新增一行并聚焦」的信号（自增可保证重复双击也能触发） */
  quickAddToken: number
  showLunar: boolean
  showSolarTerm: boolean
  bgColor: string
  editorWidth: number
  editorHeight: number
}>()

const emit = defineEmits<{
  close: []
  saveLines: [lines: { title: string; completed: boolean }[]]
  setBg: [color: string]
  resize: [width: number, height: number]
}>()

const weekdays = ['周日', '周一', '周二', '周三', '周四', '周五', '周六']

const bgColors = ['#E74C3C', '#3498DB', '#2ECC71']

const text = ref('')
const editorRef = ref<HTMLDivElement>()
const lineNumbersRef = ref<HTMLDivElement>()
const popupRef = ref<HTMLDivElement>()
const lineNumRefs = ref<HTMLElement[]>([])
const dirty = ref(false)
// 本地尺寸状态：打开时从 props 初始化，拖拽缩放时实时更新，关闭时保存到父组件
const localWidth = ref(props.editorWidth)
const localHeight = ref(props.editorHeight)
// 由程序（切换日期 / 双击新增行）改写 text 时，不应把内容标记为「已修改」，
// 避免关闭时触发一次无意义的保存
let suppressDirty = false

const MIN_WIDTH = 300
const MIN_HEIGHT = 250

// 控制 contenteditable 内容何时从数据同步到 DOM。
// 父容器整体可编辑，子 div 仅做视觉分行。只在结构性变化时重建 DOM，
// 普通输入不触发重建，保证光标和选区不被打断。
const lineSyncVersion = ref(0)

function renderEditorContent(el: HTMLElement, txt: string) {
  const lineArr = txt ? txt.split('\n') : ['']
  // 保存选区
  const sel = window.getSelection()
  let savedOffset = -1
  if (sel && sel.rangeCount > 0 && el.contains(sel.anchorNode)) {
    const range = sel.getRangeAt(0)
    const preRange = range.cloneRange()
    preRange.selectNodeContents(el)
    preRange.setEnd(range.startContainer, range.startOffset)
    savedOffset = preRange.toString().length
  }

  el.innerHTML = ''
  for (const line of lineArr) {
    const div = document.createElement('div')
    div.className = 'editor-line'
    const isCompleted = line.startsWith('✓ ') || line.startsWith('✓')
    if (isCompleted) {
      div.classList.add('completed')
      div.textContent = line || '\u00A0'
    } else {
      div.textContent = line || '\u00A0'
    }
    el.appendChild(div)
  }

  // 恢复选区
  if (savedOffset >= 0 && sel) {
    restoreCaretOffset(el, savedOffset)
  }
  // 同步行号高度
  nextTick(() => syncLineHeights())
}

function restoreCaretOffset(root: HTMLElement, targetOffset: number) {
  const walker = document.createTreeWalker(root, NodeFilter.SHOW_TEXT)
  let offset = 0
  let node: Text | null = null
  while (walker.nextNode()) {
    node = walker.currentNode as Text
    const len = node.textContent?.length ?? 0
    if (offset + len >= targetOffset) {
      const range = document.createRange()
      const sel = window.getSelection()
      range.setStart(node, targetOffset - offset)
      range.collapse(true)
      sel?.removeAllRanges()
      sel?.addRange(range)
      return
    }
    offset += len
  }
  // 超出范围，放到最后
  const lastDiv = root.lastElementChild as HTMLElement | null
  if (lastDiv) {
    const range = document.createRange()
    const sel = window.getSelection()
    range.selectNodeContents(lastDiv)
    range.collapse(false)
    sel?.removeAllRanges()
    sel?.addRange(range)
  }
}

const vEditorContent: ObjectDirective<HTMLElement, { text: string; version: number }> = {
  mounted(el, binding) {
    renderEditorContent(el, binding.value.text)
    ;(el as any).__lv = binding.value.version
  },
  updated(el, binding) {
    if (binding.value.version !== (el as any).__lv) {
      renderEditorContent(el, binding.value.text)
      ;(el as any).__lv = binding.value.version
    }
  },
}

function startResize(e: MouseEvent) {
  const startX = e.clientX
  const startY = e.clientY
  const startW = localWidth.value
  const startH = localHeight.value
  const maxW = window.innerWidth * 0.9
  const maxH = window.innerHeight * 0.9

  // rAF 节流：mousemove 高频触发，直接写响应式变量会每事件重渲染一次
  let raf = 0
  let pendingW = 0
  let pendingH = 0
  function onMove(ev: MouseEvent) {
    pendingW = Math.min(maxW, Math.max(MIN_WIDTH, startW + ev.clientX - startX))
    pendingH = Math.min(maxH, Math.max(MIN_HEIGHT, startH + ev.clientY - startY))
    if (raf) return
    raf = requestAnimationFrame(() => {
      raf = 0
      localWidth.value = pendingW
      localHeight.value = pendingH
    })
  }
  function onUp() {
    // 有未落帧的尺寸时先落地，避免松手时停在上一帧
    if (raf) {
      cancelAnimationFrame(raf)
      raf = 0
      localWidth.value = pendingW
      localHeight.value = pendingH
    }
    document.removeEventListener('mousemove', onMove)
    document.removeEventListener('mouseup', onUp)
  }
  document.addEventListener('mousemove', onMove)
  document.addEventListener('mouseup', onUp)
}

const displayDay = computed(() => {
  if (!props.date) return ''
  return parseLocalDate(props.date).getDate()
})

const weekdayName = computed(() => {
  if (!props.date) return ''
  return weekdays[parseLocalDate(props.date).getDay()]
})

const lines = computed(() => {
  if (!text.value) return ['']
  return text.value.split('\n')
})

const lineCount = computed(() => lines.value.length)

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
  if (lineNumbersRef.value && editorRef.value) {
    lineNumbersRef.value.scrollTop = editorRef.value.scrollTop
  }
}

// 使用 requestAnimationFrame 节流行号高度同步，避免频繁重排
// 双重 rAF 确保布局完全完成后再测量
let syncLineHeightsRaf: number | undefined
function syncLineHeights() {
  if (syncLineHeightsRaf) return
  syncLineHeightsRaf = requestAnimationFrame(() => {
    requestAnimationFrame(() => {
      syncLineHeightsRaf = undefined
      const editor = editorRef.value
      if (!editor) return
      const lineEls = editor.querySelectorAll('.editor-line')
      lineNumRefs.value.forEach((lineNum, idx) => {
        const lineEl = lineEls[idx]
        if (lineNum && lineEl) {
          const height = lineEl.getBoundingClientRect().height
          lineNum.style.height = height + 'px'
          lineNum.style.lineHeight = height + 'px'
        }
      })
    })
  })
}

/** 用事件列表刷新文本，并保持「未修改」状态 */
function syncTextFromEvents() {
  suppressDirty = true
  text.value = eventsToText(props.events)
  lineSyncVersion.value++
  dirty.value = false
  nextTick(() => {
    suppressDirty = false
  })
}

function focusEnd() {
  const editor = editorRef.value
  if (!editor) return
  editor.focus()
  const lineEls = editor.querySelectorAll('.editor-line')
  const lastLine = lineEls[lineEls.length - 1] as HTMLElement | undefined
  if (lastLine) {
    const range = document.createRange()
    const sel = window.getSelection()
    range.selectNodeContents(lastLine)
    range.collapse(false)
    sel?.removeAllRanges()
    sel?.addRange(range)
  }
}

/** 双击日期：在末尾追加一个空行并把光标移到行尾，方便直接输入 */
function appendNewLine() {
  nextTick(() => {
    suppressDirty = true
    if (text.value.length > 0 && !text.value.endsWith('\n')) {
      text.value += '\n'
    }
    lineSyncVersion.value++
    nextTick(() => {
      focusEnd()
      suppressDirty = false
    })
  })
}

/** 从 DOM 读回各行内容，同步到 text.value */
function readLinesFromDOM(): string[] {
  const editor = editorRef.value
  if (!editor) return ['']
  const lineEls = editor.querySelectorAll('.editor-line')
  if (lineEls.length === 0) return ['']
  const result: string[] = []
  lineEls.forEach(el => {
    const content = (el.textContent || '').replace(/\u00A0/g, '')
    // 如果该行有 completed 类，保留 ✓ 前缀
    const wasCompleted = el.classList.contains('completed')
    if (wasCompleted && content) {
      result.push('✓ ' + content)
    } else {
      result.push(content)
    }
  })
  return result
}

/**
 * 自动删除空行（非末行、非光标所在行）。空行保存时本就进不了任务列表，
 * 编辑期就地清掉保证所见即所得。只移除光标所在行之外的节点，
 * 光标所在的文本节点不受影响，无需恢复选区。
 */
function removeEmptyLines(editor: HTMLElement) {
  const lineEls = Array.from(editor.querySelectorAll('.editor-line'))
  if (lineEls.length === 0) return
  const caretIdx = getCaretLineIndex(editor)
  for (let i = 0; i < lineEls.length - 1; i++) {
    if (i === caretIdx) continue
    const el = lineEls[i]
    const content = (el.textContent || '').replace(/\u00A0/g, '')
    if (content === '') el.remove()
  }
}

function onEditorInput() {
  const editor = editorRef.value
  if (!editor) return
  // 规范化：确保每个直接子元素是 div.editor-line
  normalizeEditorDOM(editor)
  removeEmptyLines(editor)
  const newLines = readLinesFromDOM()
  text.value = newLines.join('\n')
  // 同步行号高度
  nextTick(() => syncLineHeights())
}

/** 确保编辑器每个直接子节点是 div.editor-line */
function normalizeEditorDOM(editor: HTMLElement) {
  const children = Array.from(editor.childNodes)
  for (const child of children) {
    if (child.nodeType === Node.TEXT_NODE) {
      // 裸文本节点 → 包进 div
      const text = child.textContent || ''
      if (text.trim() === '') {
        child.remove()
        continue
      }
      const div = document.createElement('div')
      div.className = 'editor-line'
      div.textContent = text
      editor.replaceChild(div, child)
    } else if (child.nodeType === Node.ELEMENT_NODE) {
      const el = child as HTMLElement
      if (el.tagName === 'BR') {
        // <br> → 空行 div
        const div = document.createElement('div')
        div.className = 'editor-line'
        div.textContent = '\u00A0'
        editor.replaceChild(div, el)
      } else if (!el.classList.contains('editor-line')) {
        el.classList.add('editor-line')
      }
    }
  }
  // 保证至少有一行
  if (editor.children.length === 0) {
    const div = document.createElement('div')
    div.className = 'editor-line'
    div.textContent = '\u00A0'
    editor.appendChild(div)
  }
}

function handleEditorKeydown(e: KeyboardEvent) {
  if (e.key === 'Escape') {
    e.preventDefault()
    saveAndClose()
    return
  }
  if (e.key === 'Enter' && (e.ctrlKey || e.metaKey)) {
    e.preventDefault()
    saveAndClose()
    return
  }
  if (e.key === 'Enter') {
    e.preventDefault()
    insertNewLine()
    return
  }
  // Backspace 在空行开头时删除该行（除非只剩一行）
  if (e.key === 'Backspace') {
    const editor = editorRef.value
    if (!editor) return
    const lineEls = editor.querySelectorAll('.editor-line')
    if (lineEls.length <= 1) return
    const caretInfo = getCaretLineIndex(editor)
    if (caretInfo >= 0) {
      const curLine = lineEls[caretInfo] as HTMLElement
      const content = (curLine.textContent || '').replace(/\u00A0/g, '')
      if (content === '') {
        e.preventDefault()
        curLine.remove()
        // 确保至少一行
        if (editor.querySelectorAll('.editor-line').length === 0) {
          const div = document.createElement('div')
          div.className = 'editor-line'
          div.textContent = '\u00A0'
          editor.appendChild(div)
        }
        const newLines = readLinesFromDOM()
        text.value = newLines.join('\n')
        // 焦点移到上一行末尾
        nextTick(() => {
          const updatedLines = editor.querySelectorAll('.editor-line')
          const prevIdx = Math.max(0, caretInfo - 1)
          const prevLine = updatedLines[prevIdx] as HTMLElement | undefined
          if (prevLine) {
            prevLine.focus()
            const range = document.createRange()
            const sel = window.getSelection()
            range.selectNodeContents(prevLine)
            range.collapse(false)
            sel?.removeAllRanges()
            sel?.addRange(range)
          }
        })
      }
    }
  }
}

/** 在当前光标位置所在行的下方插入新行 */
function insertNewLine() {
  const editor = editorRef.value
  if (!editor) return
  const sel = window.getSelection()
  if (!sel || sel.rangeCount === 0) return
  const range = sel.getRangeAt(0)

  // 找到光标所在的行 div
  let lineEl: Node | null = range.startContainer
  while (lineEl && lineEl !== editor && (lineEl as HTMLElement)?.className !== 'editor-line') {
    lineEl = lineEl.parentNode
  }
  if (!lineEl || lineEl === editor) return

  const lineDiv = lineEl as HTMLElement
  const newDiv = document.createElement('div')
  newDiv.className = 'editor-line'
  newDiv.textContent = '\u00A0'

  // 如果光标在行尾，在后面插入；否则在当前行后插入
  if (lineDiv.nextSibling) {
    editor.insertBefore(newDiv, lineDiv.nextSibling)
  } else {
    editor.appendChild(newDiv)
  }

  // 移动光标到新行开头
  const newRange = document.createRange()
  newRange.setStart(newDiv, 0)
  newRange.collapse(true)
  sel.removeAllRanges()
  sel.addRange(newRange)

  const newLines = readLinesFromDOM()
  text.value = newLines.join('\n')
}

/** 获取光标所在行在编辑器中的索引 */
function getCaretLineIndex(editor: HTMLElement): number {
  const sel = window.getSelection()
  if (!sel || sel.rangeCount === 0) return -1
  let node: Node | null = sel.getRangeAt(0).startContainer
  while (node && node !== editor && (node as HTMLElement)?.className !== 'editor-line') {
    node = node.parentNode
  }
  if (!node || node === editor) return -1
  const lineEls = Array.from(editor.querySelectorAll('.editor-line'))
  return lineEls.indexOf(node as HTMLElement)
}

function handleContextMenu(e: MouseEvent) {
  e.preventDefault()
  let target = e.target as Node
  if (target.nodeType === Node.TEXT_NODE) {
    target = target.parentElement!
  }
  const lineEl = (target as HTMLElement).closest('.editor-line') as HTMLElement | undefined
  if (!lineEl) return
  const editor = editorRef.value
  if (!editor) return
  const lineEls = Array.from(editor.querySelectorAll('.editor-line'))
  const idx = lineEls.indexOf(lineEl)
  if (idx < 0) return

  const newLines = [...lines.value]
  const line = newLines[idx] ?? ''
  if (line.startsWith('✓ ')) {
    newLines[idx] = line.slice(2)
  } else if (line.startsWith('✓')) {
    newLines[idx] = line.slice(1).trimStart()
  } else {
    newLines[idx] = '✓ ' + line
  }
  text.value = newLines.join('\n')
  lineSyncVersion.value++
}

function saveAndClose() {
  if (popupRef.value) {
    const w = Math.round(popupRef.value.offsetWidth)
    const h = Math.round(popupRef.value.offsetHeight)
    if (w !== props.editorWidth || h !== props.editorHeight) {
      emit('resize', w, h)
    }
  }
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
  syncTextFromEvents()
})

watch(() => props.visible, (val) => {
  if (val) {
    localWidth.value = props.editorWidth
    localHeight.value = props.editorHeight
    nextTick(() => {
      focusEnd()
    })
  }
})

watch(() => props.events, (newEvents) => {
  if (dirty.value) return
  suppressDirty = true
  text.value = eventsToText(newEvents)
  lineSyncVersion.value++
  nextTick(() => {
    suppressDirty = false
  })
}, { deep: true })

watch(() => props.quickAddToken, () => {
  appendNewLine()
})

watch(text, () => {
  if (!suppressDirty) dirty.value = true
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
}

.popup {
  display: flex;
  flex-direction: column;
  background: var(--bg-secondary);
  color: var(--text-primary);
  border-radius: 12px;
  border: 1px solid var(--border);
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.5), 0 0 0 1px rgba(255, 255, 255, 0.05);
  overflow: hidden;
  position: relative;
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
  color: var(--text-muted);
  text-align: right;
  min-width: 20px;
  font-variant-numeric: tabular-nums;
  box-sizing: border-box;
  padding-top: 0;
  margin: 0;
}

.events-editor {
  flex: 1;
  background: transparent;
  color: var(--text-primary);
  font-size: var(--editor-font-size, 13px);
  line-height: 1.6;
  outline: none;
  padding: 0 12px 0 4px;
  font-family: inherit;
  min-height: 120px;
  overflow-y: auto;
  cursor: text;
}

.editor-line {
  min-height: 1.6em;
  white-space: pre-wrap;
  word-break: break-word;
}

.editor-line.completed {
  text-decoration: line-through;
  color: var(--text-secondary);
}

.popup-footer {
  padding: 8px 16px;
  border-top: 1px solid var(--border);
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}

.bg-row {
  display: flex;
  align-items: center;
  gap: 6px;
  flex-shrink: 0;
}

.bg-label {
  font-size: 11px;
  color: var(--text-muted);
}

.bg-swatch {
  width: 16px;
  height: 16px;
  border-radius: 50%;
  border: 2px solid transparent;
  transition: transform 0.15s;
}

.bg-swatch:hover {
  transform: scale(1.15);
}

.bg-swatch.active {
  border-color: var(--text-primary);
}

.bg-clear {
  font-size: 11px;
  color: var(--text-secondary);
  padding: 2px 8px;
  border-radius: 4px;
  background: var(--bg-hover);
}

.bg-clear:hover {
  color: var(--text-primary);
}

.bg-clear.active {
  color: var(--text-primary);
}

.footer-hint {
  font-size: 11px;
  color: var(--text-muted);
  white-space: nowrap;
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

.resize-handle {
  position: absolute;
  right: 0;
  bottom: 0;
  width: 16px;
  height: 16px;
  cursor: nwse-resize;
  z-index: 10;
}

.resize-handle::after {
  content: '';
  position: absolute;
  right: 3px;
  bottom: 3px;
  width: 8px;
  height: 8px;
  border-right: 2px solid var(--text-muted);
  border-bottom: 2px solid var(--text-muted);
  opacity: 0.5;
}

.resize-handle:hover::after {
  opacity: 1;
}
</style>
