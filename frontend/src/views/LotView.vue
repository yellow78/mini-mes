<template>
  <div class="lot-view">
    <div class="lot-header">
      <h2>Lot / WIP 看板</h2>
      <div class="lot-stats">
        <span class="stat-item">在製品 WIP：<b>{{ wipCount }}</b></span>
        <span class="stat-item">排隊中：<b>{{ queuedCount }}</b></span>
      </div>
    </div>

    <!-- Lot 狀態分欄 -->
    <div class="lot-columns">
      <div
        v-for="col in columns"
        :key="col.status"
        class="lot-column"
      >
        <div class="col-header" :style="{ borderTopColor: col.color }">
          <span class="col-title">{{ col.label }}</span>
          <span class="col-count">{{ col.lots.length }}</span>
        </div>
        <div class="col-body">
          <div
            v-for="lot in col.lots"
            :key="lot.id"
            class="lot-card"
          >
            <div class="lot-number">{{ lot.lotNumber }}</div>
            <div class="lot-meta">
              <span>{{ lot.product }}</span>
              <span class="priority-badge" :class="'p' + lot.priority">P{{ lot.priority }}</span>
            </div>
            <div class="lot-wafer">{{ lot.waferCount }} wafers</div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useLotStore } from '../stores/lot'

const lotStore    = useLotStore()
const wipCount    = computed(() => lotStore.wipCount)
const queuedCount = computed(() => lotStore.queuedCount)

const columns = computed(() => [
  {
    status: 'QUEUED',
    label: '排隊中',
    color: 'var(--mes-idle)',
    lots: lotStore.lots.filter(l => l.status === 'QUEUED'),
  },
  {
    status: 'RUNNING',
    label: 'Running',
    color: 'var(--mes-running)',
    lots: lotStore.lots.filter(l => l.status === 'RUNNING'),
  },
  {
    status: 'ON_HOLD',
    label: 'On Hold',
    color: 'var(--mes-alarm)',
    lots: lotStore.lots.filter(l => l.status === 'ON_HOLD'),
  },
  {
    status: 'COMPLETED',
    label: '完成',
    color: 'var(--el-color-primary)',
    lots: lotStore.lots.filter(l => l.status === 'COMPLETED'),
  },
])

onMounted(() => {
  lotStore.fetchLots()
})
</script>

<style scoped>
.lot-view {
  padding: 16px;
  height: 100%;
  display: flex;
  flex-direction: column;
  gap: 16px;
  overflow: hidden;
}

.lot-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.lot-header h2 {
  font-size: 18px;
  font-weight: 700;
  color: var(--el-text-color-primary);
}

.lot-stats {
  display: flex;
  gap: 24px;
  font-size: 14px;
  color: var(--el-text-color-regular);
}

.lot-stats b {
  color: var(--el-text-color-primary);
}

.lot-columns {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 12px;
  flex: 1;
  overflow: hidden;
}

.lot-column {
  display: flex;
  flex-direction: column;
  background: var(--mes-surface);
  border: 1px solid var(--mes-border);
  border-radius: 8px;
  overflow: hidden;
}

.col-header {
  padding: 10px 14px;
  border-top: 3px solid;
  display: flex;
  justify-content: space-between;
  align-items: center;
  background: var(--mes-surface-deep);
}

.col-title {
  font-weight: 600;
  font-size: 13px;
  color: var(--el-text-color-primary);
}

.col-count {
  font-size: 12px;
  color: var(--el-text-color-regular);
  background: var(--mes-surface);
  padding: 2px 8px;
  border-radius: 10px;
}

.col-body {
  flex: 1;
  overflow-y: auto;
  padding: 8px;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.lot-card {
  background: var(--mes-surface-deep);
  border: 1px solid var(--mes-border);
  border-radius: 6px;
  padding: 10px 12px;
}

.lot-number {
  font-family: monospace;
  font-weight: 600;
  color: var(--el-color-primary);
  font-size: 13px;
  margin-bottom: 4px;
}

.lot-meta {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 12px;
  color: var(--el-text-color-regular);
}

.priority-badge {
  font-size: 11px;
  font-weight: 700;
  padding: 1px 6px;
  border-radius: 4px;
  background: var(--mes-border);
  color: var(--el-text-color-primary);
}

.priority-badge.p1 { background: var(--mes-alarm); color: #fff; }
.priority-badge.p2 { background: var(--mes-pm);    color: #fff; }

.lot-wafer {
  font-size: 11px;
  color: var(--el-text-color-regular);
  margin-top: 4px;
}
</style>
