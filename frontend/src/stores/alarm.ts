// 告警事件 store
import { defineStore } from 'pinia'
import { computed, ref } from 'vue'
import { getAlarms, acknowledgeAlarm as apiAcknowledge } from '../api/alarm'
import type { AlarmEvent } from '../types/mes'

export const useAlarmStore = defineStore('alarm', () => {
  const alarms  = ref<AlarmEvent[]>([])
  const loading = ref(false)

  const unacknowledgedCount = computed(
    () => alarms.value.filter(a => !a.acknowledged).length
  )

  // WebSocket spc_alarm 事件觸發時新增告警
  function addAlarm(alarm: AlarmEvent) {
    alarms.value.unshift(alarm)
  }

  // 確認告警（呼叫 API + 更新本地狀態）
  async function acknowledgeAlarm(id: number) {
    await apiAcknowledge(id)
    const alarm = alarms.value.find(a => a.id === id)
    if (alarm) alarm.acknowledged = true
  }

  // 從後端 API 取得告警列表
  async function fetchAlarms(all = false) {
    loading.value = true
    try {
      alarms.value = await getAlarms(all)
    } finally {
      loading.value = false
    }
  }

  return {
    alarms,
    loading,
    unacknowledgedCount,
    addAlarm,
    acknowledgeAlarm,
    fetchAlarms,
  }
})
