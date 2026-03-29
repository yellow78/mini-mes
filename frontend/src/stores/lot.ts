// Lot / WIP 狀態 store
import { defineStore } from 'pinia'
import { computed, ref } from 'vue'
import type { Lot, LotStatus } from '../types/mes'

const MOCK_LOTS: Lot[] = [
  { id: 101, lotNumber: 'LOT-2024-001', product: 'NAND-128G', recipeId: 1, status: 'RUNNING',   priority: 1, waferCount: 25, createdAt: new Date(Date.now() - 4 * 3600 * 1000).toISOString() },
  { id: 102, lotNumber: 'LOT-2024-002', product: 'DRAM-8G',   recipeId: 2, status: 'RUNNING',   priority: 2, waferCount: 25, createdAt: new Date(Date.now() - 3 * 3600 * 1000).toISOString() },
  { id: 103, lotNumber: 'LOT-2024-003', product: 'NAND-128G', recipeId: 1, status: 'RUNNING',   priority: 1, waferCount: 25, createdAt: new Date(Date.now() - 2 * 3600 * 1000).toISOString() },
  { id: 201, lotNumber: 'LOT-2024-010', product: 'Logic-28nm',recipeId: 3, status: 'RUNNING',   priority: 2, waferCount: 25, createdAt: new Date(Date.now() - 5 * 3600 * 1000).toISOString() },
  { id: 202, lotNumber: 'LOT-2024-011', product: 'DRAM-8G',   recipeId: 4, status: 'RUNNING',   priority: 3, waferCount: 25, createdAt: new Date(Date.now() - 1 * 3600 * 1000).toISOString() },
  { id: 301, lotNumber: 'LOT-2024-020', product: 'NAND-256G', recipeId: 5, status: 'RUNNING',   priority: 1, waferCount: 25, createdAt: new Date(Date.now() - 6 * 3600 * 1000).toISOString() },
  { id: 401, lotNumber: 'LOT-2024-030', product: 'Logic-14nm',recipeId: 6, status: 'RUNNING',   priority: 1, waferCount: 25, createdAt: new Date(Date.now() - 8 * 3600 * 1000).toISOString() },
  { id: 501, lotNumber: 'LOT-2024-040', product: 'NAND-128G', recipeId: 1, status: 'QUEUED',    priority: 3, waferCount: 25, createdAt: new Date(Date.now() - 30 * 60 * 1000).toISOString() },
  { id: 502, lotNumber: 'LOT-2024-041', product: 'DRAM-16G',  recipeId: 2, status: 'QUEUED',    priority: 2, waferCount: 25, createdAt: new Date(Date.now() - 20 * 60 * 1000).toISOString() },
  { id: 601, lotNumber: 'LOT-2024-050', product: 'Logic-28nm',recipeId: 3, status: 'COMPLETED', priority: 2, waferCount: 25, createdAt: new Date(Date.now() - 24 * 3600 * 1000).toISOString() },
  { id: 602, lotNumber: 'LOT-2024-051', product: 'NAND-128G', recipeId: 1, status: 'ON_HOLD',   priority: 1, waferCount: 25, createdAt: new Date(Date.now() - 10 * 3600 * 1000).toISOString() },
]

export const useLotStore = defineStore('lot', () => {
  const lots = ref<Lot[]>(JSON.parse(JSON.stringify(MOCK_LOTS)))

  const wipCount    = computed(() => lots.value.filter(l => l.status === 'RUNNING').length)
  const queuedCount = computed(() => lots.value.filter(l => l.status === 'QUEUED').length)

  function updateLotStatus(id: number, status: LotStatus) {
    const lot = lots.value.find(l => l.id === id)
    if (lot) lot.status = status
  }

  async function fetchLots() {
    lots.value = JSON.parse(JSON.stringify(MOCK_LOTS))
  }

  return {
    lots,
    wipCount,
    queuedCount,
    updateLotStatus,
    fetchLots,
  }
})
