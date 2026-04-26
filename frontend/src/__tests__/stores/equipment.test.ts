import { beforeEach, describe, expect, it, vi } from 'vitest'
import { createPinia, setActivePinia } from 'pinia'
import { useEquipmentStore } from '../../stores/equipment'
import type { Equipment } from '../../types/mes'
import { getEquipments } from '../../api/equipment'

vi.mock('../../api/equipment', () => ({
  getEquipments: vi.fn(),
}))

// 22 台設備：CVD×6, Etch×6, CMP×5, Diffusion×5
const makeEquip = (override: Partial<Equipment>): Equipment => ({
  id: 1, name: 'CVD-01', type: 'CVD', status: 'RUNNING',
  currentLotId: null, currentLot: null, recipeName: null,
  utilization: 80, temperature: 650, pressure: 120,
  ucl_temp: 680, lcl_temp: 600, ucl_pressure: 150, lcl_pressure: 80,
  isAlarm: false, updatedAt: '2024-01-01T00:00:00Z',
  ...override,
})

const mockEquipments: Equipment[] = [
  // CVD (6)
  makeEquip({ id: 1,  name: 'CVD-01', type: 'CVD', status: 'RUNNING', currentLot: 'LOT-2024-001', currentLotId: 101 }),
  makeEquip({ id: 2,  name: 'CVD-02', type: 'CVD', status: 'RUNNING', isAlarm: true, temperature: 710 }),
  makeEquip({ id: 3,  name: 'CVD-03', type: 'CVD', status: 'IDLE' }),
  makeEquip({ id: 4,  name: 'CVD-04', type: 'CVD', status: 'DOWN' }),
  makeEquip({ id: 5,  name: 'CVD-05', type: 'CVD', status: 'PM' }),
  makeEquip({ id: 6,  name: 'CVD-06', type: 'CVD', status: 'IDLE' }),
  // Etch (6)
  makeEquip({ id: 7,  name: 'Etch-01', type: 'Etch', status: 'RUNNING' }),
  makeEquip({ id: 8,  name: 'Etch-02', type: 'Etch', status: 'RUNNING', isAlarm: true }),
  makeEquip({ id: 9,  name: 'Etch-03', type: 'Etch', status: 'IDLE' }),
  makeEquip({ id: 10, name: 'Etch-04', type: 'Etch', status: 'DOWN' }),
  makeEquip({ id: 11, name: 'Etch-05', type: 'Etch', status: 'PM' }),
  makeEquip({ id: 12, name: 'Etch-06', type: 'Etch', status: 'IDLE' }),
  // CMP (5)
  makeEquip({ id: 13, name: 'CMP-01', type: 'CMP', status: 'RUNNING' }),
  makeEquip({ id: 14, name: 'CMP-02', type: 'CMP', status: 'RUNNING' }),
  makeEquip({ id: 15, name: 'CMP-03', type: 'CMP', status: 'IDLE' }),
  makeEquip({ id: 16, name: 'CMP-04', type: 'CMP', status: 'DOWN' }),
  makeEquip({ id: 17, name: 'CMP-05', type: 'CMP', status: 'PM' }),
  // Diffusion (5)
  makeEquip({ id: 18, name: 'Diff-01', type: 'Diffusion', status: 'RUNNING' }),
  makeEquip({ id: 19, name: 'Diff-02', type: 'Diffusion', status: 'RUNNING' }),
  makeEquip({ id: 20, name: 'Diff-03', type: 'Diffusion', status: 'IDLE' }),
  makeEquip({ id: 21, name: 'Diff-04', type: 'Diffusion', status: 'DOWN' }),
  makeEquip({ id: 22, name: 'Diff-05', type: 'Diffusion', status: 'PM' }),
]

describe('useEquipmentStore', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.mocked(getEquipments).mockResolvedValue(mockEquipments.map(e => ({ ...e })))
  })

  // --- 初始化 ---

  it('初始載入 22 台設備', async () => {
    const store = useEquipmentStore()
    await store.fetchEquipments()
    expect(store.equipments).toHaveLength(22)
  })

  it('初始篩選條件為全部 / 非 alarmOnly / 無關鍵字', () => {
    const store = useEquipmentStore()
    expect(store.filterType).toBe('ALL')
    expect(store.filterAlarmOnly).toBe(false)
    expect(store.filterKeyword).toBe('')
  })

  // --- groupedEquipments ---

  it('設備依 CVD / Etch / CMP / Diffusion 分為四群', async () => {
    const store = useEquipmentStore()
    await store.fetchEquipments()
    const groups = store.groupedEquipments
    expect(groups).toHaveLength(4)
    expect(groups.map(g => g.type)).toEqual(['CVD', 'Etch', 'CMP', 'Diffusion'])
  })

  it('每群組的 equipments 數量與 statusCount 加總一致', async () => {
    const store = useEquipmentStore()
    await store.fetchEquipments()
    store.groupedEquipments.forEach(g => {
      const total = g.statusCount.running + g.statusCount.idle + g.statusCount.down + g.statusCount.pm
      expect(total).toBe(g.equipments.length)
    })
  })

  it('alarm 設備排在群組最頂部', async () => {
    const store = useEquipmentStore()
    await store.fetchEquipments()
    store.groupedEquipments.forEach(g => {
      const firstNonAlarm = g.equipments.findIndex(e => !e.isAlarm)
      const lastAlarm     = g.equipments.map(e => e.isAlarm).lastIndexOf(true)
      if (firstNonAlarm !== -1 && lastAlarm !== -1) {
        expect(lastAlarm).toBeLessThan(firstNonAlarm)
      }
    })
  })

  it('群組 alarmCount 與實際 alarm 設備數一致', async () => {
    const store = useEquipmentStore()
    await store.fetchEquipments()
    store.groupedEquipments.forEach(g => {
      expect(g.alarmCount).toBe(g.equipments.filter(e => e.isAlarm).length)
    })
  })

  // --- statusCount / overallUtilization ---

  it('statusCount 總和等於設備總數', async () => {
    const store = useEquipmentStore()
    await store.fetchEquipments()
    const { running, idle, down, pm } = store.statusCount
    expect(running + idle + down + pm).toBe(store.equipments.length)
  })

  it('overallUtilization 在 0–100 之間', async () => {
    const store = useEquipmentStore()
    await store.fetchEquipments()
    expect(store.overallUtilization).toBeGreaterThanOrEqual(0)
    expect(store.overallUtilization).toBeLessThanOrEqual(100)
  })

  // --- setFilter ---

  it('類型篩選 CVD 只回傳 CVD 群組', async () => {
    const store = useEquipmentStore()
    await store.fetchEquipments()
    store.setFilter('CVD', false, '')
    expect(store.filteredGroups).toHaveLength(1)
    expect(store.filteredGroups[0].type).toBe('CVD')
  })

  it('alarmOnly 篩選只回傳含 alarm 的群組', async () => {
    const store = useEquipmentStore()
    await store.fetchEquipments()
    store.setFilter('ALL', true, '')
    store.filteredGroups.forEach(g => {
      expect(g.alarmCount).toBeGreaterThan(0)
    })
  })

  it('關鍵字搜尋 CVD-03 只回傳含該設備的群組', async () => {
    const store = useEquipmentStore()
    await store.fetchEquipments()
    store.setFilter('ALL', false, 'CVD-03')
    expect(store.filteredGroups).toHaveLength(1)
    expect(store.filteredGroups[0].equipments[0].name).toBe('CVD-03')
  })

  it('不存在的關鍵字回傳空群組列表', async () => {
    const store = useEquipmentStore()
    await store.fetchEquipments()
    store.setFilter('ALL', false, 'NONEXISTENT-999')
    expect(store.filteredGroups).toHaveLength(0)
  })

  it('關鍵字搜尋 Lot 編號可找到對應設備', async () => {
    const store = useEquipmentStore()
    await store.fetchEquipments()
    store.setFilter('ALL', false, 'LOT-2024-001')
    const allEqs = store.filteredGroups.flatMap(g => g.equipments)
    expect(allEqs.some(e => e.currentLot === 'LOT-2024-001')).toBe(true)
  })

  // --- updateEquipmentStatus ---

  it('updateEquipmentStatus 正確更新指定設備狀態', async () => {
    const store = useEquipmentStore()
    await store.fetchEquipments()
    store.updateEquipmentStatus(1, 'DOWN')
    const eq = store.equipments.find(e => e.id === 1)
    expect(eq?.status).toBe('DOWN')
  })

  it('updateEquipmentStatus 傳入不存在的 id 不拋錯', async () => {
    const store = useEquipmentStore()
    await store.fetchEquipments()
    expect(() => store.updateEquipmentStatus(9999, 'IDLE')).not.toThrow()
  })
})
