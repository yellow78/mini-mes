import { beforeEach, describe, expect, it, vi } from 'vitest'
import { createPinia, setActivePinia } from 'pinia'
import { useLotStore } from '../../stores/lot'
import type { Lot } from '../../types/mes'
import { getLots } from '../../api/lot'

vi.mock('../../api/lot', () => ({
  getLots: vi.fn(),
  createLot: vi.fn(),
  dispatchLot: vi.fn(),
}))

const mockLots: Lot[] = [
  { id: 101, lotNumber: 'LOT-2024-001', product: 'DRAM-256M', recipeId: 1, status: 'RUNNING',   priority: 1, waferCount: 25, createdAt: '2024-01-01T00:00:00Z' },
  { id: 102, lotNumber: 'LOT-2024-002', product: 'DRAM-256M', recipeId: 1, status: 'RUNNING',   priority: 2, waferCount: 25, createdAt: '2024-01-01T00:00:00Z' },
  { id: 103, lotNumber: 'LOT-2024-003', product: 'SRAM-512K', recipeId: 2, status: 'RUNNING',   priority: 3, waferCount: 20, createdAt: '2024-01-01T00:00:00Z' },
  { id: 104, lotNumber: 'LOT-2024-004', product: 'SRAM-512K', recipeId: 2, status: 'QUEUED',    priority: 4, waferCount: 20, createdAt: '2024-01-01T00:00:00Z' },
  { id: 105, lotNumber: 'LOT-2024-005', product: 'Flash-1G',  recipeId: 3, status: 'QUEUED',    priority: 5, waferCount: 30, createdAt: '2024-01-01T00:00:00Z' },
  { id: 106, lotNumber: 'LOT-2024-006', product: 'Flash-1G',  recipeId: 3, status: 'QUEUED',    priority: 1, waferCount: 15, createdAt: '2024-01-01T00:00:00Z' },
  { id: 107, lotNumber: 'LOT-2024-007', product: 'DRAM-256M', recipeId: 1, status: 'COMPLETED', priority: 2, waferCount: 25, createdAt: '2024-01-01T00:00:00Z' },
  { id: 108, lotNumber: 'LOT-2024-008', product: 'DRAM-256M', recipeId: 1, status: 'COMPLETED', priority: 3, waferCount: 25, createdAt: '2024-01-01T00:00:00Z' },
  { id: 109, lotNumber: 'LOT-2024-009', product: 'SRAM-512K', recipeId: 2, status: 'ON_HOLD',   priority: 4, waferCount: 20, createdAt: '2024-01-01T00:00:00Z' },
  { id: 110, lotNumber: 'LOT-2024-010', product: 'Flash-1G',  recipeId: 3, status: 'COMPLETED', priority: 5, waferCount: 25, createdAt: '2024-01-01T00:00:00Z' },
  { id: 111, lotNumber: 'LOT-2024-011', product: 'Flash-1G',  recipeId: 3, status: 'QUEUED',    priority: 1, waferCount: 20, createdAt: '2024-01-01T00:00:00Z' },
]

describe('useLotStore', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.mocked(getLots).mockImplementation(() => Promise.resolve(mockLots.map(l => ({ ...l }))))
  })

  it('初始載入 11 筆 mock Lot', async () => {
    const store = useLotStore()
    await store.fetchLots()
    expect(store.lots).toHaveLength(11)
  })

  it('wipCount 正確計算 RUNNING 狀態 Lot 數', async () => {
    const store = useLotStore()
    await store.fetchLots()
    const expected = store.lots.filter(l => l.status === 'RUNNING').length
    expect(store.wipCount).toBe(expected)
  })

  it('queuedCount 正確計算 QUEUED 狀態 Lot 數', async () => {
    const store = useLotStore()
    await store.fetchLots()
    const expected = store.lots.filter(l => l.status === 'QUEUED').length
    expect(store.queuedCount).toBe(expected)
  })

  it('updateLotStatus 正確更新指定 Lot 狀態', async () => {
    const store = useLotStore()
    await store.fetchLots()
    store.updateLotStatus(101, 'COMPLETED')
    expect(store.lots.find(l => l.id === 101)?.status).toBe('COMPLETED')
  })

  it('updateLotStatus RUNNING → ON_HOLD 後 wipCount 減少 1', async () => {
    const store = useLotStore()
    await store.fetchLots()
    const before = store.wipCount
    const running = store.lots.find(l => l.status === 'RUNNING')!
    store.updateLotStatus(running.id, 'ON_HOLD')
    expect(store.wipCount).toBe(before - 1)
  })

  it('updateLotStatus QUEUED → RUNNING 後 wipCount 增加 1', async () => {
    const store = useLotStore()
    await store.fetchLots()
    const before = store.wipCount
    const queued = store.lots.find(l => l.status === 'QUEUED')!
    store.updateLotStatus(queued.id, 'RUNNING')
    expect(store.wipCount).toBe(before + 1)
  })

  it('updateLotStatus 傳入不存在的 id 不拋錯', async () => {
    const store = useLotStore()
    await store.fetchLots()
    expect(() => store.updateLotStatus(9999, 'COMPLETED')).not.toThrow()
  })

  it('fetchLots 重置回初始 mock 資料', async () => {
    const store = useLotStore()
    await store.fetchLots()
    store.updateLotStatus(101, 'COMPLETED')
    await store.fetchLots()
    expect(store.lots.find(l => l.id === 101)?.status).toBe('RUNNING')
  })

  it('所有 Lot priority 在 1–5 之間', async () => {
    const store = useLotStore()
    await store.fetchLots()
    store.lots.forEach(l => {
      expect(l.priority).toBeGreaterThanOrEqual(1)
      expect(l.priority).toBeLessThanOrEqual(5)
    })
  })

  it('所有 Lot waferCount 大於 0', async () => {
    const store = useLotStore()
    await store.fetchLots()
    store.lots.forEach(l => {
      expect(l.waferCount).toBeGreaterThan(0)
    })
  })
})
