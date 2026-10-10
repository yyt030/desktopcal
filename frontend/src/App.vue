<template>
  <div class="app-container">
    <AppHeader
      :year="currentYear"
      :month="currentMonth"
      @prev="prevMonth"
      @next="nextMonth"
      @today="goToday"
      @export="onExport"
      @import="onImport"
      @settings="showSettings = true"
    />
    <div class="app-body">
      <CalendarGrid
        :days="days"
        :selected-date="selectedDate"
        :first-day-of-week="currentSettings.firstDayOfWeek"
        :show-lunar="currentSettings.showLunar"
        :show-solar-term="currentSettings.showSolarTerm"
        @select-day="onSelectDay"
        @quick-add="onDblClickDay"
      />
    </div>
    <EventPopup
      :visible="showPopup"
      :date="selectedDate"
      :events="selectedEvents"
      :lunar="selectedLunar"
      :solar-term="selectedSolarTerm"
      :quick-add-token="quickAddToken"
      :show-lunar="currentSettings.showLunar"
      :show-solar-term="currentSettings.showSolarTerm"
      :bg-color="selectedDay?.bgColor ?? ''"
      :editor-width="currentSettings.editorWidth"
      :editor-height="currentSettings.editorHeight"
      @close="showPopup = false"
      @save-lines="onSaveLines"
      @set-bg="onSetBg"
      @resize="onEditorResize"
    />
    <SettingsPopup
      :visible="showSettings"
      :settings="currentSettings"
      @close="showSettings = false"
      @save="onSaveSettings"
      @toast="showToast"
    />
    <Teleport to="body">
      <Transition name="toast">
        <div v-if="toast" class="toast">{{ toast }}</div>
      </Transition>
    </Teleport>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onBeforeUnmount } from 'vue'
import AppHeader from './components/AppHeader.vue'
import CalendarGrid from './components/CalendarGrid.vue'
import EventPopup from './components/EventPopup.vue'
import SettingsPopup from './components/SettingsPopup.vue'
import { getMonthDays, createEvent, updateEvent, deleteEvent, batchSaveEvents, exportData, importData, getSettings, saveSettings, saveEditorSize, setDayBackground, setWindowSize, setWindowAlwaysOnTop } from './api'
import { formatDate, isSameMonth } from './utils/date'
import { applyAppearance } from './theme'
import type { DayInfo, Settings } from './types'

const now = new Date()
const currentYear = ref(now.getFullYear())
const currentMonth = ref(now.getMonth() + 1)
const days = ref<DayInfo[]>([])
const selectedDate = ref(formatDate(now))
// 每次双击日格自增，作为弹窗「新增一行」的信号
const quickAddToken = ref(0)
const showPopup = ref(false)
const showSettings = ref(false)
const toast = ref('')
let toastTimer: number | undefined
let editorResizeTimer: number | undefined

// 已渲染数据对应的「今天」。跨日检测据此判断，避免每分钟都重新拉数据。
let anchorToday = formatDate(new Date())
let dayCheckTimer: number | undefined

/** 跨日检测间隔：每 60 秒做一次纯日期比较，只有日期真的变了才重新加载。 */
const ROLLOVER_CHECK_MS = 180000

function showToast(msg: string) {
  toast.value = msg
  if (toastTimer) window.clearTimeout(toastTimer)
  toastTimer = window.setTimeout(() => {
    toast.value = ''
  }, 3000)
}

function errMsg(err: unknown): string {
  return err instanceof Error ? err.message : String(err)
}
const currentSettings = ref<Settings>({
  opacity: 100,
  width: 1100,
  height: 750,
  theme: 'dark',
  alwaysTop: true,
  fontSize: 13,
  firstDayOfWeek: 1,
  showLunar: true,
  showSolarTerm: true,
  layoutMode: 'loose',
  hotkey: 'cmd+shift+c',
  closeToTray: true,
  exportPath: '',
  editorWidth: 420,
  editorHeight: 400,
})

const selectedDay = computed(() =>
  days.value.find(d => d.solarDate === selectedDate.value)
)
const selectedEvents = computed(() => selectedDay.value?.events ?? [])
const selectedLunar = computed(() => selectedDay.value?.lunar)
const selectedSolarTerm = computed(() => selectedDay.value?.solarTerm ?? '')

async function loadMonth() {
  try {
    days.value = await getMonthDays(currentYear.value, currentMonth.value)
    // 拉取成功才更新锚点：数据已反映这里取的日期，失败时保留旧值以便下次重试。
    anchorToday = formatDate(new Date())
  } catch (err) {
    showToast('加载日历失败：' + errMsg(err))
  }
}

/**
 * 跨日检测：应用常驻托盘、长期不关闭，「今天」标记是后端按取数时刻算的，
 * 跨过午夜后若不重新取数就会停留在昨天。定时器本身可能被隐藏窗口节流，
 * 所以判断依据始终是「已渲染的日期 ≠ 当前日期」，而非定时器是否准点触发。
 */
function checkRollover() {
  const today = formatDate(new Date())
  if (today === anchorToday) return

  // 视图仍停在旧「今天」所在月份时，跟随跨日进入新月份；
  // 用户已手动翻到其它月份则不打断，只刷新数据让标记回到正确位置。
  if (isSameMonth(currentYear.value, currentMonth.value, anchorToday)) {
    const t = new Date()
    currentYear.value = t.getFullYear()
    currentMonth.value = t.getMonth() + 1
  }
  loadMonth()
}

function prevMonth() {
  if (currentMonth.value === 1) {
    currentMonth.value = 12
    currentYear.value--
  } else {
    currentMonth.value--
  }
  loadMonth()
}

function nextMonth() {
  if (currentMonth.value === 12) {
    currentMonth.value = 1
    currentYear.value++
  } else {
    currentMonth.value++
  }
  loadMonth()
}

function goToday() {
  const t = new Date()
  currentYear.value = t.getFullYear()
  currentMonth.value = t.getMonth() + 1
  selectedDate.value = formatDate(t)
  loadMonth()
}

function onSelectDay(date: string) {
  selectedDate.value = date
}

function onDblClickDay(date: string) {
  selectedDate.value = date
  showPopup.value = true
  quickAddToken.value++
}

async function onSaveLines(lines: { title: string; completed: boolean }[]) {
  try {
    const date = selectedDate.value
    // 将编辑器行转换为事件对象，使用批量保存（单次事务）
    const events = lines.map(line => ({
      id: '',
      title: line.title,
      description: '',
      date,
      startTime: '',
      endTime: '',
      color: '#4A90D9',
      completed: line.completed,
      createdAt: '',
      updatedAt: '',
    }))
    await batchSaveEvents(date, events)
  } catch (err) {
    showToast('保存事件失败：' + errMsg(err))
  }

  await loadMonth()
}

async function onSetBg(color: string) {
  try {
    await setDayBackground(selectedDate.value, color)
    await loadMonth()
    showToast(color ? '已设置当天背景色' : '已清除背景色')
  } catch (err) {
    showToast('设置背景色失败：' + errMsg(err))
  }
}

// 编辑器缩放防抖：500ms 内多次触发只保存最后一次，避免拖拽时频繁写数据库
function onEditorResize(width: number, height: number) {
  currentSettings.value.editorWidth = width
  currentSettings.value.editorHeight = height
  if (editorResizeTimer) window.clearTimeout(editorResizeTimer)
  editorResizeTimer = window.setTimeout(async () => {
    try {
      await saveEditorSize(width, height)
    } catch (err) {
      console.error('保存编辑器大小失败:', err)
    }
  }, 500)
}

async function onExport() {
  try {
    await exportData(currentYear.value, currentMonth.value)
    showToast('数据已导出')
  } catch (err) {
    showToast('导出失败：' + errMsg(err))
  }
}

async function onImport() {
  const input = document.createElement('input')
  input.type = 'file'
  input.accept = '.json'
  input.onchange = async () => {
    const file = input.files?.[0]
    if (!file) return
    try {
      const text = await file.text()
      await importData(text)
      await loadMonth()
      showToast('数据导入成功')
    } catch (err) {
      showToast('导入失败，文件格式可能不正确：' + errMsg(err))
    }
  }
  input.click()
}

async function onSaveSettings(s: Settings) {
  try {
    const oldSettings = { ...currentSettings.value }
    await saveSettings(s)
    currentSettings.value = s
    // 只在窗口尺寸设置改变时才调整窗口大小
    if (oldSettings.width !== s.width || oldSettings.height !== s.height) {
      setWindowSize(s.width, s.height)
    }
    // 置顶状态改变时才设置
    if (oldSettings.alwaysTop !== s.alwaysTop) {
      setWindowAlwaysOnTop(s.alwaysTop)
    }
    // 外观设置（主题、透明度、字号）总是应用
    applyAppearance(s)
    await loadMonth()
  } catch (err) {
    showToast('保存设置失败：' + errMsg(err))
  }
}

onMounted(async () => {
  await loadMonth()
  const s = await getSettings()
  currentSettings.value = s
  applyAppearance(s)
  setWindowSize(s.width, s.height)
  setWindowAlwaysOnTop(s.alwaysTop)
  dayCheckTimer = window.setInterval(checkRollover, ROLLOVER_CHECK_MS)
})

onBeforeUnmount(() => {
  if (dayCheckTimer !== undefined) window.clearInterval(dayCheckTimer)
})
</script>

<style scoped>
.app-container {
  display: flex;
  flex-direction: column;
  height: 100vh;
  padding: 8px 12px;
  gap: 8px;
  opacity: var(--app-opacity, 1);
}

.app-body {
  display: flex;
  flex: 1;
  min-height: 0;
}

.toast {
  position: fixed;
  left: 50%;
  bottom: 24px;
  transform: translateX(-50%);
  z-index: 2000;
  padding: 8px 16px;
  border-radius: 8px;
  background: var(--bg-secondary);
  color: var(--text-primary);
  border: 1px solid var(--border);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.35);
  font-size: 13px;
  max-width: 80vw;
}

.toast-enter-active,
.toast-leave-active {
  transition: all 0.2s ease;
}

.toast-enter-from,
.toast-leave-to {
  opacity: 0;
  transform: translateX(-50%) translateY(6px);
}
</style>
