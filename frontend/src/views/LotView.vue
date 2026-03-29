<template>
  <div class="lot-view">
    <div class="lot-header">
      <h2>Lot / WIP 看板</h2>
      <div class="lot-stats">
        <span class="stat-item">在製品 WIP：<b>{{ wipCount }}</b></span>
        <span class="stat-item">排隊中：<b>{{ queuedCount }}</b></span>
      </div>
      <el-button type="primary" @click="dialogVisible = true">+ 建立 Lot</el-button>
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
            <!-- 排隊中的 Lot 顯示派工按鈕 -->
            <el-button
              v-if="lot.status === 'QUEUED'"
              size="small"
              type="success"
              plain
              class="dispatch-btn"
              :loading="dispatchingId === lot.id"
              @click.stop="handleDispatch(lot.id)"
            >
              派工
            </el-button>
          </div>
        </div>
      </div>
    </div>

    <!-- 建立 Lot 對話框 -->
    <el-dialog
      v-model="dialogVisible"
      title="建立新 Lot"
      width="400px"
    >
      <el-form :model="form" label-width="80px">
        <el-form-item label="Lot 編號">
          <el-input v-model="form.lot_number" placeholder="例：LOT-2024-001" />
        </el-form-item>
        <el-form-item label="產品">
          <el-input v-model="form.product" placeholder="例：DRAM-32G" />
        </el-form-item>
        <el-form-item label="Recipe ID">
          <el-input-number v-model="form.recipe_id" :min="1" />
        </el-form-item>
        <el-form-item label="優先度">
          <el-select v-model="form.priority">
            <el-option label="P1（最高）" :value="1" />
            <el-option label="P2" :value="2" />
            <el-option label="P3" :value="3" />
            <el-option label="P4" :value="4" />
            <el-option label="P5（最低）" :value="5" />
          </el-select>
        </el-form-item>
        <el-form-item label="Wafer 數">
          <el-input-number v-model="form.wafer_count" :min="1" :max="25" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="creating" @click="handleCreate">建立</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import { ElMessage } from 'element-plus'
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

// 建立 Lot
const dialogVisible = ref(false)
const creating = ref(false)
const form = reactive({
  lot_number: '',
  product: '',
  recipe_id: 1,
  priority: 3,
  wafer_count: 25,
})

async function handleCreate() {
  if (!form.lot_number || !form.product) {
    ElMessage.warning('請填寫 Lot 編號與產品')
    return
  }
  creating.value = true
  try {
    await lotStore.addLot({ ...form })
    ElMessage.success(`Lot ${form.lot_number} 建立成功`)
    dialogVisible.value = false
    form.lot_number = ''
    form.product = ''
  } catch (e: unknown) {
    ElMessage.error(e instanceof Error ? e.message : '建立失敗')
  } finally {
    creating.value = false
  }
}

// 派工
const dispatchingId = ref<number | null>(null)

async function handleDispatch(id: number) {
  dispatchingId.value = id
  try {
    await lotStore.dispatch(id)
    ElMessage.success('派工成功')
  } catch (e: unknown) {
    ElMessage.error(e instanceof Error ? e.message : '派工失敗')
  } finally {
    dispatchingId.value = null
  }
}

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

.dispatch-btn {
  margin-top: 8px;
  width: 100%;
}
</style>
