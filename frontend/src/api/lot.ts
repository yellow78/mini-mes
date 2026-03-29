// Lot 相關 API 呼叫封裝
import axios from 'axios'
import type { Lot } from '../types/mes'

const BASE = '/api/v1'

function mapLot(raw: Record<string, unknown>): Lot {
  return {
    id:         raw.id as number,
    lotNumber:  raw.lot_number as string,
    product:    raw.product as string,
    recipeId:   raw.recipe_id as number,
    status:     raw.status as Lot['status'],
    priority:   raw.priority as number,
    waferCount: raw.wafer_count as number,
    createdAt:  raw.created_at as string,
  }
}

export async function getLots(): Promise<Lot[]> {
  const { data } = await axios.get<{ data: Record<string, unknown>[] }>(`${BASE}/lots`)
  return (data.data ?? []).map(mapLot)
}

export async function getLot(id: number): Promise<Lot> {
  const { data } = await axios.get<{ data: Record<string, unknown> }>(`${BASE}/lots/${id}`)
  return mapLot(data.data)
}

export interface CreateLotPayload {
  lot_number: string
  product: string
  recipe_id: number
  priority: number
  wafer_count: number
}

export async function createLot(payload: CreateLotPayload): Promise<Lot> {
  const { data } = await axios.post<{ data: Record<string, unknown> }>(`${BASE}/lots`, payload)
  return mapLot(data.data)
}

export async function dispatchLot(id: number): Promise<{ lot_id: number; equipment_id: number }> {
  const { data } = await axios.post<{ data: { lot_id: number; equipment_id: number } }>(
    `${BASE}/lots/${id}/dispatch`
  )
  return data.data
}
