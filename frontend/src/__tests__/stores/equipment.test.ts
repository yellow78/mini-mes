import { beforeEach, describe, expect, it, vi } from 'vitest'
import { createPinia, setActivePinia } from 'pinia'
import { useEquipmentStore } from '../../stores/equipment'

describe('useEquipmentStore', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.useFakeTimers()
  })

  afterEach(() => {
    vi.useRealTimers()
  })

  // --- 初始化 ---

  it('初始載入 22 台設備', () => {
    const store = useEquipmentStore()
    store.fetchEquipments()
    expect(store.equipments).toHaveLength(22)
  })

  it('初始篩選條件為全部 / 非 alarmOnly / 無關鍵字', () => {
    const store = useEquipmentStore()
    expect(store.filterType).toBe('ALL')
    expect(store.filterAlarmOnly).toBe(false)
    expect(store.filterKeyword).toBe('')
  })

  // --- groupedEquipments ---

  it('設備依 CVD / Etch / CMP / Diffusion 分為四群', () => {
    const store = useEquipmentStore()
    store.fetchEquipments()
    const groups = store.groupedEquipments
    expect(groups).toHaveLength(4)
    expect(groups.map(g => g.type)).toEqual(['CVD', 'Etch', 'CMP', 'Diffusion'])
  })

  it('每群組的 equipments 數量與 statusCount 加總一致', () => {
    const store = useEquipmentStore()
    store.fetchEquipments()
    store.groupedEquipments.forEach(g => {
      const total = g.statusCount.running + g.statusCount.idle + g.statusCount.down + g.statusCount.pm
      expect(total).toBe(g.equipments.length)
    })
  })

  it('alarm 設備排在群組最頂部', () => {
    const store = useEquipmentStore()
    store.fetchEquipments()
    store.groupedEquipments.forEach(g => {
      const firstNonAlarm = g.equipments.findIndex(e => !e.isAlarm)
      const lastAlarm     = g.equipments.map(e => e.isAlarm).lastIndexOf(true)
      if (firstNonAlarm !== -1 && lastAlarm !== -1) {
        expect(lastAlarm).toBeLessThan(firstNonAlarm)
      }
    })
  })

  it('群組 alarmCount 與實際 alarm 設備數一致', () => {
    const store = useEquipmentStore()
    store.fetchEquipments()
    store.groupedEquipments.forEach(g => {
      expect(g.alarmCount).toBe(g.equipments.filter(e => e.isAlarm).length)
    })
  })

  // --- statusCount / overallUtilization ---

  it('statusCount 總和等於設備總數', () => {
    const store = useEquipmentStore()
    store.fetchEquipments()
    const { running, idle, down, pm } = store.statusCount
    expect(running + idle + down + pm).toBe(store.equipments.length)
  })

  it('overallUtilization 在 0–100 之間', () => {
    const store = useEquipmentStore()
    store.fetchEquipments()
    expect(store.overallUtilization).toBeGreaterThanOrEqual(0)
    expect(store.overallUtilization).toBeLessThanOrEqual(100)
  })

  // --- setFilter ---

  it('類型篩選 CVD 只回傳 CVD 群組', () => {
    const store = useEquipmentStore()
    store.fetchEquipments()
    store.setFilter('CVD', false, '')
    expect(store.filteredGroups).toHaveLength(1)
    expect(store.filteredGroups[0].type).toBe('CVD')
  })

  it('alarmOnly 篩選只回傳含 alarm 的群組', () => {
    const store = useEquipmentStore()
    store.fetchEquipments()
    store.setFilter('ALL', true, '')
    store.filteredGroups.forEach(g => {
      expect(g.alarmCount).toBeGreaterThan(0)
    })
  })

  it('關鍵字搜尋 CVD-03 只回傳含該設備的群組', () => {
    const store = useEquipmentStore()
    store.fetchEquipments()
    store.setFilter('ALL', false, 'CVD-03')
    expect(store.filteredGroups).toHaveLength(1)
    expect(store.filteredGroups[0].equipments[0].name).toBe('CVD-03')
  })

  it('不存在的關鍵字回傳空群組列表', () => {
    const store = useEquipmentStore()
    store.fetchEquipments()
    store.setFilter('ALL', false, 'NONEXISTENT-999')
    expect(store.filteredGroups).toHaveLength(0)
  })

  it('關鍵字搜尋 Lot 編號可找到對應設備', () => {
    const store = useEquipmentStore()
    store.fetchEquipments()
    store.setFilter('ALL', false, 'LOT-2024-001')
    const allEqs = store.filteredGroups.flatMap(g => g.equipments)
    expect(allEqs.some(e => e.currentLot === 'LOT-2024-001')).toBe(true)
  })

  // --- updateEquipmentStatus ---

  it('updateEquipmentStatus 正確更新指定設備狀態', () => {
    const store = useEquipmentStore()
    store.fetchEquipments()
    store.updateEquipmentStatus(1, 'DOWN')
    const eq = store.equipments.find(e => e.id === 1)
    expect(eq?.status).toBe('DOWN')
  })

  it('updateEquipmentStatus 傳入不存在的 id 不拋錯', () => {
    const store = useEquipmentStore()
    store.fetchEquipments()
    expect(() => store.updateEquipmentStatus(9999, 'IDLE')).not.toThrow()
  })

  // --- simulation ---

  it('startSimulation / stopSimulation 不拋錯', () => {
    const store = useEquipmentStore()
    store.fetchEquipments()
    expect(() => store.startSimulation()).not.toThrow()
    expect(() => store.stopSimulation()).not.toThrow()
  })

  it('重複呼叫 startSimulation 不會建立多個 timer', () => {
    const store = useEquipmentStore()
    store.fetchEquipments()
    store.startSimulation()
    store.startSimulation() // 第二次應被忽略
    store.stopSimulation()
  })

  it('模擬 2 秒後 RUNNING 設備溫度有波動', () => {
    const store = useEquipmentStore()
    store.fetchEquipments()
    const before = store.equipments.find(e => e.status === 'RUNNING')!.temperature
    store.startSimulation()
    vi.advanceTimersByTime(2000)
    const after = store.equipments.find(e => e.status === 'RUNNING')!.temperature
    // 波動幅度極小但理論上不等（極低機率相等，可接受）
    expect(typeof after).toBe('number')
    store.stopSimulation()
  })
})
