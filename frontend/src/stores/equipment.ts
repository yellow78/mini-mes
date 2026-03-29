// 設備狀態管理 store，Phase 1 使用 mock 資料 + 模擬波動
import { defineStore } from 'pinia'
import { computed, ref } from 'vue'
import type { Equipment, EquipmentGroup, EquipmentStatus, EquipmentType } from '../types/mes'

// Mock 資料：每種類型 4–7 台，含至少 2 台 Alarm 設備
const MOCK_EQUIPMENTS: Equipment[] = [
  // CVD 群組（6台）
  { id: 1,  name: 'CVD-01', type: 'CVD', status: 'RUNNING', currentLotId: 101, currentLot: 'LOT-2024-001', recipeName: 'CVD-SiO2-STD', utilization: 85, temperature: 680, pressure: 350, ucl_temp: 720, lcl_temp: 640, ucl_pressure: 400, lcl_pressure: 300, isAlarm: false, updatedAt: new Date().toISOString() },
  { id: 2,  name: 'CVD-02', type: 'CVD', status: 'RUNNING', currentLotId: 102, currentLot: 'LOT-2024-002', recipeName: 'CVD-SiN-STD',  utilization: 92, temperature: 695, pressure: 370, ucl_temp: 720, lcl_temp: 640, ucl_pressure: 400, lcl_pressure: 300, isAlarm: false, updatedAt: new Date().toISOString() },
  { id: 3,  name: 'CVD-03', type: 'CVD', status: 'RUNNING', currentLotId: 103, currentLot: 'LOT-2024-003', recipeName: 'CVD-SiO2-STD', utilization: 78, temperature: 735, pressure: 360, ucl_temp: 720, lcl_temp: 640, ucl_pressure: 400, lcl_pressure: 300, isAlarm: true,  updatedAt: new Date().toISOString() },
  { id: 4,  name: 'CVD-04', type: 'CVD', status: 'IDLE',    currentLotId: null, currentLot: null, recipeName: null, utilization: 45, temperature: 25,  pressure: 50,  ucl_temp: 720, lcl_temp: 640, ucl_pressure: 400, lcl_pressure: 300, isAlarm: false, updatedAt: new Date().toISOString() },
  { id: 5,  name: 'CVD-05', type: 'CVD', status: 'PM',      currentLotId: null, currentLot: null, recipeName: null, utilization: 0,  temperature: 30,  pressure: 10,  ucl_temp: 720, lcl_temp: 640, ucl_pressure: 400, lcl_pressure: 300, isAlarm: false, updatedAt: new Date().toISOString() },
  { id: 6,  name: 'CVD-06', type: 'CVD', status: 'DOWN',    currentLotId: null, currentLot: null, recipeName: null, utilization: 0,  temperature: 28,  pressure: 8,   ucl_temp: 720, lcl_temp: 640, ucl_pressure: 400, lcl_pressure: 300, isAlarm: false, updatedAt: new Date().toISOString() },

  // Etch 群組（6台）
  { id: 7,  name: 'Etch-01', type: 'Etch', status: 'RUNNING', currentLotId: 201, currentLot: 'LOT-2024-010', recipeName: 'Etch-Poly-DRY',  utilization: 88, temperature: 45,  pressure: 8,   ucl_temp: 80,  lcl_temp: 20,  ucl_pressure: 15,  lcl_pressure: 3,  isAlarm: false, updatedAt: new Date().toISOString() },
  { id: 8,  name: 'Etch-02', type: 'Etch', status: 'RUNNING', currentLotId: 202, currentLot: 'LOT-2024-011', recipeName: 'Etch-Oxide-DRY', utilization: 76, temperature: 38,  pressure: 6,   ucl_temp: 80,  lcl_temp: 20,  ucl_pressure: 15,  lcl_pressure: 3,  isAlarm: false, updatedAt: new Date().toISOString() },
  { id: 9,  name: 'Etch-03', type: 'Etch', status: 'RUNNING', currentLotId: 203, currentLot: 'LOT-2024-012', recipeName: 'Etch-Poly-DRY',  utilization: 91, temperature: 52,  pressure: 18,  ucl_temp: 80,  lcl_temp: 20,  ucl_pressure: 15,  lcl_pressure: 3,  isAlarm: true,  updatedAt: new Date().toISOString() },
  { id: 10, name: 'Etch-04', type: 'Etch', status: 'IDLE',    currentLotId: null, currentLot: null, recipeName: null, utilization: 50, temperature: 22,  pressure: 2,   ucl_temp: 80,  lcl_temp: 20,  ucl_pressure: 15,  lcl_pressure: 3,  isAlarm: false, updatedAt: new Date().toISOString() },
  { id: 11, name: 'Etch-05', type: 'Etch', status: 'IDLE',    currentLotId: null, currentLot: null, recipeName: null, utilization: 42, temperature: 25,  pressure: 3,   ucl_temp: 80,  lcl_temp: 20,  ucl_pressure: 15,  lcl_pressure: 3,  isAlarm: false, updatedAt: new Date().toISOString() },
  { id: 12, name: 'Etch-06', type: 'Etch', status: 'DOWN',    currentLotId: null, currentLot: null, recipeName: null, utilization: 0,  temperature: 20,  pressure: 1,   ucl_temp: 80,  lcl_temp: 20,  ucl_pressure: 15,  lcl_pressure: 3,  isAlarm: false, updatedAt: new Date().toISOString() },

  // CMP 群組（5台）
  { id: 13, name: 'CMP-01', type: 'CMP', status: 'RUNNING', currentLotId: 301, currentLot: 'LOT-2024-020', recipeName: 'CMP-STI-STD',  utilization: 82, temperature: 60,  pressure: 180, ucl_temp: 90,  lcl_temp: 40,  ucl_pressure: 220, lcl_pressure: 140, isAlarm: false, updatedAt: new Date().toISOString() },
  { id: 14, name: 'CMP-02', type: 'CMP', status: 'RUNNING', currentLotId: 302, currentLot: 'LOT-2024-021', recipeName: 'CMP-Metal-STD', utilization: 95, temperature: 65,  pressure: 200, ucl_temp: 90,  lcl_temp: 40,  ucl_pressure: 220, lcl_pressure: 140, isAlarm: false, updatedAt: new Date().toISOString() },
  { id: 15, name: 'CMP-03', type: 'CMP', status: 'IDLE',    currentLotId: null, currentLot: null, recipeName: null, utilization: 55, temperature: 30,  pressure: 50,  ucl_temp: 90,  lcl_temp: 40,  ucl_pressure: 220, lcl_pressure: 140, isAlarm: false, updatedAt: new Date().toISOString() },
  { id: 16, name: 'CMP-04', type: 'CMP', status: 'PM',      currentLotId: null, currentLot: null, recipeName: null, utilization: 0,  temperature: 25,  pressure: 10,  ucl_temp: 90,  lcl_temp: 40,  ucl_pressure: 220, lcl_pressure: 140, isAlarm: false, updatedAt: new Date().toISOString() },
  { id: 17, name: 'CMP-05', type: 'CMP', status: 'IDLE',    currentLotId: null, currentLot: null, recipeName: null, utilization: 38, temperature: 28,  pressure: 45,  ucl_temp: 90,  lcl_temp: 40,  ucl_pressure: 220, lcl_pressure: 140, isAlarm: false, updatedAt: new Date().toISOString() },

  // Diffusion 群組（5台）
  { id: 18, name: 'Diff-01', type: 'Diffusion', status: 'RUNNING', currentLotId: 401, currentLot: 'LOT-2024-030', recipeName: 'Diff-N-Well',  utilization: 90, temperature: 950, pressure: 760, ucl_temp: 1000, lcl_temp: 900, ucl_pressure: 800, lcl_pressure: 720, isAlarm: false, updatedAt: new Date().toISOString() },
  { id: 19, name: 'Diff-02', type: 'Diffusion', status: 'RUNNING', currentLotId: 402, currentLot: 'LOT-2024-031', recipeName: 'Diff-P-Well',  utilization: 87, temperature: 960, pressure: 755, ucl_temp: 1000, lcl_temp: 900, ucl_pressure: 800, lcl_pressure: 720, isAlarm: false, updatedAt: new Date().toISOString() },
  { id: 20, name: 'Diff-03', type: 'Diffusion', status: 'IDLE',    currentLotId: null, currentLot: null, recipeName: null, utilization: 60, temperature: 25,  pressure: 760, ucl_temp: 1000, lcl_temp: 900, ucl_pressure: 800, lcl_pressure: 720, isAlarm: false, updatedAt: new Date().toISOString() },
  { id: 21, name: 'Diff-04', type: 'Diffusion', status: 'DOWN',    currentLotId: null, currentLot: null, recipeName: null, utilization: 0,  temperature: 25,  pressure: 760, ucl_temp: 1000, lcl_temp: 900, ucl_pressure: 800, lcl_pressure: 720, isAlarm: false, updatedAt: new Date().toISOString() },
  { id: 22, name: 'Diff-05', type: 'Diffusion', status: 'RUNNING', currentLotId: 403, currentLot: 'LOT-2024-032', recipeName: 'Diff-N-Well',  utilization: 83, temperature: 945, pressure: 758, ucl_temp: 1000, lcl_temp: 900, ucl_pressure: 800, lcl_pressure: 720, isAlarm: false, updatedAt: new Date().toISOString() },
]

export const useEquipmentStore = defineStore('equipment', () => {
  const equipments = ref<Equipment[]>(JSON.parse(JSON.stringify(MOCK_EQUIPMENTS)))

  // 篩選條件
  const filterType    = ref<EquipmentType | 'ALL'>('ALL')
  const filterAlarmOnly = ref(false)
  const filterKeyword = ref('')

  // 模擬計時器
  let simulationTimer: ReturnType<typeof setInterval> | null = null

  // 依類型分群，alarm 設備置頂
  const groupedEquipments = computed<EquipmentGroup[]>(() => {
    const types: EquipmentType[] = ['CVD', 'Etch', 'CMP', 'Diffusion']
    return types.map(type => {
      const group = equipments.value.filter(e => e.type === type)
      // alarm 設備排最前面
      const sorted = [...group].sort((a, b) => (b.isAlarm ? 1 : 0) - (a.isAlarm ? 1 : 0))
      const alarmCount = group.filter(e => e.isAlarm).length
      const utilization = group.length
        ? Math.round(group.reduce((sum, e) => sum + e.utilization, 0) / group.length)
        : 0
      return {
        type,
        equipments: sorted,
        alarmCount,
        utilization,
        statusCount: {
          running: group.filter(e => e.status === 'RUNNING').length,
          idle:    group.filter(e => e.status === 'IDLE').length,
          down:    group.filter(e => e.status === 'DOWN').length,
          pm:      group.filter(e => e.status === 'PM').length,
        },
      }
    })
  })

  // 套用搜尋 + 篩選後的群組列表
  const filteredGroups = computed<EquipmentGroup[]>(() => {
    const keyword = filterKeyword.value.toLowerCase()
    return groupedEquipments.value
      .filter(g => {
        // 類型篩選
        if (filterType.value !== 'ALL' && g.type !== filterType.value) return false
        // 僅看 Alarm
        if (filterAlarmOnly.value && g.alarmCount === 0) return false
        return true
      })
      .map(g => {
        // 搜尋 filter：過濾設備名或 Lot 編號
        if (!keyword) return g
        const filtered = g.equipments.filter(
          e =>
            e.name.toLowerCase().includes(keyword) ||
            (e.currentLot?.toLowerCase().includes(keyword) ?? false)
        )
        return { ...g, equipments: filtered }
      })
      .filter(g => g.equipments.length > 0)
  })

  // KPI 計算
  const overallUtilization = computed(() => {
    const running = equipments.value.filter(e => e.status === 'RUNNING')
    if (!running.length) return 0
    return Math.round(running.reduce((sum, e) => sum + e.utilization, 0) / running.length)
  })

  const statusCount = computed(() => ({
    running: equipments.value.filter(e => e.status === 'RUNNING').length,
    idle:    equipments.value.filter(e => e.status === 'IDLE').length,
    down:    equipments.value.filter(e => e.status === 'DOWN').length,
    pm:      equipments.value.filter(e => e.status === 'PM').length,
    alarm:   equipments.value.filter(e => e.isAlarm).length,
  }))

  // 更新篩選條件
  function setFilter(type: EquipmentType | 'ALL', alarmOnly: boolean, keyword: string) {
    filterType.value    = type
    filterAlarmOnly.value = alarmOnly
    filterKeyword.value = keyword
  }

  // WebSocket 事件觸發時更新設備狀態
  function updateEquipmentStatus(id: number, status: EquipmentStatus) {
    const eq = equipments.value.find(e => e.id === id)
    if (eq) {
      eq.status    = status
      eq.updatedAt = new Date().toISOString()
    }
  }

  // Phase 1 Mock 資料模擬：每 2 秒微幅波動溫度/壓力/稼動率
  function startSimulation() {
    if (simulationTimer) return
    simulationTimer = setInterval(() => {
      equipments.value.forEach(eq => {
        if (eq.status !== 'RUNNING') return
        // 微幅波動 ±1%
        const delta = (Math.random() - 0.5) * 2
        eq.temperature  = parseFloat((eq.temperature  + delta * 0.5).toFixed(1))
        eq.pressure     = parseFloat((eq.pressure     + delta * 0.3).toFixed(1))
        eq.utilization  = Math.max(0, Math.min(100, parseFloat((eq.utilization + delta * 0.2).toFixed(1))))
        // 更新 alarm 狀態
        eq.isAlarm = eq.temperature > eq.ucl_temp || eq.temperature < eq.lcl_temp ||
                     eq.pressure    > eq.ucl_pressure || eq.pressure < eq.lcl_pressure
        eq.updatedAt = new Date().toISOString()
      })
    }, 2000)
  }

  function stopSimulation() {
    if (simulationTimer) {
      clearInterval(simulationTimer)
      simulationTimer = null
    }
  }

  // Phase 2 後換成真實 API 呼叫
  async function fetchEquipments() {
    // TODO: Phase 2 改為 import { getEquipments } from '../api/equipment'
    equipments.value = JSON.parse(JSON.stringify(MOCK_EQUIPMENTS))
  }

  return {
    equipments,
    filterType,
    filterAlarmOnly,
    filterKeyword,
    groupedEquipments,
    filteredGroups,
    overallUtilization,
    statusCount,
    setFilter,
    updateEquipmentStatus,
    startSimulation,
    stopSimulation,
    fetchEquipments,
  }
})
