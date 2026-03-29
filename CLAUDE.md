# Mini-MES — Claude CLI 開發工作手冊

## Claude 互動規範
- **回答語言：繁體中文**
- 程式碼內的註解使用繁體中文
- 變數名稱、函式名稱、檔案名稱維持英文（程式碼慣例）
- 架構說明、錯誤排查、開發建議一律用繁體中文回答

---

## 專案定位
半導體製造執行系統（MES）展示專案。
目標：展示 MES 領域知識 + 全端開發能力 + 系統設計思維。
Demo 時間約 3–5 分鐘，需要能即時跑起來。

---

## 技術棧

| 層級       | 技術                                      | 目錄           |
|------------|-------------------------------------------|----------------|
| 前端       | Vue 3 + TypeScript + Element Plus + Vite  | `frontend/`    |
| 後端 API   | Go + Gin + WebSocket (gorilla)            | `backend/`     |
| 分析引擎   | Python + FastAPI + pandas                 | `analytics/`   |
| 資料庫     | PostgreSQL + Redis                        | `migrations/`  |
| 容器化     | Docker Compose                            | 根目錄         |

---

## 專案目錄結構

```
mini-mes/
├── CLAUDE.md                        ← 本文件
├── docker-compose.yml
├── migrations/
│   ├── 001_init.sql
│   └── 002_add_spc.sql
├── docs/
│   ├── api-spec.md
│   └── db-schema.md
├── frontend/                        ← Vue 3 專案
│   ├── src/
│   │   ├── types/
│   │   │   └── mes.ts               ← 所有型別定義集中於此
│   │   ├── stores/
│   │   │   ├── equipment.ts
│   │   │   ├── alarm.ts
│   │   │   └── lot.ts
│   │   ├── api/
│   │   │   ├── equipment.ts
│   │   │   ├── lot.ts
│   │   │   └── alarm.ts
│   │   ├── composables/
│   │   │   └── useWebSocket.ts
│   │   ├── components/
│   │   │   ├── layout/
│   │   │   │   ├── AppHeader.vue
│   │   │   │   └── AppSidebar.vue
│   │   │   ├── equipment/
│   │   │   │   ├── EquipmentGroupList.vue   ← 群組折疊主元件
│   │   │   │   ├── EquipmentGroup.vue       ← 單一群組（可折疊）
│   │   │   │   ├── EquipmentRow.vue         ← 單台設備列
│   │   │   │   └── EquipmentDrawer.vue      ← 點入後右側詳細抽屜
│   │   │   ├── alarm/
│   │   │   │   └── AlarmList.vue
│   │   │   └── spc/
│   │   │       └── SpcMiniChart.vue
│   │   ├── views/
│   │   │   ├── DashboardView.vue
│   │   │   ├── LotView.vue
│   │   │   └── SpcView.vue
│   │   └── router/
│   │       └── index.ts
│   ├── package.json
│   └── vite.config.ts
├── backend/                         ← Go 專案
│   ├── cmd/server/main.go
│   ├── internal/
│   │   ├── handler/
│   │   │   ├── equipment.go
│   │   │   ├── lot.go
│   │   │   └── alarm.go
│   │   ├── service/
│   │   │   ├── equipment.go
│   │   │   ├── lot.go
│   │   │   └── dispatch.go
│   │   ├── repository/
│   │   │   ├── equipment.go
│   │   │   └── lot.go
│   │   └── model/
│   │       ├── equipment.go
│   │       ├── lot.go
│   │       └── spc_record.go
│   ├── pkg/websocket/
│   │   └── hub.go
│   └── go.mod
└── analytics/                       ← Python 分析引擎
    ├── main.py
    ├── spc/
    │   └── control_chart.py
    └── requirements.txt
```

---

## 開發執行順序（Phase）

> 每個 Phase 結束都是獨立可展示的狀態，不依賴下一個 Phase。

| Phase | 內容                              | 狀態      |
|-------|-----------------------------------|-----------|
| 1     | Vue Dashboard + Mock 資料         | 進行中    |
| 2     | DB Schema + Go REST API           | 待開始    |
| 3     | WebSocket 即時推送                | 待開始    |
| 4     | Python SPC 告警引擎               | 待開始    |
| 5     | Lot 派工流程                      | 待開始    |
| 6     | Docker Compose 整合 + Demo 準備   | 待開始    |

**最低可展示目標：Phase 1–3 完成。**

---

## 前端畫面規劃（已確認）

### Equipment 頁面設計邏輯

**核心問題：設備高達 100 台，卡片式排列失控。**

採用「群組折疊 + 列表行」設計：

#### 整體佈局
- 頂部 Header：Logo、導覽列、Alarm 快捷按鈕、即時時鐘
- KPI 列：整體稼動率、Running / Idle / Down / PM / Alarm 數量
- 工具列：搜尋框 + 類型篩選 Pill + 「僅看 Alarm」快篩
- 群組折疊區：依設備類型（CVD / Etch / CMP / Diffusion）分群

#### 群組折疊規則
- 每個群組 Header 收合時顯示：設備類型名、各狀態數量小圓點、稼動率橫條
- 有 Alarm 的群組：Header 邊框變紅 + 顯示「N Alarm」紅色標籤
- 群組展開後為表格列表，欄位：設備名 / 狀態 / Current Lot / 溫度 / 壓力 / 稼動率
- **Alarm 設備永遠排在群組最頂部**，左側有紅色邊框標記

#### 設備詳細 Drawer
- 點擊任一設備列 → 右側滑出 Drawer（不跳頁）
- Drawer 內容：設備名稱、狀態 Badge、Current Lot + Recipe、製程參數（溫度/壓力含 UCL/LCL）、SPC 迷你趨勢圖（最近 20 點）、操作按鈕（Hold 設備 / 查看歷史）

#### 篩選行為
- 「僅看 Alarm」：隱藏無 Alarm 的群組，Alarm 設備直接展開
- 類型篩選（CVD / Etch…）：只顯示該類型群組
- 搜尋：即時 filter 設備名稱或 Lot 編號

---

## 產業術語對照（程式碼與 UI 必須使用）

| 術語        | 說明                                         |
|-------------|----------------------------------------------|
| Equipment   | 設備（CVD / Etch / CMP / Diffusion）         |
| Lot         | 批次，一批 wafer 的生產單位                  |
| Wafer       | 晶圓，Lot 內最小追蹤單位                     |
| Recipe      | 製程配方，定義溫度 / 壓力 / 時間參數         |
| WIP         | Work In Progress，在製品                     |
| SPC         | Statistical Process Control，統計製程管制    |
| UCL / LCL   | Upper / Lower Control Limit，管制上下限      |
| Dispatch    | 派工，將 Lot 指派給可用設備                  |
| Downtime    | 設備停機時間                                 |
| PM          | Preventive Maintenance，預防性保養           |
| Hold        | 暫停該設備或 Lot 的生產                      |

---

## 資料模型（核心 Schema）

```sql
-- 設備
Equipment: id, name, type, status(RUNNING|IDLE|DOWN|PM), current_lot_id, updated_at

-- 批次
Lot: id, lot_number, product, recipe_id, priority, status(QUEUED|RUNNING|COMPLETED|ON_HOLD), wafer_count, created_at

-- 晶圓
Wafer: id, lot_id, sequence, status, defect_flag

-- 製程配方
Recipe: id, name, equipment_type, target_temp, target_pressure, duration_min

-- SPC 紀錄
SPC_Record: id, equipment_id, parameter, value, ucl, lcl, is_alarm, timestamp
```

---

## TypeScript 型別定義（frontend/src/types/mes.ts）

```typescript
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
```

---

## API 介面約定

所有 Go API 前綴：`/api/v1/`

```
GET    /api/v1/equipment                    設備列表（含狀態）
GET    /api/v1/equipment/:id                單台設備詳細
PUT    /api/v1/equipment/:id/status         更新設備狀態
POST   /api/v1/equipment/:id/hold           Hold 設備

GET    /api/v1/lots                         Lot 列表
POST   /api/v1/lots                         建立新 Lot
GET    /api/v1/lots/:id                     單一 Lot 詳細
POST   /api/v1/lots/:id/dispatch            觸發派工

GET    /api/v1/alarms                       告警列表
PUT    /api/v1/alarms/:id/acknowledge       確認告警

GET    /api/v1/spc/:equipment_id            SPC 歷史資料

WS     /ws                                  WebSocket 連線
```

### WebSocket 事件格式
```json
{
  "event": "equipment_status_changed",
  "payload": { "equipment_id": 3, "status": "DOWN" }
}
{
  "event": "spc_alarm",
  "payload": { "equipment_id": 3, "parameter": "temperature", "value": 692.4, "ucl": 680 }
}
{
  "event": "lot_dispatched",
  "payload": { "lot_id": 17, "equipment_id": 5 }
}
```

---

## 各服務 Port

| 服務            | Port  |
|-----------------|-------|
| Vue dev server  | 5173  |
| Go API          | 8080  |
| Python FastAPI  | 8001  |
| PostgreSQL      | 5432  |
| Redis           | 6379  |

---

## 程式碼規範

### Vue（frontend/）
- Composition API + `<script setup lang="ts">` 語法
- 型別定義統一放 `src/types/mes.ts`，不在元件內重複定義
- Pinia store 每個 domain 一個檔案
- API 呼叫封裝在 `src/api/`，元件不直接呼叫 axios
- 元件命名 PascalCase、composable 用 useXxx、store 用 useXxxStore
- Element Plus 元件優先，不重複造輪子

### Go（backend/）
- 標準三層：handler → service → repository
- 錯誤統一回傳 `{"error": "訊息"}` JSON 格式
- WebSocket Hub 統一在 `pkg/websocket/hub.go`
- 環境變數透過 `.env` 注入，不 hardcode 連線字串

### Python（analytics/）
- FastAPI 對外提供 HTTP endpoint，Go 透過 HTTP 呼叫
- SPC 計算邏輯封裝在 `spc/control_chart.py`
- 不與 Go 共用 DB 連線，透過 API 解耦

### SQL（migrations/）
- 檔名格式：`001_init.sql`、`002_add_xxx.sql`
- 每個 migration 使用 `IF NOT EXISTS`，可重複執行

---

## Demo 腳本（展示用）

1. 開啟 Dashboard，展示 100 台設備群組折疊視圖
2. 點擊「僅看 Alarm」篩選，快速定位異常設備
3. 點入異常設備，右側 Drawer 展開，說明 SPC 超標細節
4. 切換到 Lot 頁面，建立新 Lot 並觸發自動派工
5. 設備狀態即時透過 WebSocket 更新（畫面動起來）
6. 口頭補充：「實際環境這裡會接 SECS/GEM 協定或 OPC-UA，SPC 引擎可以接真實製程資料」

---

## 給 Claude CLI 的開發指令範例

```bash
# 開始新任務時，請先說明目前在哪個 Phase
> 現在在 Phase 1，幫我實作 frontend/src/components/equipment/EquipmentGroup.vue
  依照 CLAUDE.md 的設計規範，群組 Header 要顯示類型名稱、各狀態數量、稼動率橫條

> Phase 2 開始，依照 CLAUDE.md 的 Schema 建立 migrations/001_init.sql

> 依照 CLAUDE.md 的 API 規範，實作 backend/internal/handler/equipment.go 的 listEquipment handler
```