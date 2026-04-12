# Mini-MES

半導體製造執行系統（MES）展示專案。

展示 MES 領域知識 + 全端開發能力 + 系統設計思維，Demo 時間約 3–5 分鐘。

---

## 技術棧

| 層級 | 技術 |
|------|------|
| 前端 | Vue 3 + TypeScript + Element Plus + Vite |
| 後端 API | Go + Gin + WebSocket (gorilla/websocket) |
| 分析引擎 | Python + FastAPI + numpy |
| 資料庫 | PostgreSQL 16 |
| 容器化 | Docker Compose |

---

## 功能特色

### 設備監控 Dashboard
- 100 台設備依類型（CVD / Etch / CMP / Diffusion）群組折疊顯示
- 群組 Header 顯示各狀態數量（Running / Idle / Down / PM）與整體稼動率橫條
- 有 Alarm 的群組：Header 邊框變紅 + 閃爍 Alarm 標籤
- KPI 列：整體稼動率、各狀態設備數量即時統計

### 設備詳細 Drawer
- 點擊任一設備列 → 右側滑出 Drawer
- 顯示設備名稱、狀態、Current Lot + Recipe
- 製程參數（溫度 / 壓力）含 UCL/LCL 管制線
- SPC 迷你趨勢圖（最近 20 筆量測值）
- 操作按鈕：Hold 設備 / 查看完整 SPC

### 多維篩選
- 關鍵字搜尋設備名稱 / Lot 編號（即時 filter）
- 類型篩選（CVD / Etch / CMP / Diffusion）
- 「僅看 Alarm」快篩，有 Alarm 群組自動展開

### Lot / WIP 看板
- Kanban 風格顯示在製品各階段（排隊中 / Running / On Hold / 完成）
- 建立新 Lot：填入 Lot 編號、產品、Recipe、優先度、Wafer 數
- 派工按鈕：自動指派可用設備，觸發 WebSocket 即時通知

### Python SPC 告警引擎
- Nelson Rules 實作：
  - Rule 1（CRITICAL）：最新量測值超出 UCL / LCL
  - Rule 2（WARNING）：連續 9 點落在平均同側
  - Rule 3（WARNING）：連續 6 點單調遞增或遞減
- Go 後端模擬器每 8 秒呼叫 Python 分析端點，結果寫入 DB 並廣播到前端

### WebSocket 即時推送
- 設備狀態變更（`equipment_status_changed`）
- SPC 告警觸發（`spc_alarm`）
- Lot 派工完成（`lot_dispatched`）

---

## 開發進度（Phase）

| Phase | 內容 | 狀態 |
|-------|------|------|
| 1 | Vue Dashboard + Mock 資料 | ✅ 完成 |
| 2 | DB Schema + Go REST API | ✅ 完成 |
| 3 | WebSocket 即時推送 | ✅ 完成 |
| 4 | Python SPC 告警引擎 | ✅ 完成 |
| 5 | Lot 派工流程 | ✅ 完成 |
| 6 | Docker Compose 整合 + Demo 準備 | ✅ 完成 |

---

## 快速啟動（Docker Compose）

```bash
# 全服務啟動（postgres → analytics → backend → frontend）
docker compose up -d

# 確認所有服務健康
docker compose ps

# 開啟瀏覽器
# http://localhost:5173
```

> **注意：** 首次啟動約需 30–60 秒等待 PostgreSQL 初始化與 Seed 資料載入。

### 停止並清除資料

```bash
docker compose down -v
```

---

## 本地開發啟動

### 前置需求
- Node.js 18+
- Go 1.25+
- Python 3.12+
- PostgreSQL 16（或透過 Docker 只起 DB）

### 步驟

```bash
# 1. 啟動 PostgreSQL
docker compose up -d postgres

# 2. 啟動 Python 分析引擎
cd analytics
pip install -r requirements.txt
uvicorn main:app --port 8001

# 3. 啟動 Go 後端（新 terminal）
cd backend
go run cmd/server/main.go

# 4. 啟動 Vue 前端（新 terminal）
cd frontend
npm install
npm run dev
# 開啟 http://localhost:5173
```

---

## 服務 Port 對照

| 服務 | Port |
|------|------|
| Vue 前端 | 5173 |
| Go API | 8080 |
| Python 分析引擎 | 8001 |
| PostgreSQL | 5432 |

---

## API 端點概覽

```
GET    /api/v1/equipment               設備列表（含狀態、當前 Lot）
GET    /api/v1/equipment/:id           單台設備詳細
PUT    /api/v1/equipment/:id/status    更新設備狀態
POST   /api/v1/equipment/:id/hold      Hold 設備

GET    /api/v1/lots                    Lot 列表
POST   /api/v1/lots                    建立新 Lot
POST   /api/v1/lots/:id/dispatch       觸發派工

GET    /api/v1/alarms                  告警列表
PUT    /api/v1/alarms/:id/acknowledge  確認告警

GET    /api/v1/spc/:equipment_id       SPC 歷史資料
GET    /api/v1/recipes                 Recipe 列表

WS     /ws                             WebSocket 即時事件
```

---

## 目錄結構

```
mini-mes/
├── frontend/           # Vue 3 前端
│   ├── src/
│   │   ├── components/ # EquipmentGroup / AlarmList / SpcMiniChart ...
│   │   ├── stores/     # Pinia store（equipment / alarm / lot）
│   │   ├── api/        # axios 封裝
│   │   ├── composables/# useWebSocket
│   │   └── views/      # DashboardView / LotView / SpcView
│   └── nginx.conf      # 反向代理設定（Docker 部署用）
├── backend/            # Go API 後端
│   ├── cmd/server/     # 程式進入點
│   ├── internal/
│   │   ├── handler/    # HTTP 路由層
│   │   ├── service/    # 業務邏輯層
│   │   ├── repository/ # DB 存取層
│   │   └── model/      # 資料結構定義
│   ├── pkg/websocket/  # WebSocket Hub
│   └── internal/simulator/ # 產線模擬器（Demo 用）
├── analytics/          # Python SPC 分析引擎
│   ├── main.py         # FastAPI 服務
│   └── spc/            # Nelson Rules 實作
├── migrations/
│   ├── 001_init.sql    # DB Schema
│   └── seed.sql        # 初始資料（100 台設備 / Recipe / Lot）
├── docs/
│   ├── api-spec.md     # API 規格文件
│   └── db-schema.md    # DB Schema 說明
└── docker-compose.yml
```

---

## MES 產業術語

| 術語 | 說明 |
|------|------|
| Equipment | 設備（CVD / Etch / CMP / Diffusion） |
| Lot | 批次，一批 wafer 的生產單位 |
| Wafer | 晶圓，Lot 內最小追蹤單位 |
| Recipe | 製程配方，定義溫度 / 壓力 / 時間參數 |
| WIP | Work In Progress，在製品 |
| SPC | Statistical Process Control，統計製程管制 |
| UCL / LCL | Upper / Lower Control Limit，管制上下限 |
| Dispatch | 派工，將 Lot 指派給可用設備 |
| Hold | 暫停設備或 Lot 的生產 |
| PM | Preventive Maintenance，預防性保養 |
| Downtime | 設備停機時間 |

---

## Demo 腳本（3–5 分鐘）

1. 開啟 Dashboard，展示 100 台設備群組折疊視圖
2. 點擊「僅看 Alarm」快篩，快速定位 SPC 異常設備
3. 點入異常設備，右側 Drawer 展開，說明 SPC 超標細節與趨勢圖
4. 切換到 Lot 頁面，建立新 Lot 並觸發自動派工
5. 觀察設備狀態即時透過 WebSocket 更新（Dashboard 畫面動起來）
6. 補充說明：「實際環境這裡會接 SECS/GEM 協定或 OPC-UA，SPC 引擎可以接真實製程資料」
