<template>
  <div class="calendar">
    <div class="grid" :style="{ gridTemplateRows: `repeat(${rowCount}, 1fr)` }">
      <DayCell
        v-for="(day, i) in days"
        :key="i"
        :day="day"
        :selected="day.solarDate === selectedDate"
        :show-lunar="showLunar"
        :show-solar-term="showSolarTerm"
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
  showLunar: boolean
  showSolarTerm: boolean
}>()

defineEmits<{
  selectDay: [date: string]
  quickAdd: [date: string]
}>()

const rowCount = computed(() => Math.ceil(props.days.length / 7))
</script>

<style scoped>
.calendar {
  flex: 1;
  display: flex;
  flex-direction: column;
  min-width: 0;
}

.grid {
  display: grid;
  grid-template-columns: repeat(7, 1fr);
  gap: 0;
  padding: 0;
  flex: 1;
  min-height: 0;
}
</style>
