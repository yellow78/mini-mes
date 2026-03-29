<template>
  <div class="alarm-list">
    <div v-if="alarms.length === 0" class="alarm-empty">
      目前無告警
    </div>
    <div
      v-for="alarm in alarms"
      :key="alarm.id"
      class="alarm-item"
      :class="{ acknowledged: alarm.acknowledged, critical: alarm.severity === 'CRITICAL' }"
    >
      <div class="alarm-header">
        <span class="alarm-equipment">{{ alarm.equipmentName }}</span>
        <el-tag :type="alarm.severity === 'CRITICAL' ? 'danger' : 'warning'" size="small">
          {{ alarm.severity }}
        </el-tag>
      </div>
      <div class="alarm-detail">
        {{ alarm.parameter === 'temperature' ? '溫度' : '壓力' }}
        超標：{{ alarm.value }} (UCL: {{ alarm.ucl }})
      </div>
      <div class="alarm-footer">
        <span class="alarm-time">{{ formatTime(alarm.timestamp) }}</span>
        <el-button
          v-if="!alarm.acknowledged"
          size="small"
          type="primary"
          plain
          @click="alarmStore.acknowledgeAlarm(alarm.id)"
        >
          確認
        </el-button>
        <span v-else class="alarm-acked">已確認</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useAlarmStore } from '../../stores/alarm'

const alarmStore = useAlarmStore()
const alarms     = computed(() => alarmStore.alarms)

function formatTime(iso: string): string {
  const d = new Date(iso)
  return d.toLocaleString('zh-TW', { month: '2-digit', day: '2-digit', hour: '2-digit', minute: '2-digit', second: '2-digit' })
}
</script>

<style scoped>
.alarm-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.alarm-empty {
  color: var(--el-text-color-regular);
  text-align: center;
  padding: 24px;
}

.alarm-item {
  background: var(--mes-surface);
  border: 1px solid var(--mes-border);
  border-left: 3px solid var(--mes-pm);
  border-radius: 6px;
  padding: 10px 12px;
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.alarm-item.critical {
  border-left-color: var(--mes-alarm);
}

.alarm-item.acknowledged {
  opacity: 0.5;
}

.alarm-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.alarm-equipment {
  font-family: monospace;
  font-weight: 600;
  color: var(--el-text-color-primary);
}

.alarm-detail {
  font-size: 13px;
  color: var(--el-text-color-regular);
}

.alarm-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.alarm-time {
  font-size: 12px;
  color: var(--el-text-color-regular);
}

.alarm-acked {
  font-size: 12px;
  color: var(--mes-idle);
}
</style>
