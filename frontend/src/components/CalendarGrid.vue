<template>
  <div class="calendar" :style="{ fontSize: fontSize + 'px', fontFamily: fontFamily || 'inherit' }">
    <div class="weekday-row">
      <span
        v-for="(w, i) in weekdayLabels"
        :key="i"
        class="weekday"
        :class="{ weekend: w.isWeekend }"
      >{{ w.label }}</span>
    </div>
    <div ref="gridRef" class="grid" :style="{ gridTemplateRows: `repeat(${rowCount}, 1fr)` }">
      <DayCell
        v-for="(day, i) in days"
        :key="i"
        :ref="setDayCellRef"
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
import { computed, nextTick, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import DayCell from './DayCell.vue'
import type { DayInfo } from '../types'

const props = defineProps<{
  days: DayInfo[]
  selectedDate: string
  firstDayOfWeek: number
  showLunar: boolean
  showSolarTerm: boolean
  fontSize: number
  fontFamily: string
}>()

defineEmits<{
  selectDay: [date: string]
  quickAdd: [date: string]
}>()

const allWeekdays = ['日', '一', '二', '三', '四', '五', '六']

const weekdayLabels = computed(() =>
  Array.from({ length: 7 }, (_, i) => {
    const dow = (props.firstDayOfWeek + i) % 7
    return { label: '周' + allWeekdays[dow], isWeekend: dow === 0 || dow === 6 }
  })
)

const rowCount = computed(() => Math.ceil(props.days.length / 7))

const gridRef = ref<HTMLElement>()
const dayCellRefs = ref<(InstanceType<typeof DayCell> | null)[]>([])

function setDayCellRef(el: any, index: number) {
  dayCellRefs.value[index] = el
}

let observer: ResizeObserver | undefined
let measureRaf: number | undefined

function measureAll() {
  if (measureRaf) return
  measureRaf = requestAnimationFrame(() => {
    measureRaf = undefined
    dayCellRefs.value.forEach(cell => cell?.measure())
  })
}

onMounted(() => {
  observer = new ResizeObserver(measureAll)
  if (gridRef.value) observer.observe(gridRef.value)
  nextTick(measureAll)
})

onBeforeUnmount(() => observer?.disconnect())

watch(() => props.days, () => nextTick(measureAll))
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
  padding: 0 2px 4px;
  border-bottom: 1px solid var(--border);
  margin-bottom: 2px;
}

.weekday {
  text-align: center;
  font-size: 0.85em;
  color: var(--text-secondary);
  font-weight: 500;
}

.weekday.weekend {
  color: var(--weekend);
}

.grid {
  display: grid;
  grid-template-columns: repeat(7, 1fr);
  gap: 0;
  padding: 0;
  flex: 1;
  min-height: 0;
  height: 0;
  overflow: hidden;
}
</style>
