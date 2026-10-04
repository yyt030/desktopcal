<template>
  <div
    v-if="day.isEmpty"
    class="day-cell empty"
  />
  <div
    v-else
    ref="cellRef"
    class="day-cell"
    :class="{
      today: day.isToday,
      selected,
      'other-month': !day.isCurrentMonth,
      weekend: isWeekend,
    }"
    :style="day.bgColor ? { background: day.bgColor + '18', borderColor: day.bgColor + '40' } : {}"
    :title="cellTitle"
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
        v-for="(evt, i) in day.events"
        :key="evt.id"
        class="event-text"
      >
        <span class="event-idx">{{ i + 1 }}.</span>
        <span class="event-title" :class="{ completed: evt.completed }">{{ evt.title }}</span>
      </div>
      <div v-if="hiddenCount > 0" class="more-badge">+{{ hiddenCount }}</div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, nextTick, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import type { DayInfo } from '../types'
import { parseLocalDate } from '../utils/date'

const props = defineProps<{
  day: DayInfo
  selected: boolean
  showLunar: boolean
  showSolarTerm: boolean
}>()

const dayNum = computed(() => parseLocalDate(props.day.solarDate).getDate())

const isWeekend = computed(() => {
  if (!props.day.solarDate) return false
  const dow = parseLocalDate(props.day.solarDate).getDay()
  return dow === 0 || dow === 6
})

const cellTitle = computed(() => {
  if (!props.day.solarDate) return ''
  const parts = [props.day.solarDate]
  if (props.day.lunar) parts.push(`${props.day.lunar.monthName}${props.day.lunar.dayName}`)
  if (props.day.solarTerm) parts.push(props.day.solarTerm)
  return parts.join(' · ')
})

const lunarDisplay = computed(() => {
  if (!props.day.lunar) return ''
  if (props.day.lunar.day === 1) {
    return props.day.lunar.monthName
  }
  return props.day.lunar.dayName
})

// 事件预览按格子高度撑满显示，放不下的裁剪并用 +N 提示剩余条数。
// 观察根节点（而非随事件数增删的 .events-text），保证格子在任意时机都能
// 被监听到尺寸变化并重新测量。
const cellRef = ref<HTMLElement>()
const fitCount = ref(props.day.events.length)
const hiddenCount = computed(() => Math.max(0, props.day.events.length - fitCount.value))

const MORE_BADGE_HEIGHT = 13

let observer: ResizeObserver | undefined

function measure() {
  const cell = cellRef.value
  const box = cell?.querySelector<HTMLElement>('.events-text')
  if (!box) {
    fitCount.value = 0
    return
  }
  const total = props.day.events.length
  const children = Array.from(box.children).filter(
    (c): c is HTMLElement => c.classList.contains('event-text')
  )
  if (total === 0 || children.length === 0) {
    fitCount.value = 0
    return
  }
  const countFit = (avail: number): number => {
    let n = 0
    for (const child of children) {
      if (child.offsetTop + child.offsetHeight <= avail) n++
      else break
    }
    return n
  }
  const avail = box.clientHeight
  let count = countFit(avail)
  if (count < total) {
    count = countFit(avail - MORE_BADGE_HEIGHT)
    if (count === 0) count = 1
  }
  fitCount.value = Math.min(count, total)
}

onMounted(() => {
  observer = new ResizeObserver(() => measure())
  if (cellRef.value) observer.observe(cellRef.value)
  nextTick(measure)
})

onBeforeUnmount(() => observer?.disconnect())

// 根节点会随 day.isEmpty 在 v-if/v-else 间切换，需要重新绑定观察目标
watch(cellRef, (el, oldEl) => {
  if (!observer) return
  if (oldEl) observer.unobserve(oldEl)
  if (el) observer.observe(el)
  nextTick(measure)
})

watch(
  () => props.day,
  () => {
    fitCount.value = props.day.events.length
    nextTick(measure)
  }
)
</script>

<style scoped>
.day-cell {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 4px 3px;
  cursor: pointer;
  transition: all 0.15s;
  position: relative;
  min-height: 0;
  gap: 1px;
  border-right: 1px solid var(--border);
  border-bottom: 1px solid var(--border);
}

.day-cell:hover {
  background: var(--bg-hover);
}

.day-cell.selected {
  background: var(--bg-active);
}

.day-cell.today {
  background: var(--today-bg);
  box-shadow: inset 0 0 0 1px var(--today-border, var(--accent));
}

.day-cell.today.selected {
  background: var(--bg-active);
}

.day-cell.other-month {
  opacity: 0.35;
}

.day-cell.empty {
  cursor: default;
}

.day-cell.empty:hover {
  background: transparent;
}

.day-number-row {
  display: flex;
  align-items: center;
  gap: 4px;
  position: relative;
}

.day-number {
  font-size: var(--editor-font-size, 13px);
  font-weight: 500;
  line-height: 1;
}

/* 周末日期数字着色，与表头 --weekend 保持一致 */
.day-cell.weekend .day-number {
  color: var(--weekend);
}

/* 放在 weekend 之后，保证「今天」始终使用强调色 */
.day-cell.today .day-number {
  color: var(--accent);
  font-weight: 700;
}

.day-cell.today.weekend .day-number {
  color: var(--accent);
}

.today-dot {
  width: 4px;
  height: 4px;
  border-radius: 50%;
  background: var(--accent);
}

.lunar-row {
  font-size: calc(var(--editor-font-size, 13px) * 0.7);
  line-height: 1.1;
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
  position: relative;
}

.more-badge {
  position: absolute;
  right: 2px;
  bottom: 0;
  font-size: calc(var(--editor-font-size, 13px) * 0.7);
  line-height: 1.2;
  color: var(--text-secondary);
  background: var(--bg-primary);
  padding: 0 4px;
  border-radius: 4px;
}

.event-text {
  font-size: calc(var(--editor-font-size, 13px) * 0.7);
  line-height: 1.2;
  word-break: break-word;
  color: var(--text-primary);
  text-align: left;
  padding: 1px 2px;
  border-radius: 2px;
  margin-bottom: 1px;
}

.event-idx {
  color: var(--text-muted);
  font-variant-numeric: tabular-nums;
  margin-right: 2px;
}

.event-title.completed {
  text-decoration: line-through;
  color: var(--text-secondary);
}
</style>
