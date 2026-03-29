// WebSocket 連線管理，統一訂閱事件後分發到各 store
import { onUnmounted, ref } from 'vue'
import { useEquipmentStore } from '../stores/equipment'
import { useAlarmStore } from '../stores/alarm'
import { useLotStore } from '../stores/lot'
import type { AlarmEvent, EquipmentStatus, WSMessage } from '../types/mes'

export function useWebSocket() {
  const connected = ref(false)
  let ws: WebSocket | null = null

  const equipmentStore = useEquipmentStore()
  const alarmStore     = useAlarmStore()
  const lotStore       = useLotStore()

  function connect() {
    const url = `ws://${window.location.hostname}:8080/ws`
    ws = new WebSocket(url)

    ws.onopen = () => {
      connected.value = true
      console.log('[WS] 連線成功')
    }

    ws.onmessage = (event) => {
      try {
        const msg: WSMessage = JSON.parse(event.data)
        handleMessage(msg)
      } catch (e) {
        console.error('[WS] 訊息解析失敗', e)
      }
    }

    ws.onclose = () => {
      connected.value = false
      console.log('[WS] 連線關閉，5 秒後重連...')
      setTimeout(connect, 5000)
    }

    ws.onerror = (err) => {
      console.error('[WS] 連線錯誤', err)
    }
  }

  function handleMessage(msg: WSMessage) {
    switch (msg.event) {
      case 'equipment_status_changed': {
        const { equipment_id, status } = msg.payload as { equipment_id: number; status: EquipmentStatus }
        equipmentStore.updateEquipmentStatus(equipment_id, status)
        break
      }
      case 'spc_alarm': {
        const p = msg.payload as {
          equipment_id: number
          equipment_name: string
          parameter: string
          value: number
          ucl: number
          lcl: number
        }
        const alarm: AlarmEvent = {
          id: Date.now(),
          equipmentId:   p.equipment_id,
          equipmentName: p.equipment_name,
          parameter:     p.parameter,
          value:         p.value,
          ucl:           p.ucl,
          lcl:           p.lcl,
          severity:      p.value > p.ucl * 1.05 ? 'CRITICAL' : 'WARNING',
          timestamp:     new Date().toISOString(),
          acknowledged:  false,
        }
        alarmStore.addAlarm(alarm)
        break
      }
      case 'lot_dispatched': {
        const { lot_id } = msg.payload as { lot_id: number; equipment_id: number }
        lotStore.updateLotStatus(lot_id, 'RUNNING')
        break
      }
    }
  }

  function disconnect() {
    ws?.close()
    ws = null
  }

  onUnmounted(disconnect)

  return { connected, connect, disconnect }
}
