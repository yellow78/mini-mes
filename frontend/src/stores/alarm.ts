// 告警事件 store
import { defineStore } from 'pinia'
import { computed, ref } from 'vue'
import type { AlarmEvent } from '../types/mes'

const MOCK_ALARMS: AlarmEvent[] = [
  {
    id: 1,
    equipmentId: 3,
    equipmentName: 'CVD-03',
    parameter: 'temperature',
    value: 735,
    ucl: 720,
    lcl: 640,
    severity: 'CRITICAL',
    timestamp: new Date(Date.now() - 5 * 60 * 1000).toISOString(),
    acknowledged: false,
  },
  {
    id: 2,
    equipmentId: 9,
    equipmentName: 'Etch-03',
    parameter: 'pressure',
    value: 18,
    ucl: 15,
    lcl: 3,
    severity: 'WARNING',
    timestamp: new Date(Date.now() - 12 * 60 * 1000).toISOString(),
    acknowledged: false,
  },
  {
    id: 3,
    equipmentId: 3,
    equipmentName: 'CVD-03',
    parameter: 'temperature',
    value: 728,
    ucl: 720,
    lcl: 640,
    severity: 'WARNING',
    timestamp: new Date(Date.now() - 25 * 60 * 1000).toISOString(),
    acknowledged: true,
  },
]

export const useAlarmStore = defineStore('alarm', () => {
  const alarms = ref<AlarmEvent[]>(JSON.parse(JSON.stringify(MOCK_ALARMS)))

  const unacknowledgedCount = computed(
    () => alarms.value.filter(a => !a.acknowledged).length
  )

  // WebSocket spc_alarm 事件觸發時新增告警
  function addAlarm(alarm: AlarmEvent) {
    alarms.value.unshift(alarm)
  }

  // 確認告警
  function acknowledgeAlarm(id: number) {
    const alarm = alarms.value.find(a => a.id === id)
    if (alarm) alarm.acknowledged = true
  }

  // Phase 2 後換成真實 API
  async function fetchAlarms() {
    alarms.value = JSON.parse(JSON.stringify(MOCK_ALARMS))
  }

  return {
    alarms,
    unacknowledgedCount,
    addAlarm,
    acknowledgeAlarm,
    fetchAlarms,
  }
})
