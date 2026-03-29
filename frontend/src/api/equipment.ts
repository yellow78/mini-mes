// 設備相關 API 呼叫封裝
import axios from 'axios'
import type { Equipment, EquipmentStatus, SpcDataPoint } from '../types/mes'

const BASE = '/api/v1'

// 後端回傳 snake_case，對應前端 camelCase 型別
function mapEquipment(raw: Record<string, unknown>): Equipment {
  return {
    id:           raw.id as number,
    name:         raw.name as string,
    type:         raw.type as Equipment['type'],
    status:       raw.status as EquipmentStatus,
    currentLotId: raw.current_lot_id as number | null,
    currentLot:   raw.current_lot as string | null,
    recipeName:   raw.recipe_name as string | null,
    utilization:  raw.utilization as number,
    temperature:  raw.temperature as number,
    pressure:     raw.pressure as number,
    ucl_temp:     raw.ucl_temp as number,
    lcl_temp:     raw.lcl_temp as number,
    ucl_pressure: raw.ucl_pressure as number,
    lcl_pressure: raw.lcl_pressure as number,
    isAlarm:      raw.is_alarm as boolean,
    updatedAt:    raw.updated_at as string,
  }
}

export async function getEquipments(): Promise<Equipment[]> {
  const { data } = await axios.get<{ data: Record<string, unknown>[] }>(`${BASE}/equipment`)
  return (data.data ?? []).map(mapEquipment)
}

export async function getEquipment(id: number): Promise<Equipment> {
  const { data } = await axios.get<{ data: Record<string, unknown> }>(`${BASE}/equipment/${id}`)
  return mapEquipment(data.data)
}

export async function updateEquipmentStatus(id: number, status: EquipmentStatus): Promise<void> {
  await axios.put(`${BASE}/equipment/${id}/status`, { status })
}

export async function holdEquipment(id: number): Promise<void> {
  await axios.post(`${BASE}/equipment/${id}/hold`)
}

// SPC 歷史資料（最近 N 筆）
export async function getSpcHistory(equipmentId: number, limit = 20): Promise<SpcDataPoint[]> {
  const { data } = await axios.get<{ data: Record<string, unknown>[] }>(
    `${BASE}/spc/${equipmentId}?limit=${limit}`
  )
  return (data.data ?? []).map(r => ({
    timestamp: r.timestamp as string,
    value:     r.value as number,
    ucl:       r.ucl as number,
    lcl:       r.lcl as number,
    isAlarm:   r.is_alarm as boolean,
  }))
}
