<template>
  <header class="app-header">
    <!-- Logo + 標題 -->
    <div class="header-logo">
      <span class="logo-icon">⬡</span>
      <span class="logo-text">Mini-MES</span>
    </div>

    <!-- 導覽列 -->
    <nav class="header-nav">
      <router-link to="/" class="nav-link" :class="{ active: route.path === '/' }">
        Dashboard
      </router-link>
      <router-link to="/lots" class="nav-link" :class="{ active: route.path === '/lots' }">
        Lot / WIP
      </router-link>
      <router-link to="/spc" class="nav-link" :class="{ active: route.path === '/spc' }">
        SPC
      </router-link>
    </nav>

    <div class="header-right">
      <!-- Alarm 快捷按鈕 -->
      <el-badge :value="unacknowledgedCount" :hidden="unacknowledgedCount === 0" type="danger">
        <el-button
          class="alarm-btn"
          :class="{ 'alarm-active': unacknowledgedCount > 0 }"
          @click="showAlarmDrawer = true"
        >
          ⚠ Alarm
        </el-button>
      </el-badge>

      <!-- 即時時鐘 -->
      <span class="realtime-clock">{{ currentTime }}</span>
    </div>

    <!-- Alarm 清單 Drawer -->
    <el-drawer
      v-model="showAlarmDrawer"
      title="告警清單"
      direction="rtl"
      size="400px"
    >
      <AlarmList />
    </el-drawer>
  </header>
</template>

<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref } from 'vue'
import { useRoute } from 'vue-router'
import { useAlarmStore } from '../../stores/alarm'
import AlarmList from '../alarm/AlarmList.vue'

const route             = useRoute()
const alarmStore        = useAlarmStore()
const showAlarmDrawer   = ref(false)
const currentTime       = ref('')

const unacknowledgedCount = computed(() => alarmStore.unacknowledgedCount)

let clockTimer: ReturnType<typeof setInterval> | null = null

function updateClock() {
  const now  = new Date()
  const date = now.toLocaleDateString('zh-TW', { year: 'numeric', month: '2-digit', day: '2-digit' })
  const time = now.toLocaleTimeString('zh-TW', { hour: '2-digit', minute: '2-digit', second: '2-digit' })
  currentTime.value = `${date} ${time}`
}

onMounted(() => {
  updateClock()
  clockTimer = setInterval(updateClock, 1000)
})

onUnmounted(() => {
  if (clockTimer) clearInterval(clockTimer)
})
</script>

<style scoped>
.app-header {
  display: flex;
  align-items: center;
  height: 56px;
  padding: 0 20px;
  background: var(--mes-surface-deep);
  border-bottom: 1px solid var(--mes-border);
  gap: 24px;
  flex-shrink: 0;
}

.header-logo {
  display: flex;
  align-items: center;
  gap: 8px;
  color: var(--el-color-primary);
  font-weight: 700;
  font-size: 18px;
  white-space: nowrap;
}

.logo-icon {
  font-size: 22px;
}

.header-nav {
  display: flex;
  gap: 4px;
  flex: 1;
}

.nav-link {
  padding: 6px 16px;
  border-radius: 6px;
  color: var(--el-text-color-regular);
  text-decoration: none;
  font-size: 14px;
  transition: all 0.2s;
}

.nav-link:hover,
.nav-link.active {
  background: rgba(59, 130, 246, 0.15);
  color: var(--el-color-primary);
}

.header-right {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-left: auto;
}

.alarm-btn {
  background: transparent;
  border-color: var(--mes-border);
  color: var(--el-text-color-regular);
}

.alarm-btn.alarm-active {
  border-color: var(--mes-alarm);
  color: var(--mes-alarm);
  animation: alarm-pulse 1.5s ease-in-out infinite;
}

.realtime-clock {
  font-size: 13px;
  color: var(--el-text-color-regular);
  font-variant-numeric: tabular-nums;
  white-space: nowrap;
}

@keyframes alarm-pulse {
  0%, 100% { opacity: 1; }
  50%       { opacity: 0.6; }
}
</style>
