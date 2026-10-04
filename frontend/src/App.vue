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
      :quick-add-date="quickAddDate"
      :show-lunar="currentSettings.showLunar"
      :show-solar-term="currentSettings.showSolarTerm"
      @close="showPopup = false"
      @save-lines="onSaveLines"
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
import { ref, computed, onMounted } from 'vue'
import AppHeader from './components/AppHeader.vue'
import CalendarGrid from './components/CalendarGrid.vue'
import EventPopup from './components/EventPopup.vue'
import SettingsPopup from './components/SettingsPopup.vue'
import { getMonthDays, createEvent, updateEvent, deleteEvent, exportData, importData, getSettings, saveSettings, setWindowSize, setWindowAlwaysOnTop } from './api'
import type { DayInfo, Settings } from './types'

const now = new Date()
const currentYear = ref(now.getFullYear())
const currentMonth = ref(now.getMonth() + 1)
const days = ref<DayInfo[]>([])
const selectedDate = ref(formatDate(now))
const quickAddDate = ref('')
const showPopup = ref(false)
const showSettings = ref(false)
const toast = ref('')
let toastTimer: number | undefined

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
  } catch (err) {
    showToast('加载日历失败：' + errMsg(err))
  }
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
}

async function onSaveLines(lines: { title: string; completed: boolean }[]) {
  try {
    const existing = [...selectedEvents.value]
    const date = selectedDate.value

    const matched = new Set<string>()
    const toCreate: { title: string; completed: boolean }[] = []

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

    for (const e of existing) {
      if (!matched.has(e.id)) {
        await deleteEvent(e.id)
      }
    }

    for (const line of toCreate) {
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

const themeVars: Record<string, Record<string, string>> = {
  dark: {
    '--bg-primary': '#1e1e2e',
    '--bg-secondary': '#2a2a3e',
    '--bg-hover': '#33334d',
    '--bg-active': 'rgba(74,144,217,0.28)',
    '--today-bg': 'rgba(74,144,217,0.18)',
    '--weekend': 'rgba(232,120,120,0.85)',
    '--text-primary': '#ffffff',
    '--text-secondary': 'rgba(255,255,255,0.6)',
    '--text-muted': 'rgba(255,255,255,0.3)',
    '--border': 'rgba(255,255,255,0.1)',
    '--accent': '#4A90D9',
    '--accent-hover': '#5BA0E9',
  },
  light: {
    '--bg-primary': '#f5f5f5',
    '--bg-secondary': '#e8e8e8',
    '--bg-hover': '#ddd',
    '--bg-active': 'rgba(37,99,235,0.16)',
    '--today-bg': 'rgba(37,99,235,0.12)',
    '--weekend': 'rgba(200,70,70,0.9)',
    '--text-primary': '#1a1a1a',
    '--text-secondary': 'rgba(0,0,0,0.55)',
    '--text-muted': 'rgba(0,0,0,0.35)',
    '--border': 'rgba(0,0,0,0.12)',
    '--accent': '#2563eb',
    '--accent-hover': '#1d4ed8',
  },
  blue: {
    '--bg-primary': '#1a3a5c',
    '--bg-secondary': '#1e4470',
    '--bg-hover': '#254e80',
    '--bg-active': 'rgba(91,163,245,0.3)',
    '--today-bg': 'rgba(91,163,245,0.18)',
    '--weekend': 'rgba(232,140,140,0.85)',
    '--text-primary': '#e8f0ff',
    '--text-secondary': 'rgba(232,240,255,0.6)',
    '--text-muted': 'rgba(232,240,255,0.35)',
    '--border': 'rgba(255,255,255,0.12)',
    '--accent': '#5ba3f5',
    '--accent-hover': '#7bb8f8',
  },
  green: {
    '--bg-primary': '#1a3c2e',
    '--bg-secondary': '#1e4a36',
    '--bg-hover': '#256040',
    '--bg-active': 'rgba(76,175,80,0.3)',
    '--today-bg': 'rgba(76,175,80,0.18)',
    '--weekend': 'rgba(232,140,140,0.85)',
    '--text-primary': '#e8ffe8',
    '--text-secondary': 'rgba(232,255,232,0.6)',
    '--text-muted': 'rgba(232,255,232,0.35)',
    '--border': 'rgba(255,255,255,0.12)',
    '--accent': '#4caf50',
    '--accent-hover': '#66bb6a',
  },
}

function applySettings(s: Settings) {
  const root = document.documentElement
  const vars = themeVars[s.theme] || themeVars.dark
  for (const [k, v] of Object.entries(vars)) {
    root.style.setProperty(k, v)
  }
  root.style.setProperty('--app-opacity', String(s.opacity / 100))
  root.style.setProperty('--editor-font-size', s.fontSize + 'px')
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

function formatDate(d: Date): string {
  const y = d.getFullYear()
  const m = String(d.getMonth() + 1).padStart(2, '0')
  const day = String(d.getDate()).padStart(2, '0')
  return `${y}-${m}-${day}`
}

onMounted(async () => {
  await loadMonth()
  const s = await getSettings()
  currentSettings.value = s
  applySettings(s)
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
