<template>
  <header class="header">
    <div class="header-left">
      <h1 class="title">桌面日历</h1>
      <span class="date-display">{{ year }}年{{ month }}月</span>
    </div>
    <div class="header-right">
      <button class="btn btn-today" @click="$emit('today')">今天</button>
      <div class="data-group">
        <button class="btn btn-data" @click="$emit('export')" title="导出数据">
          <svg width="14" height="14" viewBox="0 0 14 14" fill="currentColor">
            <path d="M7 1a.5.5 0 0 1 .5.5v7.793l2.146-2.147a.5.5 0 0 1 .708.708l-3 3a.5.5 0 0 1-.708 0l-3-3a.5.5 0 1 1 .708-.708L6.5 9.293V1.5A.5.5 0 0 1 7 1zM2 11.5a.5.5 0 0 1 .5.5v1h9v-1a.5.5 0 0 1 1 0v1.5a.5.5 0 0 1-.5.5h-10a.5.5 0 0 1-.5-.5V12a.5.5 0 0 1 .5-.5z"/>
          </svg>
        </button>
        <button class="btn btn-data" @click="$emit('import')" title="导入数据">
          <svg width="14" height="14" viewBox="0 0 14 14" fill="currentColor">
            <path d="M7 10a.5.5 0 0 1-.5-.5V1.707L4.354 3.854a.5.5 0 1 1-.708-.708l3-3a.5.5 0 0 1 .708 0l3 3a.5.5 0 0 1-.708.708L7.5 1.707V9.5a.5.5 0 0 1-.5.5zM2 11.5a.5.5 0 0 1 .5.5v1h9v-1a.5.5 0 0 1 1 0v1.5a.5.5 0 0 1-.5.5h-10a.5.5 0 0 1-.5-.5V12a.5.5 0 0 1 .5-.5z"/>
          </svg>
        </button>
      </div>
      <div class="nav-group">
        <button class="btn btn-nav" @click="$emit('prev')">
          <svg width="16" height="16" viewBox="0 0 16 16" fill="currentColor">
            <path d="M10.354 3.354a.5.5 0 0 0-.708-.708l-5 5a.5.5 0 0 0 0 .708l5 5a.5.5 0 0 0 .708-.708L5.707 8l4.647-4.646z"/>
          </svg>
        </button>
        <button class="btn btn-nav" @click="$emit('next')">
          <svg width="16" height="16" viewBox="0 0 16 16" fill="currentColor">
            <path d="M5.646 3.354a.5.5 0 0 1 .708-.708l5 5a.5.5 0 0 1 0 .708l-5 5a.5.5 0 0 1-.708-.708L10.293 8 5.646 3.354z"/>
          </svg>
        </button>
      </div>
    </div>
  </header>
</template>

<script setup lang="ts">
defineProps<{
  year: number
  month: number
}>()

defineEmits<{
  prev: []
  next: []
  today: []
  export: []
  import: []
}>()
</script>

<style scoped>
.header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 8px 4px;
  flex-shrink: 0;
}

.header-left {
  display: flex;
  align-items: baseline;
  gap: 16px;
}

.title {
  font-size: 20px;
  font-weight: 600;
  letter-spacing: 0.5px;
  background: linear-gradient(135deg, #fff, rgba(255,255,255,0.7));
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
}

.date-display {
  font-size: 15px;
  color: var(--text-secondary);
  font-weight: 400;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 12px;
}

.btn {
  padding: 6px 14px;
  border-radius: var(--radius-sm);
  font-size: 13px;
  transition: all 0.2s;
}

.btn-today {
  background: var(--bg-secondary);
  border: 1px solid var(--border);
}

.btn-today:hover {
  background: var(--bg-hover);
  border-color: var(--accent);
  color: var(--accent);
}

.data-group {
  display: flex;
  gap: 2px;
  background: var(--bg-secondary);
  border-radius: var(--radius-sm);
  border: 1px solid var(--border);
  overflow: hidden;
}

.btn-data {
  padding: 6px 8px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.btn-data:hover {
  background: var(--bg-hover);
  color: var(--accent);
}

.nav-group {
  display: flex;
  gap: 2px;
  background: var(--bg-secondary);
  border-radius: var(--radius-sm);
  border: 1px solid var(--border);
  overflow: hidden;
}

.btn-nav {
  padding: 6px 10px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.btn-nav:hover {
  background: var(--bg-hover);
}
</style>
