// Lot 相關 API 呼叫封裝（Phase 2 後啟用）
import axios from 'axios'
import type { Lot } from '../types/mes'

const BASE = '/api/v1'

export async function getLots(): Promise<Lot[]> {
  const { data } = await axios.get(`${BASE}/lots`)
  return data
}

export async function getLot(id: number): Promise<Lot> {
  const { data } = await axios.get(`${BASE}/lots/${id}`)
  return data
}

export async function createLot(payload: Partial<Lot>): Promise<Lot> {
  const { data } = await axios.post(`${BASE}/lots`, payload)
  return data
}

export async function dispatchLot(id: number): Promise<void> {
  await axios.post(`${BASE}/lots/${id}/dispatch`)
}
