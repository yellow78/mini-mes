# Mini-MES

半導體製造執行系統（MES）。

展示 MES 領域知識 + 全端開發能力 + 系統設計思維。

---

## 技術棧

| 層級 | 技術 |
|------|------|
| 前端 | Vue 3 + TypeScript + Element Plus + Vite |
| 後端 API | Go + Gin + WebSocket (gorilla) |
| 分析引擎 | Python + FastAPI + pandas |
| 資料庫 | PostgreSQL + Redis |
| 容器化 | Docker Compose |

---

## 功能特色

- **設備監控 Dashboard**：100 台設備依類型（CVD / Etch / CMP / Diffusion）群組折疊顯示
- **即時告警**：SPC 超標設備自動置頂，紅色邊框 + 閃爍標籤
- **設備詳細 Drawer**：點擊設備列滑出右側抽屜，顯示製程參數（含 UCL/LCL）與 SPC 趨勢圖
- **多維篩選**：搜尋設備名 / Lot 編號、類型篩選、「僅看 Alarm」快篩
- **Lot / WIP 看板**：Kanban 風格顯示在製品流程狀態
- **WebSocket 即時推送**：設備狀態變更、SPC 告警事件即時更新畫面

---

## 開發進度（Phase）

| Phase | 內容 | 狀態 |
|-------|------|------|
| 1 | Vue Dashboard + Mock 資料 | ✅ 完成 |
| 2 | DB Schema + Go REST API | 🔲 待開始 |
| 3 | WebSocket 即時推送 | 🔲 待開始 |
| 4 | Python SPC 告警引擎 | 🔲 待開始 |
| 5 | Lot 派工流程 | 🔲 待開始 |
| 6 | Docker Compose 整合 | 🔲 待開始 |

---

## 本地啟動（Phase 1）

```bash
cd frontend
npm install
npm run dev
# 開啟 http://localhost:5173
```

---

## 目錄結構

```
mini-mes/
├── frontend/       # Vue 3 前端
├── backend/        # Go API 後端
├── analytics/      # Python SPC 分析引擎
├── migrations/     # SQL Schema
└── docs/           # API 規格文件
```

---

## MES 產業術語

| 術語 | 說明 |
|------|------|
| Equipment | 設備（CVD / Etch / CMP / Diffusion） |
| Lot | 批次，一批 wafer 的生產單位 |
| Recipe | 製程配方，定義溫度 / 壓力 / 時間參數 |
| SPC | Statistical Process Control，統計製程管制 |
| UCL / LCL | Upper / Lower Control Limit，管制上下限 |
| WIP | Work In Progress，在製品 |
| Dispatch | 派工，將 Lot 指派給可用設備 |
| PM | Preventive Maintenance，預防性保養 |
