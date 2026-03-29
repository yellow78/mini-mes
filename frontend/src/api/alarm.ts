// 告警相關 API 呼叫封裝
import axios from 'axios'
import type { AlarmEvent } from '../types/mes'

const BASE = '/api/v1'

function mapAlarm(raw: Record<string, unknown>): AlarmEvent {
  return {
    id:            raw.id as number,
    equipmentId:   raw.equipment_id as number,
    equipmentName: raw.equipment_name as string,
    parameter:     raw.parameter as string,
    value:         raw.value as number,
    ucl:           raw.ucl as number,
    lcl:           raw.lcl as number,
    severity:      raw.severity as 'WARNING' | 'CRITICAL',
    timestamp:     raw.timestamp as string,
    acknowledged:  raw.acknowledged as boolean,
  }
}

export async function getAlarms(all = false): Promise<AlarmEvent[]> {
  const { data } = await axios.get<{ data: Record<string, unknown>[] }>(
    `${BASE}/alarms${all ? '?all=true' : ''}`
  )
  return (data.data ?? []).map(mapAlarm)
}

export async function acknowledgeAlarm(id: number): Promise<void> {
  await axios.put(`${BASE}/alarms/${id}/acknowledge`)
}
