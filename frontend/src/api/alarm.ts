// 告警相關 API 呼叫封裝（Phase 2 後啟用）
import axios from 'axios'
import type { AlarmEvent } from '../types/mes'

const BASE = '/api/v1'

export async function getAlarms(): Promise<AlarmEvent[]> {
  const { data } = await axios.get(`${BASE}/alarms`)
  return data
}

export async function acknowledgeAlarm(id: number): Promise<void> {
  await axios.put(`${BASE}/alarms/${id}/acknowledge`)
}
