<template>
  <div class="calendar">
    <div class="weekday-row">
      <span
        v-for="(w, i) in weekdays"
        :key="i"
        class="weekday"
        :class="{ weekend: weekendIndicesces.has(i) }"
      >{{ w }}</span>
    </div>
    <div class="grid" :style="{ gridTemplateRows: `repeat(${rowCount}, 1fr)` }">
      <DayCell
        v-for="(day, i) in days"
        :key="i"
        :day="day"
        :selected="day.solarDate === selectedDate"
        @select="$emit('selectDay', $event)"
        @quick-add="$emit('quickAdd', $event)"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import DayCell from './DayCell.vue'
import type { DayInfo } from '../types'

const props = defineProps<{
  days: DayInfo[]
  selectedDate: string
  firstDayOfWeek: number
}>()

defineEmits<{
  selectDay: [date: string]
  quickAdd: [date: string]
}>()

const allWeekdays = ['日', '一', '二', '三', '四', '五', '六']
const weekendIndicesces = computed(() => {
  const indices = new Set<number>()
  for (let i = 0; i < 7; i++) {
    const dayOfWeek = (props.firstDayOfWeek + i) % 7
    if (dayOfWeek === 0 || dayOfWeek === 6) {
      indices.add(i)
    }
  }
  return indices
})

const weekdays = computed(() => {
  const result: string[] = []
  for (let i = 0; i < 7; i++) {
    result.push(allWeekdays[(props.firstDayOfWeek + i) % 7])
  }
  return result
})

const rowCount = computed(() => props.days.length / 7)
</script>

<style scoped>
.calendar {
  flex: 1;
  display: flex;
  flex-direction: column;
  min-width: 0;
}

.weekday-row {
  display: grid;
  grid-template-columns: repeat(7, 1fr);
  padding: 0 4px 8px;
  border-bottom: 1px solid var(--border);
  margin-bottom: 4px;
}

.weekday {
  text-align: center;
  font-size: 12px;
  color: var(--text-secondary);
  font-weight: 500;
  text-transform: uppercase;
  letter-spacing: 1px;
}

.weekday.weekend {
  color: var(--weekend);
}

.grid {
  display: grid;
  grid-template-columns: repeat(7, 1fr);
  gap: 2px;
  padding: 4px;
  flex: 1;
  min-height: 0;
}
</style>
