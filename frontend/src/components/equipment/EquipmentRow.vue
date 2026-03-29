<template>
  <!-- 單台設備列，alarm 設備有左側紅色邊框 -->
  <div
    class="eq-row"
    :class="{ 'eq-row--alarm': equipment.isAlarm }"
    @click="emit('select', equipment)"
  >
    <!-- 設備名稱 -->
    <div class="eq-col eq-name">{{ equipment.name }}</div>

    <!-- 狀態 Badge -->
    <div class="eq-col eq-status">
      <el-tag :type="statusTagType" size="small" effect="dark">
        {{ equipment.status }}
      </el-tag>
    </div>

    <!-- Current Lot -->
    <div class="eq-col eq-lot">
      <span v-if="equipment.currentLot" class="lot-number">{{ equipment.currentLot }}</span>
      <span v-else class="lot-empty">—</span>
    </div>

    <!-- 溫度 -->
    <div class="eq-col eq-temp" :class="{ 'param-alarm': isTempAlarm }">
      {{ equipment.temperature }}°C
      <span v-if="isTempAlarm" class="alarm-indicator">↑</span>
    </div>

    <!-- 壓力 -->
    <div class="eq-col eq-pressure" :class="{ 'param-alarm': isPressureAlarm }">
      {{ equipment.pressure }} mT
      <span v-if="isPressureAlarm" class="alarm-indicator">↑</span>
    </div>

    <!-- 稼動率 -->
    <div class="eq-col eq-util">
      <div class="util-bar-wrap">
        <div class="util-bar" :style="{ width: equipment.utilization + '%', background: utilizationColor }" />
      </div>
      <span class="util-pct">{{ equipment.utilization }}%</span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { Equipment } from '../../types/mes'

const props = defineProps<{ equipment: Equipment }>()
const emit  = defineEmits<{ (e: 'select', eq: Equipment): void }>()

const statusTagType = computed(() => {
  switch (props.equipment.status) {
    case 'RUNNING': return 'success'
    case 'IDLE':    return 'info'
    case 'DOWN':    return 'danger'
    case 'PM':      return 'warning'
  }
})

const isTempAlarm     = computed(() =>
  props.equipment.temperature > props.equipment.ucl_temp ||
  props.equipment.temperature < props.equipment.lcl_temp
)
const isPressureAlarm = computed(() =>
  props.equipment.pressure > props.equipment.ucl_pressure ||
  props.equipment.pressure < props.equipment.lcl_pressure
)

const utilizationColor = computed(() => {
  const u = props.equipment.utilization
  if (u >= 80) return 'var(--mes-running)'
  if (u >= 50) return 'var(--mes-pm)'
  return 'var(--mes-idle)'
})
</script>

<style scoped>
.eq-row {
  display: flex;
  align-items: center;
  padding: 8px 16px;
  border-bottom: 1px solid var(--mes-border);
  cursor: pointer;
  transition: background 0.15s;
  border-left: 2px solid transparent;
}

.eq-row:hover {
  background: rgba(59, 130, 246, 0.05);
}

.eq-row--alarm {
  border-left-color: var(--mes-alarm);
  background: rgba(239, 68, 68, 0.04);
}

.eq-col {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  font-size: 13px;
  color: var(--el-text-color-primary);
}

.eq-name {
  width: 80px;
  font-family: monospace;
  font-weight: 600;
}

.eq-status { width: 100px; }

.eq-lot {
  flex: 1;
  min-width: 0;
}

.eq-temp,
.eq-pressure {
  width: 90px;
  text-align: right;
}

.eq-util {
  width: 100px;
  display: flex;
  align-items: center;
  gap: 6px;
}

.lot-number {
  font-family: monospace;
  color: var(--el-color-primary);
}

.lot-empty {
  color: var(--el-text-color-regular);
}

.param-alarm {
  color: var(--mes-alarm) !important;
  font-weight: 600;
}

.alarm-indicator {
  font-size: 12px;
}

.util-bar-wrap {
  flex: 1;
  height: 6px;
  background: var(--mes-surface-deep);
  border-radius: 3px;
  overflow: hidden;
}

.util-bar {
  height: 100%;
  border-radius: 3px;
  transition: width 0.4s ease;
}

.util-pct {
  font-size: 12px;
  color: var(--el-text-color-regular);
  width: 32px;
  text-align: right;
}
</style>
