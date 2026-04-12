// Recipe 相關 API 呼叫封裝
import axios from 'axios'

const BASE = '/api/v1'

export interface Recipe {
  id: number
  name: string
  equipment_type: string
  target_temp: number
  target_pressure: number
  duration_min: number
}

export async function getRecipes(): Promise<Recipe[]> {
  const { data } = await axios.get<{ data: Recipe[] }>(`${BASE}/recipes`)
  return data.data ?? []
}
