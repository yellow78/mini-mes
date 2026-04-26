import { beforeEach, describe, expect, it, vi } from 'vitest'
import { createPinia, setActivePinia } from 'pinia'
import { useAlarmStore } from '../../stores/alarm'
import type { AlarmEvent } from '../../types/mes'
import { getAlarms } from '../../api/alarm'

vi.mock('../../api/alarm', () => ({
  getAlarms: vi.fn(),
  acknowledgeAlarm: vi.fn().mockResolvedValue(undefined),
}))

const mockAlarms: AlarmEvent[] = [
  {
    id: 1, equipmentId: 1, equipmentName: 'CVD-01', parameter: 'temperature',
    value: 700, ucl: 680, lcl: 600, severity: 'CRITICAL',
    timestamp: '2024-01-01T00:00:00Z', acknowledged: false,
  },
  {
    id: 2, equipmentId: 2, equipmentName: 'CVD-02', parameter: 'pressure',
    value: 420, ucl: 400, lcl: 300, severity: 'WARNING',
    timestamp: '2024-01-01T00:01:00Z', acknowledged: false,
  },
  {
    id: 3, equipmentId: 3, equipmentName: 'Etch-01', parameter: 'temperature',
    value: 560, ucl: 540, lcl: 460, severity: 'WARNING',
    timestamp: '2024-01-01T00:02:00Z', acknowledged: true,
  },
]

describe('useAlarmStore', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.mocked(getAlarms).mockResolvedValue(mockAlarms.map(a => ({ ...a })))
  })

  it('初始載入 3 筆 mock 告警', async () => {
    const store = useAlarmStore()
    await store.fetchAlarms()
    expect(store.alarms).toHaveLength(3)
  })

  it('unacknowledgedCount 正確計算未確認數量', async () => {
    const store = useAlarmStore()
    await store.fetchAlarms()
    const expected = store.alarms.filter(a => !a.acknowledged).length
    expect(store.unacknowledgedCount).toBe(expected)
  })

  it('acknowledgeAlarm 將指定告警標記為已確認', async () => {
    const store = useAlarmStore()
    await store.fetchAlarms()
    const target = store.alarms.find(a => !a.acknowledged)!
    await store.acknowledgeAlarm(target.id)
    expect(store.alarms.find(a => a.id === target.id)?.acknowledged).toBe(true)
  })

  it('acknowledgeAlarm 後 unacknowledgedCount 減少 1', async () => {
    const store = useAlarmStore()
    await store.fetchAlarms()
    const before = store.unacknowledgedCount
    const target = store.alarms.find(a => !a.acknowledged)!
    await store.acknowledgeAlarm(target.id)
    expect(store.unacknowledgedCount).toBe(before - 1)
  })

  it('acknowledgeAlarm 傳入不存在的 id 不拋錯', async () => {
    const store = useAlarmStore()
    await store.fetchAlarms()
    await expect(store.acknowledgeAlarm(9999)).resolves.toBeUndefined()
  })

  it('addAlarm 新增一筆告警至列表最前面', async () => {
    const store = useAlarmStore()
    await store.fetchAlarms()
    const before = store.alarms.length
    const newAlarm: AlarmEvent = {
      id: 999,
      equipmentId: 1,
      equipmentName: 'CVD-01',
      parameter: 'temperature',
      value: 740,
      ucl: 720,
      lcl: 640,
      severity: 'CRITICAL',
      timestamp: new Date().toISOString(),
      acknowledged: false,
    }
    store.addAlarm(newAlarm)
    expect(store.alarms).toHaveLength(before + 1)
    expect(store.alarms[0].id).toBe(999)
  })

  it('addAlarm 後 unacknowledgedCount 增加 1', async () => {
    const store = useAlarmStore()
    await store.fetchAlarms()
    const before = store.unacknowledgedCount
    const newAlarm: AlarmEvent = {
      id: 998,
      equipmentId: 2,
      equipmentName: 'CVD-02',
      parameter: 'pressure',
      value: 420,
      ucl: 400,
      lcl: 300,
      severity: 'WARNING',
      timestamp: new Date().toISOString(),
      acknowledged: false,
    }
    store.addAlarm(newAlarm)
    expect(store.unacknowledgedCount).toBe(before + 1)
  })

  it('fetchAlarms 重置回初始 mock 資料', async () => {
    const store = useAlarmStore()
    await store.fetchAlarms()
    await store.acknowledgeAlarm(store.alarms[0].id)
    await store.fetchAlarms()
    expect(store.alarms).toHaveLength(3)
  })
})
