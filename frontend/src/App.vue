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
    />
    <div class="app-body">
      <CalendarGrid
        :days="days"
        :selected-date="selectedDate"
        @select-day="onSelectDay"
        @quick-add="onQuickAdd"
      />
    </div>
    <EventPopup
      :visible="showPopup"
      :date="selectedDate"
      :events="selectedEvents"
      :lunar="selectedLunar"
      :solar-term="selectedSolarTerm"
      :quick-add-date="quickAddDate"
      :current-bg-color="selectedBgColor"
      @close="showPopup = false"
      @create="onCreateEvent"
      @update="onUpdateEvent"
      @delete="onDeleteEvent"
      @set-bg-color="onSetBgColor"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch, nextTick } from 'vue'
import AppHeader from './components/AppHeader.vue'
import CalendarGrid from './components/CalendarGrid.vue'
import EventPopup from './components/EventPopup.vue'
import { getMonthDays, createEvent, updateEvent, deleteEvent, setDayBackground, exportData, importData } from './api'
import type { DayInfo, Event } from './types'

const now = new Date()
const currentYear = ref(now.getFullYear())
const currentMonth = ref(now.getMonth() + 1)
const days = ref<DayInfo[]>([])
const selectedDate = ref(formatDate(now))
const quickAddDate = ref('')
const showPopup = ref(false)

const selectedDay = computed(() =>
  days.value.find(d => d.solarDate === selectedDate.value)
)
const selectedEvents = computed(() => selectedDay.value?.events ?? [])
const selectedLunar = computed(() => selectedDay.value?.lunar)
const selectedSolarTerm = computed(() => selectedDay.value?.solarTerm ?? '')
const selectedBgColor = computed(() => selectedDay.value?.bgColor ?? '')

async function loadMonth() {
  days.value = await getMonthDays(currentYear.value, currentMonth.value)
}

function prevMonth() {
  if (currentMonth.value === 1) {
    currentMonth.value = 12
    currentYear.value--
  } else {
    currentMonth.value--
  }
}

function nextMonth() {
  if (currentMonth.value === 12) {
    currentMonth.value = 1
    currentYear.value++
  } else {
    currentMonth.value++
  }
}

function goToday() {
  const t = new Date()
  currentYear.value = t.getFullYear()
  currentMonth.value = t.getMonth() + 1
  selectedDate.value = formatDate(t)
}

function onSelectDay(date: string) {
  selectedDate.value = date
}

function onQuickAdd(date: string) {
  selectedDate.value = date
  quickAddDate.value = ''
  showPopup.value = true
  nextTick(() => {
    quickAddDate.value = date
  })
}

async function onCreateEvent(event: Omit<Event, 'id' | 'createdAt' | 'updatedAt'>) {
  const created = await createEvent({
    ...event,
    id: '',
    completed: false,
    createdAt: '',
    updatedAt: '',
  })
  await loadMonth()
  selectedDate.value = created.date
}

async function onUpdateEvent(event: Event) {
  await updateEvent(event)
  await loadMonth()
}

async function onDeleteEvent(id: string) {
  await deleteEvent(id)
  await loadMonth()
}

async function onSetBgColor(color: string) {
  await setDayBackground(selectedDate.value, color)
  await loadMonth()
}

async function onExport() {
  const data = await exportData()
  const blob = new Blob([data], { type: 'application/json' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = `desktopcal-backup-${formatDate(new Date())}.json`
  a.click()
  URL.revokeObjectURL(url)
}

async function onImport() {
  const input = document.createElement('input')
  input.type = 'file'
  input.accept = '.json'
  input.onchange = async () => {
    const file = input.files?.[0]
    if (!file) return
    const text = await file.text()
    await importData(text)
    await loadMonth()
  }
  input.click()
}

function formatDate(d: Date): string {
  const y = d.getFullYear()
  const m = String(d.getMonth() + 1).padStart(2, '0')
  const day = String(d.getDate()).padStart(2, '0')
  return `${y}-${m}-${day}`
}

watch([currentYear, currentMonth], loadMonth)

onMounted(async () => {
  await loadMonth()
})
</script>

<style scoped>
.app-container {
  display: flex;
  flex-direction: column;
  height: 100vh;
  padding: 16px 20px;
  gap: 12px;
}

.app-body {
  display: flex;
  flex: 1;
  min-height: 0;
}
</style>
