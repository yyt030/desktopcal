<template>
  <aside class="panel">
    <div class="panel-header">
      <div class="panel-date">
        <span class="panel-day">{{ displayDay }}</span>
        <span class="panel-weekday">{{ weekdayName }}</span>
      </div>
      <div v-if="lunar" class="panel-lunar">
        <span>{{ lunar.monthName }}{{ lunar.dayName }}</span>
        <span class="lunar-extra">{{ lunar.ganZhi }}年 {{ lunar.animal }}</span>
      </div>
      <div v-if="solarTerm" class="panel-solar-term">
        {{ solarTerm }}
      </div>
    </div>

    <div class="events-section">
      <div class="section-header">
        <span class="section-title">日程</span>
        <button class="btn-add" @click="showForm = true">
          <svg width="14" height="14" viewBox="0 0 14 14" fill="currentColor">
            <path d="M7 1a.75.75 0 0 1 .75.75v4.5h4.5a.75.75 0 0 1 0 1.5h-4.5v4.5a.75.75 0 0 1-1.5 0v-4.5h-4.5a.75.75 0 0 1 0-1.5h4.5v-4.5A.75.75 0 0 1 7 1z"/>
          </svg>
        </button>
      </div>

      <div v-if="showForm" class="event-form">
        <input
          v-model="form.title"
          class="form-input"
          placeholder="事件标题"
          @keyup.enter="submitForm"
        />
        <div class="form-row">
          <input v-model="form.startTime" type="time" class="form-input form-time" />
          <span class="form-sep">至</span>
          <input v-model="form.endTime" type="time" class="form-input form-time" />
        </div>
        <textarea
          v-model="form.description"
          class="form-input form-desc"
          placeholder="描述（可选）"
          rows="2"
        ></textarea>
        <div class="form-colors">
          <button
            v-for="c in colors"
            :key="c"
            class="color-btn"
            :class="{ active: form.color === c }"
            :style="{ background: c }"
            @click="form.color = c"
          />
        </div>
        <div class="form-actions">
          <button class="btn-form btn-cancel" @click="cancelForm">取消</button>
          <button class="btn-form btn-save" @click="submitForm">保存</button>
        </div>
      </div>

      <div v-if="!showForm && events.length === 0" class="no-events">
        暂无日程
      </div>

      <TransitionGroup name="list" tag="div" class="events-list">
        <div
          v-for="evt in events"
          :key="evt.id"
          class="event-card"
          :class="{ editing: editingId === evt.id, completed: evt.completed }"
        >
          <div v-if="editingId !== evt.id" class="event-view" @click="startEdit(evt)" @contextmenu.prevent="toggleCompleted(evt)">
            <div class="event-color-bar" :style="{ background: evt.color }"></div>
            <div class="event-info">
              <div class="event-title">{{ evt.title }}</div>
              <div v-if="evt.startTime" class="event-time">
                {{ evt.startTime }}<span v-if="evt.endTime"> - {{ evt.endTime }}</span>
              </div>
              <div v-if="evt.description" class="event-desc">{{ evt.description }}</div>
            </div>
            <span v-if="evt.completed" class="completed-badge">✓</span>
            <button class="btn-delete" @click.stop="confirmDelete(evt.id)">
              <svg width="14" height="14" viewBox="0 0 14 14" fill="currentColor">
                <path d="M5.5 5.5a.5.5 0 0 1 .5.5v3a.5.5 0 0 1-1 0V6a.5.5 0 0 1 .5-.5zm3 0a.5.5 0 0 1 .5.5v3a.5.5 0 0 1-1 0V6a.5.5 0 0 1 .5-.5z"/>
                <path fill-rule="evenodd" d="M4.5 2A1.5 1.5 0 0 1 6 .5h2A1.5 1.5 0 0 1 9.5 2h3a.5.5 0 0 1 0 1h-.538l-.853 10.66A2 2 0 0 1 9.115 15H4.885a2 2 0 0 1-1.994-1.84L2.038 3H1.5a.5.5 0 0 1 0-1h3zM6 1.5A.5.5 0 0 0 5.5 2h3a.5.5 0 0 0-.5-.5h-2zM3.05 3l.845 10.538a1 1 0 0 0 .997.962h4.216a1 1 0 0 0 .997-.962L10.95 3h-7.9z"/>
              </svg>
            </button>
          </div>
          <div v-else class="event-edit">
            <input v-model="editForm.title" class="form-input" placeholder="标题" />
            <div class="form-row">
              <input v-model="editForm.startTime" type="time" class="form-input form-time" />
              <span class="form-sep">至</span>
              <input v-model="editForm.endTime" type="time" class="form-input form-time" />
            </div>
            <textarea
              v-model="editForm.description"
              class="form-input form-desc"
              rows="2"
            ></textarea>
            <div class="form-colors">
              <button
                v-for="c in colors"
                :key="c"
                class="color-btn"
                :class="{ active: editForm.color === c }"
                :style="{ background: c }"
                @click="editForm.color = c"
              />
            </div>
            <div class="form-actions">
              <button class="btn-form btn-cancel" @click="editingId = ''">取消</button>
              <button class="btn-form btn-save" @click="submitEdit">保存</button>
            </div>
          </div>
        </div>
      </TransitionGroup>
    </div>
  </aside>
</template>

<script setup lang="ts">
import { ref, computed, watch, reactive } from 'vue'
import type { Event, LunarDay } from '../types'

const props = defineProps<{
  date: string
  events: Event[]
  lunar: LunarDay | undefined
  solarTerm: string
  quickAddDate: string
}>()

const emit = defineEmits<{
  create: [event: Omit<Event, 'id' | 'createdAt' | 'updatedAt'>]
  update: [event: Event]
  delete: [id: string]
}>()

const colors = ['#4A90D9', '#E87878', '#6BCB77', '#F0A050', '#9B72CF', '#4ECDC4']
const weekdays = ['周日', '周一', '周二', '周三', '周四', '周五', '周六']

const showForm = ref(false)
const editingId = ref('')

const form = reactive({
  title: '',
  description: '',
  startTime: '',
  endTime: '',
  color: '#4A90D9',
})

const editForm = reactive({
  title: '',
  description: '',
  startTime: '',
  endTime: '',
  color: '#4A90D9',
})

const displayDay = computed(() => {
  if (!props.date) return ''
  const d = new Date(props.date)
  return d.getDate()
})

const weekdayName = computed(() => {
  if (!props.date) return ''
  const d = new Date(props.date)
  return weekdays[d.getDay()]
})

watch(() => props.date, () => {
  showForm.value = false
  editingId.value = ''
})

watch(() => props.quickAddDate, (val) => {
  if (val) {
    showForm.value = true
    resetForm()
  }
})

function toggleCompleted(evt: Event) {
  emit('update', { ...evt, completed: !evt.completed })
}

function submitForm() {
  if (!form.title.trim()) return
  emit('create', {
    title: form.title.trim(),
    description: form.description.trim(),
    date: props.date,
    startTime: form.startTime,
    endTime: form.endTime,
    color: form.color,
    completed: false,
  })
  resetForm()
}

function cancelForm() {
  resetForm()
}

function resetForm() {
  form.title = ''
  form.description = ''
  form.startTime = ''
  form.endTime = ''
  form.color = '#4A90D9'
  showForm.value = false
}

function startEdit(evt: Event) {
  editingId.value = evt.id
  editForm.title = evt.title
  editForm.description = evt.description
  editForm.startTime = evt.startTime
  editForm.endTime = evt.endTime
  editForm.color = evt.color
}

function submitEdit() {
  if (!editForm.title.trim()) return
  const original = props.events.find(e => e.id === editingId.value)
  if (!original) return
  emit('update', {
    ...original,
    title: editForm.title.trim(),
    description: editForm.description.trim(),
    startTime: editForm.startTime,
    endTime: editForm.endTime,
    color: editForm.color,
  })
  editingId.value = ''
}

function confirmDelete(id: string) {
  emit('delete', id)
}
</script>

<style scoped>
.panel {
  width: 300px;
  flex-shrink: 0;
  display: flex;
  flex-direction: column;
  background: var(--bg-secondary);
  border-radius: var(--radius);
  border: 1px solid var(--border);
  overflow: hidden;
}

.panel-header {
  padding: 20px 16px 16px;
  border-bottom: 1px solid var(--border);
}

.panel-date {
  display: flex;
  align-items: baseline;
  gap: 10px;
}

.panel-day {
  font-size: 42px;
  font-weight: 700;
  line-height: 1;
  background: linear-gradient(135deg, #fff, rgba(255,255,255,0.7));
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
}

.panel-weekday {
  font-size: 15px;
  color: var(--text-secondary);
}

.panel-lunar {
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

.panel-solar-term {
  margin-top: 6px;
  font-size: 12px;
  color: #e8a0a0;
  font-weight: 500;
}

.events-section {
  flex: 1;
  display: flex;
  flex-direction: column;
  padding: 12px;
  overflow-y: auto;
  min-height: 0;
}

.section-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 10px;
}

.section-title {
  font-size: 13px;
  font-weight: 600;
  color: var(--text-secondary);
  text-transform: uppercase;
  letter-spacing: 1px;
}

.btn-add {
  width: 28px;
  height: 28px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: var(--radius-sm);
  background: var(--accent);
  color: #fff;
  transition: all 0.2s;
}

.btn-add:hover {
  background: var(--accent-hover);
  transform: scale(1.05);
}

.event-form, .event-edit {
  background: rgba(0, 0, 0, 0.2);
  border-radius: var(--radius-sm);
  padding: 12px;
  margin-bottom: 10px;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.form-input {
  width: 100%;
  padding: 8px 10px;
  background: rgba(255, 255, 255, 0.06);
  border: 1px solid var(--border);
  border-radius: var(--radius-sm);
  color: var(--text-primary);
  font-size: 13px;
  transition: border-color 0.2s;
}

.form-input:focus {
  border-color: var(--accent);
}

.form-input::placeholder {
  color: var(--text-muted);
}

.form-row {
  display: flex;
  align-items: center;
  gap: 6px;
}

.form-time {
  flex: 1;
}

.form-sep {
  font-size: 12px;
  color: var(--text-muted);
}

.form-desc {
  resize: none;
  min-height: 40px;
}

.form-colors {
  display: flex;
  gap: 6px;
}

.color-btn {
  width: 20px;
  height: 20px;
  border-radius: 50%;
  border: 2px solid transparent;
  transition: all 0.15s;
  cursor: pointer;
}

.color-btn.active {
  border-color: #fff;
  transform: scale(1.15);
}

.form-actions {
  display: flex;
  gap: 8px;
  justify-content: flex-end;
}

.btn-form {
  padding: 6px 14px;
  border-radius: var(--radius-sm);
  font-size: 12px;
  font-weight: 500;
}

.btn-cancel {
  color: var(--text-secondary);
}

.btn-cancel:hover {
  color: var(--text-primary);
}

.btn-save {
  background: var(--accent);
  color: #fff;
}

.btn-save:hover {
  background: var(--accent-hover);
}

.no-events {
  text-align: center;
  padding: 30px 0;
  color: var(--text-muted);
  font-size: 13px;
}

.events-list {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.event-card {
  border-radius: var(--radius-sm);
  overflow: hidden;
  transition: all 0.2s;
}

.event-view {
  display: flex;
  gap: 10px;
  padding: 10px;
  background: rgba(255, 255, 255, 0.03);
  border-radius: var(--radius-sm);
  cursor: pointer;
  transition: background 0.15s;
}

.event-view:hover {
  background: rgba(255, 255, 255, 0.06);
}

.event-color-bar {
  width: 3px;
  border-radius: 2px;
  flex-shrink: 0;
}

.event-info {
  flex: 1;
  min-width: 0;
}

.event-title {
  font-size: 13px;
  font-weight: 500;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.event-time {
  font-size: 11px;
  color: var(--text-secondary);
  margin-top: 2px;
}

.event-desc {
  font-size: 11px;
  color: var(--text-muted);
  margin-top: 2px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.btn-delete {
  opacity: 0;
  padding: 4px;
  border-radius: 4px;
  color: var(--text-muted);
  transition: all 0.15s;
  flex-shrink: 0;
}

.event-view:hover .btn-delete {
  opacity: 1;
}

.btn-delete:hover {
  color: #e87878;
  background: rgba(232, 120, 120, 0.1);
}

.list-enter-active, .list-leave-active {
  transition: all 0.25s ease;
}

.list-enter-from {
  opacity: 0;
  transform: translateX(-10px);
}

.list-leave-to {
  opacity: 0;
  transform: translateX(10px);
}

.event-card.completed .event-view {
  opacity: 0.5;
}

.event-card.completed .event-title {
  text-decoration: line-through;
}

.completed-badge {
  color: #6BCB77;
  font-size: 14px;
  font-weight: 700;
  flex-shrink: 0;
  padding: 0 2px;
}
</style>
