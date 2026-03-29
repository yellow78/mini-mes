// Lot / WIP 狀態 store
import { defineStore } from 'pinia'
import { computed, ref } from 'vue'
import { getLots, createLot, dispatchLot, type CreateLotPayload } from '../api/lot'
import type { Lot, LotStatus } from '../types/mes'

export const useLotStore = defineStore('lot', () => {
  const lots    = ref<Lot[]>([])
  const loading = ref(false)
  const error   = ref<string | null>(null)

  const wipCount    = computed(() => lots.value.filter(l => l.status === 'RUNNING').length)
  const queuedCount = computed(() => lots.value.filter(l => l.status === 'QUEUED').length)

  // WebSocket lot_dispatched 事件觸發
  function updateLotStatus(id: number, status: LotStatus) {
    const lot = lots.value.find(l => l.id === id)
    if (lot) lot.status = status
  }

  async function fetchLots() {
    loading.value = true
    error.value   = null
    try {
      lots.value = await getLots()
    } catch (e: unknown) {
      error.value = e instanceof Error ? e.message : '載入 Lot 失敗'
    } finally {
      loading.value = false
    }
  }

  async function addLot(payload: CreateLotPayload): Promise<Lot> {
    const newLot = await createLot(payload)
    lots.value.unshift(newLot)
    return newLot
  }

  async function dispatch(id: number): Promise<void> {
    await dispatchLot(id)
    // 派工成功後更新本地狀態
    updateLotStatus(id, 'RUNNING')
    // 重新 fetch 取得最新 equipment 指派資訊
    await fetchLots()
  }

  return {
    lots,
    loading,
    error,
    wipCount,
    queuedCount,
    updateLotStatus,
    fetchLots,
    addLot,
    dispatch,
  }
})
