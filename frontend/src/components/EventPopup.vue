<template>
  <Teleport to="body">
    <Transition name="popup">
      <div v-if="visible" class="popup-overlay" @click.self="emit('close')">
        <div class="popup">
          <div class="popup-header">
            <div class="popup-header-left">
              <div class="popup-date">
                <span class="popup-day">{{ displayDay }}</span>
                <span class="popup-weekday">{{ weekdayName }}</span>
              </div>
              <div v-if="lunar" class="popup-lunar">
                <span>{{ lunar.monthName }}{{ lunar.dayName }}</span>
                <span class="lunar-extra">{{ lunar.ganZhi }}年 {{ lunar.animal }}</span>
              </div>
              <div v-if="solarTerm" class="popup-solar-term">
                {{ solarTerm }}
              </div>
            </div>
            <div class="popup-header-right">
              <button class="btn-add" @click="showForm = true" title="添加事件">
                <svg width="14" height="14" viewBox="0 0 14 14" fill="currentColor">
                  <path d="M7 1a.75.75 0 0 1 .75.75v4.5h4.5a.75.75 0 0 1 0 1.5h-4.5v4.5a.75.75 0 0 1-1.5 0v-4.5h-4.5a.75.75 0 0 1 0-1.5h4.5v-4.5A.75.75 0 0 1 7 1z"/>
                </svg>
              </button>
              <button class="btn-close" @click="emit('close')" title="关闭">
                <svg width="14" height="14" viewBox="0 0 14 14" fill="currentColor">
                  <path d="M3.7 3.7a.5.5 0 0 1 .7 0L7 6.3l2.6-2.6a.5.5 0 0 1 .7.7L7.7 7l2.6 2.6a.5.5 0 0 1-.7.7L7 7.7l-2.6 2.6a.5.5 0 0 1-.7-.7L6.3 7 3.7 4.4a.5.5 0 0 1 0-.7z"/>
                </svg>
              </button>
            </div>
          </div>

          <div class="popup-body">
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
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import { ref, computed, watch, reactive } from 'vue'
import type { Event, LunarDay } from '../types'

const props = defineProps<{
  visible: boolean
  date: string
  events: Event[]
  lunar: LunarDay | undefined
  solarTerm: string
  quickAddDate: string
}>()

const emit = defineEmits<{
  close: []
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
.popup-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  backdrop-filter: blur(4px);
}

.popup {
  width: 380px;
  max-height: 70vh;
  display: flex;
  flex-direction: column;
  background: #1e1e2e;
  border-radius: 12px;
  border: 1px solid rgba(255, 255, 255, 0.1);
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.5), 0 0 0 1px rgba(255, 255, 255, 0.05);
  overflow: hidden;
}

.popup-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  padding: 20px 16px 16px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.08);
}

.popup-header-left {
  flex: 1;
}

.popup-header-right {
  display: flex;
  gap: 6px;
  flex-shrink: 0;
}

.popup-date {
  display: flex;
  align-items: baseline;
  gap: 10px;
}

.popup-day {
  font-size: 42px;
  font-weight: 700;
  line-height: 1;
  background: linear-gradient(135deg, #fff, rgba(255,255,255,0.7));
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
}

.popup-weekday {
  font-size: 15px;
  color: rgba(255, 255, 255, 0.6);
}

.popup-lunar {
  margin-top: 8px;
  font-size: 13px;
  color: rgba(255, 255, 255, 0.6);
  display: flex;
  gap: 8px;
  align-items: center;
}

.lunar-extra {
  font-size: 11px;
  color: rgba(255, 255, 255, 0.4);
}

.popup-solar-term {
  margin-top: 6px;
  font-size: 12px;
  color: #e8a0a0;
  font-weight: 500;
}

.btn-add, .btn-close {
  width: 28px;
  height: 28px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 6px;
  color: rgba(255, 255, 255, 0.6);
  transition: all 0.2s;
}

.btn-add {
  background: #4A90D9;
  color: #fff;
}

.btn-add:hover {
  background: #3a7bc8;
  transform: scale(1.05);
}

.btn-close {
  background: rgba(255, 255, 255, 0.06);
}

.btn-close:hover {
  background: rgba(255, 255, 255, 0.12);
  color: #fff;
}

.popup-body {
  flex: 1;
  padding: 12px;
  overflow-y: auto;
  min-height: 0;
}

.event-form, .event-edit {
  background: rgba(255, 255, 255, 0.04);
  border-radius: 8px;
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
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 6px;
  color: #fff;
  font-size: 13px;
  transition: border-color 0.2s;
}

.form-input:focus {
  border-color: #4A90D9;
  outline: none;
}

.form-input::placeholder {
  color: rgba(255, 255, 255, 0.3);
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
  color: rgba(255, 255, 255, 0.4);
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
  border-radius: 6px;
  font-size: 12px;
  font-weight: 500;
}

.btn-cancel {
  color: rgba(255, 255, 255, 0.6);
}

.btn-cancel:hover {
  color: #fff;
}

.btn-save {
  background: #4A90D9;
  color: #fff;
}

.btn-save:hover {
  background: #3a7bc8;
}

.no-events {
  text-align: center;
  padding: 30px 0;
  color: rgba(255, 255, 255, 0.4);
  font-size: 13px;
}

.events-list {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.event-card {
  border-radius: 8px;
  overflow: hidden;
  transition: all 0.2s;
}

.event-view {
  display: flex;
  gap: 10px;
  padding: 10px;
  background: rgba(255, 255, 255, 0.03);
  border-radius: 8px;
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
  color: #fff;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.event-time {
  font-size: 11px;
  color: rgba(255, 255, 255, 0.6);
  margin-top: 2px;
}

.event-desc {
  font-size: 11px;
  color: rgba(255, 255, 255, 0.4);
  margin-top: 2px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.btn-delete {
  opacity: 0;
  padding: 4px;
  border-radius: 4px;
  color: rgba(255, 255, 255, 0.4);
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

.popup-enter-active, .popup-leave-active {
  transition: all 0.25s ease;
}

.popup-enter-from, .popup-leave-to {
  opacity: 0;
}

.popup-enter-from .popup, .popup-leave-to .popup {
  transform: scale(0.95) translateY(10px);
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
</style>
