<template>
  <div
    class="day-cell"
    :class="{
      today: day.isToday,
      selected,
      'other-month': !day.isCurrent,
    }"
    :style="day.bgColor ? { background: day.bgColor + '18', borderColor: day.bgColor + '40' } : {}"
    @click="$emit('select', day.solarDate)"
    @dblclick="$emit('quickAdd', day.solarDate)"
  >
    <div class="day-number-row">
      <span class="day-number">{{ dayNum }}</span>
      <span v-if="day.isToday" class="today-dot"></span>
    </div>
    <div v-if="showLunar || showSolarTerm" class="lunar-row">
      <span v-if="showSolarTerm && day.solarTerm" class="solar-term">{{ day.solarTerm }}</span>
      <span v-else-if="showLunar" class="lunar-text">{{ lunarDisplay }}</span>
    </div>
    <div v-if="day.events.length > 0" class="events-text">
      <div
        v-for="evt in day.events.slice(0, 3)"
        :key="evt.id"
        class="event-text"
        :class="{ completed: evt.completed }"
      >{{ evt.title }}</div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { DayInfo } from '../types'

const props = defineProps<{
  day: DayInfo
  selected: boolean
  showLunar: boolean
  showSolarTerm: boolean
}>()

const dayNum = computed(() => {
  const d = new Date(props.day.solarDate)
  return d.getDate()
})

const lunarDisplay = computed(() => {
  if (!props.day.lunar) return ''
  if (props.day.lunar.day === 1) {
    return props.day.lunar.monthName
  }
  return props.day.lunar.dayName
})
</script>

<style scoped>
.day-cell {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 6px 4px;
  border-radius: var(--radius-sm);
  cursor: pointer;
  transition: all 0.15s;
  position: relative;
  min-height: 0;
  gap: 2px;
  border: 1px solid var(--border);
}

.day-cell:hover {
  background: var(--bg-hover);
}

.day-cell.selected {
  background: var(--bg-active);
  border-color: var(--accent);
}

.day-cell.today {
  background: var(--today-bg);
  border-color: var(--today-border);
}

.day-cell.today.selected {
  background: var(--bg-active);
  border-color: var(--accent);
}

.day-cell.other-month {
  opacity: 0.35;
}

.day-number-row {
  display: flex;
  align-items: center;
  gap: 4px;
  position: relative;
}

.day-number {
  font-size: 14px;
  font-weight: 500;
  line-height: 1;
}

.today .day-number {
  color: var(--accent);
  font-weight: 700;
}

.today-dot {
  width: 4px;
  height: 4px;
  border-radius: 50%;
  background: var(--accent);
}

.lunar-row {
  font-size: 10px;
  line-height: 1.2;
  text-align: center;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  max-width: 100%;
}

.lunar-text {
  color: var(--text-secondary);
}

.solar-term {
  color: #e8a0a0;
  font-weight: 500;
}

.events-text {
  width: 100%;
  margin-top: 2px;
  overflow: hidden;
  flex: 1;
  min-height: 0;
}

.event-text {
  font-size: 10px;
  line-height: 1.3;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  color: var(--text-primary);
  text-align: left;
  padding: 0 2px;
  border-radius: 2px;
}

.event-text.completed {
  text-decoration: line-through;
  color: var(--text-secondary);
}
</style>
