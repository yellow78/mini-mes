// 設備狀態管理 store
import { defineStore } from 'pinia'
import { computed, ref } from 'vue'
import { getEquipments } from '../api/equipment'
import type { Equipment, EquipmentGroup, EquipmentStatus, EquipmentType } from '../types/mes'

export const useEquipmentStore = defineStore('equipment', () => {
  const equipments  = ref<Equipment[]>([])
  const loading     = ref(false)
  const error       = ref<string | null>(null)

  // 篩選條件
  const filterType      = ref<EquipmentType | 'ALL'>('ALL')
  const filterAlarmOnly = ref(false)
  const filterKeyword   = ref('')

  // 依類型分群，alarm 設備置頂
  const groupedEquipments = computed<EquipmentGroup[]>(() => {
    const types: EquipmentType[] = ['CVD', 'Etch', 'CMP', 'Diffusion']
    return types
      .map(type => {
        const group  = equipments.value.filter(e => e.type === type)
        if (!group.length) return null
        const sorted = [...group].sort((a, b) => (b.isAlarm ? 1 : 0) - (a.isAlarm ? 1 : 0))
        const running = group.filter(e => e.status === 'RUNNING')
        const utilization = running.length
          ? Math.round(running.reduce((sum, e) => sum + e.utilization, 0) / running.length)
          : 0
        return {
          type,
          equipments:  sorted,
          alarmCount:  group.filter(e => e.isAlarm).length,
          utilization,
          statusCount: {
            running: group.filter(e => e.status === 'RUNNING').length,
            idle:    group.filter(e => e.status === 'IDLE').length,
            down:    group.filter(e => e.status === 'DOWN').length,
            pm:      group.filter(e => e.status === 'PM').length,
          },
        }
      })
      .filter((g): g is EquipmentGroup => g !== null)
  })

  // 套用搜尋 + 篩選後的群組列表
  const filteredGroups = computed<EquipmentGroup[]>(() => {
    const keyword = filterKeyword.value.toLowerCase()
    return groupedEquipments.value
      .filter(g => {
        if (filterType.value !== 'ALL' && g.type !== filterType.value) return false
        if (filterAlarmOnly.value && g.alarmCount === 0) return false
        return true
      })
      .map(g => {
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

  // KPI
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

  function setFilter(type: EquipmentType | 'ALL', alarmOnly: boolean, keyword: string) {
    filterType.value      = type
    filterAlarmOnly.value = alarmOnly
    filterKeyword.value   = keyword
  }

  // WebSocket equipment_status_changed 事件觸發
  function updateEquipmentStatus(id: number, status: EquipmentStatus) {
    const eq = equipments.value.find(e => e.id === id)
    if (eq) {
      eq.status    = status
      eq.updatedAt = new Date().toISOString()
    }
  }

  // 從後端 API 取得設備列表
  async function fetchEquipments() {
    loading.value = true
    error.value   = null
    try {
      equipments.value = await getEquipments()
    } catch (e: unknown) {
      error.value = e instanceof Error ? e.message : '載入設備失敗'
    } finally {
      loading.value = false
    }
  }

  return {
    equipments,
    loading,
    error,
    filterType,
    filterAlarmOnly,
    filterKeyword,
    groupedEquipments,
    filteredGroups,
    overallUtilization,
    statusCount,
    setFilter,
    updateEquipmentStatus,
    fetchEquipments,
  }
})
