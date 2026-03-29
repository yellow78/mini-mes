<template>
  <div class="dashboard">
    <!-- KPI 列 -->
    <div class="kpi-row">
      <div class="kpi-card">
        <div class="kpi-value">{{ overallUtilization }}%</div>
        <div class="kpi-label">整體稼動率</div>
      </div>
      <div class="kpi-card kpi-running">
        <div class="kpi-value">{{ statusCount.running }}</div>
        <div class="kpi-label">Running</div>
      </div>
      <div class="kpi-card kpi-idle">
        <div class="kpi-value">{{ statusCount.idle }}</div>
        <div class="kpi-label">Idle</div>
      </div>
      <div class="kpi-card kpi-down">
        <div class="kpi-value">{{ statusCount.down }}</div>
        <div class="kpi-label">Down</div>
      </div>
      <div class="kpi-card kpi-pm">
        <div class="kpi-value">{{ statusCount.pm }}</div>
        <div class="kpi-label">PM</div>
      </div>
      <div class="kpi-card kpi-alarm" :class="{ 'kpi-alarm--active': statusCount.alarm > 0 }">
        <div class="kpi-value">{{ statusCount.alarm }}</div>
        <div class="kpi-label">Alarm</div>
      </div>
    </div>

    <!-- 設備群組折疊清單 -->
    <div class="groups-area">
      <EquipmentGroupList />
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, onUnmounted } from 'vue'
import { useEquipmentStore } from '../stores/equipment'
import EquipmentGroupList from '../components/equipment/EquipmentGroupList.vue'

const equipmentStore    = useEquipmentStore()
const overallUtilization = computed(() => equipmentStore.overallUtilization)
const statusCount       = computed(() => equipmentStore.statusCount)

onMounted(() => {
  equipmentStore.fetchEquipments()
  equipmentStore.startSimulation()
})

onUnmounted(() => {
  equipmentStore.stopSimulation()
})
</script>

<style scoped>
.dashboard {
  display: flex;
  flex-direction: column;
  gap: 16px;
  height: 100%;
  padding: 16px;
  overflow: hidden;
}

/* KPI 列 */
.kpi-row {
  display: flex;
  gap: 12px;
  flex-shrink: 0;
}

.kpi-card {
  background: var(--mes-surface);
  border: 1px solid var(--mes-border);
  border-radius: 8px;
  padding: 12px 20px;
  text-align: center;
  min-width: 100px;
}

.kpi-value {
  font-size: 28px;
  font-weight: 700;
  color: var(--el-text-color-primary);
  line-height: 1.1;
}

.kpi-label {
  font-size: 11px;
  color: var(--el-text-color-regular);
  margin-top: 4px;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.kpi-running .kpi-value { color: var(--mes-running); }
.kpi-idle    .kpi-value { color: var(--mes-idle);    }
.kpi-down    .kpi-value { color: var(--mes-down);    }
.kpi-pm      .kpi-value { color: var(--mes-pm);      }
.kpi-alarm   .kpi-value { color: var(--el-text-color-regular); }

.kpi-alarm--active .kpi-value {
  color: var(--mes-alarm);
  animation: alarm-pulse 1.5s ease-in-out infinite;
}

.groups-area {
  flex: 1;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

@keyframes alarm-pulse {
  0%, 100% { opacity: 1; }
  50%       { opacity: 0.5; }
}
</style>
