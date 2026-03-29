// MES 核心型別定義，所有型別集中於此，元件不重複定義

export type EquipmentStatus = 'RUNNING' | 'IDLE' | 'DOWN' | 'PM'
export type EquipmentType   = 'CVD' | 'Etch' | 'CMP' | 'Diffusion'

export interface Equipment {
  id: number
  name: string
  type: EquipmentType
  status: EquipmentStatus
  currentLotId: number | null
  currentLot: string | null      // lot_number，顯示用
  recipeName: string | null
  utilization: number            // 0–100
  temperature: number            // °C
  pressure: number               // mTorr
  ucl_temp: number
  lcl_temp: number
  ucl_pressure: number
  lcl_pressure: number
  isAlarm: boolean
  updatedAt: string
}

export type LotStatus = 'QUEUED' | 'RUNNING' | 'COMPLETED' | 'ON_HOLD'

export interface Lot {
  id: number
  lotNumber: string
  product: string
  recipeId: number
  status: LotStatus
  priority: number               // 1=最高, 5=最低
  waferCount: number
  createdAt: string
}

export interface AlarmEvent {
  id: number
  equipmentId: number
  equipmentName: string
  parameter: string              // 'temperature' | 'pressure'
  value: number
  ucl: number
  lcl: number
  severity: 'WARNING' | 'CRITICAL'
  timestamp: string
  acknowledged: boolean
}

export interface WSMessage {
  event: 'equipment_status_changed' | 'spc_alarm' | 'lot_dispatched'
  payload: Record<string, unknown>
}

// 群組折疊用
export interface EquipmentGroup {
  type: EquipmentType
  equipments: Equipment[]
  alarmCount: number
  utilization: number            // 群組平均稼動率
  statusCount: {
    running: number
    idle: number
    down: number
    pm: number
  }
}

// SPC 資料點
export interface SpcDataPoint {
  timestamp: string
  value: number
  ucl: number
  lcl: number
  isAlarm: boolean
}
