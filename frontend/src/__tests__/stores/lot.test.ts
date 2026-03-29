import { beforeEach, describe, expect, it } from 'vitest'
import { createPinia, setActivePinia } from 'pinia'
import { useLotStore } from '../../stores/lot'

describe('useLotStore', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  it('初始載入 11 筆 mock Lot', () => {
    const store = useLotStore()
    store.fetchLots()
    expect(store.lots).toHaveLength(11)
  })

  it('wipCount 正確計算 RUNNING 狀態 Lot 數', () => {
    const store = useLotStore()
    store.fetchLots()
    const expected = store.lots.filter(l => l.status === 'RUNNING').length
    expect(store.wipCount).toBe(expected)
  })

  it('queuedCount 正確計算 QUEUED 狀態 Lot 數', () => {
    const store = useLotStore()
    store.fetchLots()
    const expected = store.lots.filter(l => l.status === 'QUEUED').length
    expect(store.queuedCount).toBe(expected)
  })

  it('updateLotStatus 正確更新指定 Lot 狀態', () => {
    const store = useLotStore()
    store.fetchLots()
    store.updateLotStatus(101, 'COMPLETED')
    expect(store.lots.find(l => l.id === 101)?.status).toBe('COMPLETED')
  })

  it('updateLotStatus RUNNING → ON_HOLD 後 wipCount 減少 1', () => {
    const store = useLotStore()
    store.fetchLots()
    const before = store.wipCount
    const running = store.lots.find(l => l.status === 'RUNNING')!
    store.updateLotStatus(running.id, 'ON_HOLD')
    expect(store.wipCount).toBe(before - 1)
  })

  it('updateLotStatus QUEUED → RUNNING 後 wipCount 增加 1', () => {
    const store = useLotStore()
    store.fetchLots()
    const before = store.wipCount
    const queued = store.lots.find(l => l.status === 'QUEUED')!
    store.updateLotStatus(queued.id, 'RUNNING')
    expect(store.wipCount).toBe(before + 1)
  })

  it('updateLotStatus 傳入不存在的 id 不拋錯', () => {
    const store = useLotStore()
    store.fetchLots()
    expect(() => store.updateLotStatus(9999, 'COMPLETED')).not.toThrow()
  })

  it('fetchLots 重置回初始 mock 資料', async () => {
    const store = useLotStore()
    store.fetchLots()
    store.updateLotStatus(101, 'COMPLETED')
    await store.fetchLots()
    expect(store.lots.find(l => l.id === 101)?.status).toBe('RUNNING')
  })

  it('所有 Lot priority 在 1–5 之間', () => {
    const store = useLotStore()
    store.fetchLots()
    store.lots.forEach(l => {
      expect(l.priority).toBeGreaterThanOrEqual(1)
      expect(l.priority).toBeLessThanOrEqual(5)
    })
  })

  it('所有 Lot waferCount 大於 0', () => {
    const store = useLotStore()
    store.fetchLots()
    store.lots.forEach(l => {
      expect(l.waferCount).toBeGreaterThan(0)
    })
  })
})
