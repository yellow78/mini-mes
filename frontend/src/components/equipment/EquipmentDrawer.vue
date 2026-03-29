<template>
  <el-drawer
    v-model="visible"
    direction="rtl"
    size="360px"
    :title="equipment?.name ?? '設備詳情'"
    class="eq-drawer"
  >
    <template v-if="equipment">
      <!-- 基本資訊 -->
      <div class="drawer-section">
        <div class="drawer-row">
          <span class="drawer-label">類型</span>
          <span>{{ equipment.type }}</span>
        </div>
        <div class="drawer-row">
          <span class="drawer-label">狀態</span>
          <el-tag :type="statusTagType" effect="dark" size="default">{{ equipment.status }}</el-tag>
        </div>
        <div class="drawer-row">
          <span class="drawer-label">Current Lot</span>
          <span class="mono">{{ equipment.currentLot ?? '—' }}</span>
        </div>
        <div class="drawer-row">
          <span class="drawer-label">Recipe</span>
          <span class="mono">{{ equipment.recipeName ?? '—' }}</span>
        </div>
      </div>

      <!-- 製程參數卡片 -->
      <div class="drawer-section">
        <div class="section-title">製程參數</div>
        <div class="param-cards">
          <!-- 溫度 -->
          <div class="param-card" :class="{ 'param-alarm': isTempAlarm }">
            <div class="param-name">溫度</div>
            <div class="param-value">{{ equipment.temperature }} °C</div>
            <div class="param-limits">
              <span class="limit-ucl">UCL {{ equipment.ucl_temp }}</span>
              <span class="limit-lcl">LCL {{ equipment.lcl_temp }}</span>
            </div>
            <div class="param-progress-wrap">
              <div
                class="param-progress"
                :style="{ width: tempPercent + '%', background: isTempAlarm ? 'var(--mes-alarm)' : 'var(--el-color-primary)' }"
              />
            </div>
          </div>

          <!-- 壓力 -->
          <div class="param-card" :class="{ 'param-alarm': isPressureAlarm }">
            <div class="param-name">壓力</div>
            <div class="param-value">{{ equipment.pressure }} mTorr</div>
            <div class="param-limits">
              <span class="limit-ucl">UCL {{ equipment.ucl_pressure }}</span>
              <span class="limit-lcl">LCL {{ equipment.lcl_pressure }}</span>
            </div>
            <div class="param-progress-wrap">
              <div
                class="param-progress"
                :style="{ width: pressurePercent + '%', background: isPressureAlarm ? 'var(--mes-alarm)' : '#22c55e' }"
              />
            </div>
          </div>
        </div>
      </div>

      <!-- 稼動率 -->
      <div class="drawer-section">
        <div class="section-title">稼動率</div>
        <div class="util-wrap">
          <el-progress
            :percentage="equipment.utilization"
            :color="utilizationColor"
            :stroke-width="10"
          />
        </div>
      </div>

      <!-- SPC 迷你趨勢圖 -->
      <div class="drawer-section">
        <div class="section-title">SPC 趨勢（最近 20 點）</div>
        <SpcMiniChart :data="spcData" :parameter="'temperature'" />
      </div>

      <!-- 操作按鈕 -->
      <div class="drawer-actions">
        <el-button type="danger" plain @click="handleHold">
          Hold 設備
        </el-button>
        <el-button type="primary" plain @click="goToSpc">
          查看完整 SPC
        </el-button>
      </div>
    </template>
  </el-drawer>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import SpcMiniChart from '../spc/SpcMiniChart.vue'
import { holdEquipment, getSpcHistory } from '../../api/equipment'
import { useEquipmentStore } from '../../stores/equipment'
import type { Equipment, SpcDataPoint } from '../../types/mes'

const props = defineProps<{
  modelValue: boolean
  equipment: Equipment | null
}>()
const emit = defineEmits<{
  (e: 'update:modelValue', val: boolean): void
}>()

const router = useRouter()

const visible = computed({
  get: () => props.modelValue,
  set: (v) => emit('update:modelValue', v),
})

const statusTagType = computed(() => {
  switch (props.equipment?.status) {
    case 'RUNNING': return 'success'
    case 'IDLE':    return 'info'
    case 'DOWN':    return 'danger'
    case 'PM':      return 'warning'
    default:        return 'info'
  }
})

const isTempAlarm = computed(() => {
  const e = props.equipment
  if (!e) return false
  return e.temperature > e.ucl_temp || e.temperature < e.lcl_temp
})

const isPressureAlarm = computed(() => {
  const e = props.equipment
  if (!e) return false
  return e.pressure > e.ucl_pressure || e.pressure < e.lcl_pressure
})

// 溫度在 UCL 範圍內的百分比（用於進度條顯示）
const tempPercent = computed(() => {
  const e = props.equipment
  if (!e) return 0
  const range = e.ucl_temp - e.lcl_temp
  return Math.min(100, Math.max(0, ((e.temperature - e.lcl_temp) / range) * 100))
})

const pressurePercent = computed(() => {
  const e = props.equipment
  if (!e) return 0
  const range = e.ucl_pressure - e.lcl_pressure
  return Math.min(100, Math.max(0, ((e.pressure - e.lcl_pressure) / range) * 100))
})

const utilizationColor = computed(() => {
  const u = props.equipment?.utilization ?? 0
  if (u >= 80) return 'var(--mes-running)'
  if (u >= 50) return 'var(--mes-pm)'
  return 'var(--mes-idle)'
})

// 從 API 取得 SPC 歷史資料（最近 20 點）
const spcData = ref<SpcDataPoint[]>([])
watch(() => props.equipment, async (eq) => {
  if (!eq) { spcData.value = []; return }
  spcData.value = await getSpcHistory(eq.id, 20)
}, { immediate: true })

const equipmentStore = useEquipmentStore()

async function handleHold() {
  if (!props.equipment) return
  try {
    await holdEquipment(props.equipment.id)
    ElMessage.success(`${props.equipment.name} 已 Hold`)
    visible.value = false
    await equipmentStore.fetchEquipments()
  } catch (e: unknown) {
    ElMessage.error(e instanceof Error ? e.message : 'Hold 失敗')
  }
}

function goToSpc() {
  router.push('/spc')
  visible.value = false
}
</script>

<style scoped>
.drawer-section {
  padding: 12px 0;
  border-bottom: 1px solid var(--mes-border);
}

.drawer-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 6px 0;
  font-size: 14px;
  color: var(--el-text-color-primary);
}

.drawer-label {
  color: var(--el-text-color-regular);
  font-size: 13px;
}

.mono {
  font-family: monospace;
}

.section-title {
  font-size: 12px;
  color: var(--el-text-color-regular);
  text-transform: uppercase;
  letter-spacing: 0.5px;
  margin-bottom: 8px;
}

.param-cards {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 8px;
}

.param-card {
  background: var(--mes-surface-deep);
  border: 1px solid var(--mes-border);
  border-radius: 6px;
  padding: 10px;
}

.param-card.param-alarm {
  border-color: var(--mes-alarm);
}

.param-name {
  font-size: 11px;
  color: var(--el-text-color-regular);
  margin-bottom: 4px;
}

.param-value {
  font-size: 18px;
  font-weight: 700;
  color: var(--el-text-color-primary);
  margin-bottom: 4px;
}

.param-card.param-alarm .param-value {
  color: var(--mes-alarm);
}

.param-limits {
  display: flex;
  justify-content: space-between;
  font-size: 10px;
  margin-bottom: 6px;
}

.limit-ucl { color: var(--mes-alarm); }
.limit-lcl { color: var(--mes-pm); }

.param-progress-wrap {
  height: 4px;
  background: var(--mes-border);
  border-radius: 2px;
  overflow: hidden;
}

.param-progress {
  height: 100%;
  border-radius: 2px;
  transition: width 0.3s;
}

.util-wrap {
  padding: 4px 0;
}

.drawer-actions {
  display: flex;
  gap: 12px;
  padding-top: 16px;
}

.drawer-actions .el-button {
  flex: 1;
}
</style>
