<template>
  <!-- 單一設備群組，含折疊 Header + 展開後的設備列表 -->
  <div class="eq-group" :class="{ 'eq-group--alarm': group.alarmCount > 0, 'eq-group--expanded': isExpanded }">
    <!-- 群組 Header（點擊折疊/展開） -->
    <div class="group-header" @click="toggle">
      <span class="group-arrow">{{ isExpanded ? '▾' : '▸' }}</span>
      <span class="group-type">{{ group.type }}</span>
      <span class="group-count">{{ group.equipments.length }} 台</span>

      <!-- 各狀態小圓點 -->
      <div class="group-status-dots">
        <span v-if="group.statusCount.running" class="dot dot--running" :title="'Running: ' + group.statusCount.running">
          {{ group.statusCount.running }}
        </span>
        <span v-if="group.statusCount.idle" class="dot dot--idle" :title="'Idle: ' + group.statusCount.idle">
          {{ group.statusCount.idle }}
        </span>
        <span v-if="group.statusCount.down" class="dot dot--down" :title="'Down: ' + group.statusCount.down">
          {{ group.statusCount.down }}
        </span>
        <span v-if="group.statusCount.pm" class="dot dot--pm" :title="'PM: ' + group.statusCount.pm">
          {{ group.statusCount.pm }}
        </span>
      </div>

      <!-- 稼動率橫條 -->
      <div class="group-util-wrap">
        <div class="group-util-bar" :style="{ width: group.utilization + '%' }" />
      </div>
      <span class="group-util-pct">{{ group.utilization }}%</span>

      <!-- Alarm 標籤（有 Alarm 才顯示，附閃爍動畫） -->
      <span v-if="group.alarmCount > 0" class="group-alarm-tag">
        {{ group.alarmCount }} Alarm
      </span>
    </div>

    <!-- 展開後：表格 Header + 設備列 -->
    <transition name="slide">
      <div v-show="isExpanded" class="group-body">
        <!-- 表格欄位標題 -->
        <div class="eq-table-header">
          <span style="width: 80px">設備名</span>
          <span style="width: 100px">狀態</span>
          <span style="flex: 1">Current Lot</span>
          <span style="width: 90px; text-align: right">溫度</span>
          <span style="width: 90px; text-align: right">壓力</span>
          <span style="width: 100px">稼動率</span>
        </div>
        <EquipmentRow
          v-for="eq in group.equipments"
          :key="eq.id"
          :equipment="eq"
          @select="emit('selectEquipment', $event)"
        />
      </div>
    </transition>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import type { Equipment, EquipmentGroup } from '../../types/mes'
import EquipmentRow from './EquipmentRow.vue'

const props = defineProps<{
  group: EquipmentGroup
  forceExpand?: boolean   // 「僅看 Alarm」模式強制展開
}>()
const emit = defineEmits<{
  (e: 'selectEquipment', eq: Equipment): void
}>()

const isExpanded = ref(props.group.alarmCount > 0)

function toggle() {
  isExpanded.value = !isExpanded.value
}

// 外部強制展開（搜尋或 Alarm 篩選觸發）
watch(() => props.forceExpand, (v) => {
  if (v) isExpanded.value = true
})
</script>

<style scoped>
.eq-group {
  border: 1px solid var(--mes-border);
  border-radius: 8px;
  overflow: hidden;
  margin-bottom: 8px;
}

.eq-group--alarm {
  border-color: var(--mes-alarm);
}

.group-header {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 16px;
  background: var(--mes-surface);
  cursor: pointer;
  user-select: none;
  transition: background 0.15s;
}

.group-header:hover {
  background: rgba(59, 130, 246, 0.06);
}

.group-arrow {
  color: var(--el-text-color-regular);
  font-size: 12px;
  width: 14px;
}

.group-type {
  font-weight: 700;
  font-size: 14px;
  color: var(--el-text-color-primary);
  min-width: 80px;
}

.group-count {
  font-size: 12px;
  color: var(--el-text-color-regular);
  min-width: 40px;
}

.group-status-dots {
  display: flex;
  gap: 6px;
}

.dot {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 22px;
  height: 22px;
  border-radius: 50%;
  font-size: 11px;
  font-weight: 700;
  color: #fff;
}

.dot--running { background: var(--mes-running); }
.dot--idle    { background: var(--mes-idle);    }
.dot--down    { background: var(--mes-down);    }
.dot--pm      { background: var(--mes-pm);      }

.group-util-wrap {
  flex: 1;
  height: 6px;
  background: var(--mes-surface-deep);
  border-radius: 3px;
  overflow: hidden;
  max-width: 120px;
}

.group-util-bar {
  height: 100%;
  background: var(--el-color-primary);
  border-radius: 3px;
  transition: width 0.4s;
}

.group-util-pct {
  font-size: 13px;
  color: var(--el-text-color-regular);
  min-width: 36px;
  text-align: right;
}

.group-alarm-tag {
  background: var(--mes-alarm);
  color: #fff;
  font-size: 11px;
  font-weight: 700;
  padding: 2px 8px;
  border-radius: 10px;
  animation: blink 1.5s ease-in-out infinite;
  white-space: nowrap;
}

/* 表格欄位標題 */
.eq-table-header {
  display: flex;
  align-items: center;
  padding: 6px 16px;
  background: var(--mes-surface-deep);
  font-size: 11px;
  color: var(--el-text-color-regular);
  text-transform: uppercase;
  letter-spacing: 0.3px;
  border-bottom: 1px solid var(--mes-border);
  gap: 0;
}

.eq-table-header span {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

/* 折疊動畫 */
.slide-enter-active,
.slide-leave-active {
  transition: max-height 0.25s ease, opacity 0.2s;
  overflow: hidden;
  max-height: 1000px;
}

.slide-enter-from,
.slide-leave-to {
  max-height: 0;
  opacity: 0;
}

@keyframes blink {
  0%, 100% { opacity: 1; }
  50%       { opacity: 0.5; }
}
</style>
