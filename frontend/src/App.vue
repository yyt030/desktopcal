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
      @close="showPopup = false"
      @save-lines="onSaveLines"
      @set-bg="onSetBg"
    />
    <SettingsPopup
      :visible="showSettings"
      :settings="currentSettings"
      @close="showSettings = false"
      @save="onSaveSettings"
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
import { getMonthDays, createEvent, updateEvent, deleteEvent, exportData, importData, getSettings, saveSettings, setDayBackground, setWindowSize, setWindowAlwaysOnTop } from './api'
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

// 已渲染数据对应的「今天」。跨日检测据此判断，避免每分钟都重新拉数据。
let anchorToday = formatDate(new Date())
let dayCheckTimer: number | undefined

/** 跨日检测间隔：每 60 秒做一次纯日期比较，只有日期真的变了才重新加载。 */
const ROLLOVER_CHECK_MS = 60000

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
  firstDayOfWeek: 0,
  showLunar: true,
  showSolarTerm: true,
  layoutMode: 'compact',
  hotkey: 'cmd+shift+c',
  closeToTray: true,
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
    const existing = [...selectedEvents.value]
    const date = selectedDate.value

    const matched = new Set<string>()
    const toCreate: { title: string; completed: boolean }[] = []

    // 第一遍：按标题精确匹配，仅同步完成状态
    for (const line of lines) {
      const found = existing.find(e => e.title === line.title && !matched.has(e.id))
      if (found) {
        matched.add(found.id)
        if (found.completed !== line.completed) {
          await updateEvent({ ...found, completed: line.completed })
        }
      } else {
        toCreate.push(line)
      }
    }

    // 第二遍：未匹配上的行优先「改写」未匹配的已有事件，
    // 保留事件 ID 与创建时间，而不是删旧建新
    const orphans = existing.filter(e => !matched.has(e.id))
    const reused = Math.min(orphans.length, toCreate.length)
    for (let i = 0; i < reused; i++) {
      await updateEvent({
        ...orphans[i],
        title: toCreate[i].title,
        completed: toCreate[i].completed,
      })
      matched.add(orphans[i].id)
    }

    for (const e of existing) {
      if (!matched.has(e.id)) {
        await deleteEvent(e.id)
      }
    }

    for (const line of toCreate.slice(reused)) {
      await createEvent({
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
      })
    }
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

async function onExport() {
  try {
    const data = await exportData()
    const blob = new Blob([data], { type: 'application/json' })
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = `desktopcal-backup-${formatDate(new Date())}.json`
    a.click()
    URL.revokeObjectURL(url)
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

/** 应用设置：主题色/透明度/字号走 CSS 变量，窗口尺寸与置顶走原生窗口接口。 */
function applySettings(s: Settings) {
  applyAppearance(s)
  setWindowSize(s.width, s.height)
  setWindowAlwaysOnTop(s.alwaysTop)
}

async function onSaveSettings(s: Settings) {
  try {
    await saveSettings(s)
    currentSettings.value = s
    applySettings(s)
    await loadMonth()
  } catch (err) {
    showToast('保存设置失败：' + errMsg(err))
  }
}

onMounted(async () => {
  await loadMonth()
  const s = await getSettings()
  currentSettings.value = s
  applySettings(s)
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
