// 設備相關 API 呼叫封裝（Phase 2 後啟用）
import axios from 'axios'
import type { Equipment, EquipmentStatus } from '../types/mes'

const BASE = '/api/v1'

export async function getEquipments(): Promise<Equipment[]> {
  const { data } = await axios.get(`${BASE}/equipment`)
  return data
}

export async function getEquipment(id: number): Promise<Equipment> {
  const { data } = await axios.get(`${BASE}/equipment/${id}`)
  return data
}

export async function updateEquipmentStatus(id: number, status: EquipmentStatus): Promise<void> {
  await axios.put(`${BASE}/equipment/${id}/status`, { status })
}

export async function holdEquipment(id: number): Promise<void> {
  await axios.post(`${BASE}/equipment/${id}/hold`)
}
