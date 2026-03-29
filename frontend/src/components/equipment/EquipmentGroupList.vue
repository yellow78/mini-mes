<template>
  <!-- 群組折疊主容器：包含工具列 + 所有群組 -->
  <div class="eq-group-list">
    <!-- 工具列 -->
    <div class="toolbar">
      <!-- 搜尋框 -->
      <el-input
        v-model="keyword"
        placeholder="搜尋設備名稱或 Lot 編號..."
        clearable
        size="small"
        style="width: 240px"
        @input="applyFilter"
        @clear="applyFilter"
      >
        <template #prefix>
          <span style="font-size: 13px">🔍</span>
        </template>
      </el-input>

      <!-- 類型篩選 Pill -->
      <el-radio-group v-model="selectedType" size="small" @change="applyFilter">
        <el-radio-button label="ALL">全部</el-radio-button>
        <el-radio-button label="CVD">CVD</el-radio-button>
        <el-radio-button label="Etch">Etch</el-radio-button>
        <el-radio-button label="CMP">CMP</el-radio-button>
        <el-radio-button label="Diffusion">Diffusion</el-radio-button>
      </el-radio-group>

      <!-- 僅看 Alarm 快篩 -->
      <el-radio-group v-model="alarmOnly" size="small" @change="applyFilter">
        <el-radio-button :label="false">全部設備</el-radio-button>
        <el-radio-button :label="true" class="alarm-filter-btn">⚠ 僅看 Alarm</el-radio-button>
      </el-radio-group>
    </div>

    <!-- 群組列表 -->
    <div class="groups-container">
      <EquipmentGroup
        v-for="group in filteredGroups"
        :key="group.type"
        :group="group"
        :force-expand="shouldForceExpand(group)"
        @select-equipment="handleSelectEquipment"
      />

      <div v-if="filteredGroups.length === 0" class="no-result">
        沒有符合條件的設備
      </div>
    </div>

    <!-- 設備詳細 Drawer -->
    <EquipmentDrawer
      v-model="drawerVisible"
      :equipment="selectedEquipment"
    />
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import { useEquipmentStore } from '../../stores/equipment'
import type { Equipment, EquipmentGroup as EquipmentGroupType, EquipmentType } from '../../types/mes'
import EquipmentGroup from './EquipmentGroup.vue'
import EquipmentDrawer from './EquipmentDrawer.vue'

const equipmentStore = useEquipmentStore()

const keyword      = ref('')
const selectedType = ref<EquipmentType | 'ALL'>('ALL')
const alarmOnly    = ref(false)

const filteredGroups = computed(() => equipmentStore.filteredGroups)

function applyFilter() {
  equipmentStore.setFilter(selectedType.value, alarmOnly.value, keyword.value)
}

// 搜尋關鍵字 or 僅看 Alarm 時，含 alarm 或符合搜尋結果的群組自動展開
function shouldForceExpand(group: EquipmentGroupType): boolean {
  if (alarmOnly.value && group.alarmCount > 0) return true
  if (keyword.value && group.equipments.length > 0) return true
  return false
}

// Drawer 相關
const drawerVisible    = ref(false)
const selectedEquipment = ref<Equipment | null>(null)

function handleSelectEquipment(eq: Equipment) {
  selectedEquipment.value = eq
  drawerVisible.value     = true
}
</script>

<style scoped>
.eq-group-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
  height: 100%;
}

.toolbar {
  display: flex;
  align-items: center;
  gap: 12px;
  flex-wrap: wrap;
  flex-shrink: 0;
}

.groups-container {
  flex: 1;
  overflow-y: auto;
  padding-right: 4px;
}

.no-result {
  text-align: center;
  color: var(--el-text-color-regular);
  padding: 40px;
  font-size: 14px;
}

/* 「僅看 Alarm」按鈕選中時變紅 */
.alarm-filter-btn.is-active {
  --el-radio-button-checked-bg-color: var(--mes-alarm);
  --el-radio-button-checked-border-color: var(--mes-alarm);
}
</style>
