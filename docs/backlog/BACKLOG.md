# Mini-MES 技術優化待辦清單

> 依據 2026-04-20 代碼審查結果整理。
> 優先順序：🔴 Critical → 🟠 High → 🟡 Medium

---

## 第一批（本週）— 影響正確性

### 後端

- [ ] 🔴 **Dispatch Race Condition**
  - 檔案：`backend/internal/service/dispatch.go`
  - 問題：`AssignLot` 與 `UpdateStatus` 兩步操作無 Transaction 保護，第二步失敗時設備被佔用但 Lot 仍為 QUEUED
  - 修正：使用 `sqlx.BeginTxx` 包裝成單一 Transaction

- [ ] 🔴 **ClearLot 邏輯矛盾**
  - 檔案：`backend/internal/repository/equipment.go`
  - 問題：`ClearLot` 強制將設備 `status` 設為 `IDLE`，與 Hold 設備（設為 `DOWN`）邏輯衝突
  - 修正：`ClearLot` 只清 `current_lot_id`，不更動 `status`

- [ ] 🟠 **缺少 Gin Middleware**
  - 檔案：`backend/cmd/server/main.go`
  - 問題：未掛載 `gin.Logger()` 與 `gin.Recovery()`，panic 不會自動恢復
  - 修正：在 router 初始化後加入兩個 middleware

- [ ] 🟠 **Status 欄位缺少 ENUM 驗證**
  - 檔案：`backend/internal/handler/equipment.go`
  - 問題：`UpdateStatus` 接受任意字串，未驗證是否為合法值
  - 修正：在 handler 或 service 層加入 ENUM 白名單驗證

### 前端

- [ ] 🔴 **型別命名不一致（snake_case vs camelCase）**
  - 檔案：`frontend/src/types/mes.ts`
  - 問題：`ucl_temp`、`lcl_temp`、`ucl_pressure`、`lcl_pressure` 使用 snake_case，與其他欄位不一致
  - 修正：統一改為 `uclTemp`、`lclTemp`、`uclPressure`、`lclPressure`，同步更新 `api/equipment.ts` 的 mapping

---

## 第二批（下週）— 穩定性與正確性

### 前端

- [ ] 🟡 **package.json 未宣告 engines 欄位**
  - 檔案：`frontend/package.json`
  - 問題：`vitest 4.x` 要求 Node ≥ 20，但未在 `engines` 宣告，使用 Node 18 時只會出現 npm warn 而非明確錯誤
  - 修正：加入 `"engines": { "node": ">=20.0.0" }`

- [ ] 🟠 **Store API 無重試機制**
  - 檔案：`frontend/src/stores/equipment.ts`、`lot.ts`、`alarm.ts`
  - 問題：API 失敗直接報錯，短暫網路抖動會導致畫面掛掉
  - 修正：安裝 `axios-retry`，或在各 store 的 fetch 函式加入指數退避重試（最多 3 次）

- [ ] 🟠 **WebSocket 事件 payload 無型別驗證**
  - 檔案：`frontend/src/composables/useWebSocket.ts`
  - 問題：`spc_alarm`、`equipment_status_changed`、`lot_dispatched` 的 payload 為 `any`，無執行時驗證
  - 修正：安裝 `zod`，為每個 WebSocket 事件定義 schema 並在收到事件時解析驗證

### 資料庫

- [ ] 🟠 **缺少複合索引**
  - 檔案：`migrations/001_init.sql`（或新增 migration）
  - 問題：常用查詢 `WHERE status = 'RUNNING' AND type = 'CVD'` 無複合索引支援
  - 修正：新增 `003_add_compound_indices.sql`
    ```sql
    CREATE INDEX IF NOT EXISTS idx_equipment_status_type ON equipment(status, type);
    CREATE INDEX IF NOT EXISTS idx_alarm_acknowledged_ts ON alarm_event(acknowledged, timestamp DESC);
    ```

### Python 分析引擎

- [ ] 🟠 **Nelson Rule 2/3 中心線計算有誤**
  - 檔案：`analytics/spc/control_chart.py`
  - 問題：Rule 2/3 用計算出的樣本 `mean` 判斷，標準做法應用 `(UCL + LCL) / 2` 作為製程目標中心線
  - 修正：將 `mean = float(arr.mean())` 改為 `center = (ucl + lcl) / 2`，並以 `center` 判斷同側

- [ ] 🟡 **缺少輸入驗證**
  - 檔案：`analytics/main.py`
  - 問題：未驗證 `ucl > lcl`、`values` 不為空
  - 修正：在 endpoint 加入 Pydantic validator 或手動 guard

- [ ] 🟡 **pytest 未列入 requirements.txt**
  - 檔案：`analytics/requirements.txt`
  - 問題：`pytest` 未記錄，新環境建立 `.venv` 後無法直接執行測試
  - 修正：新增 `requirements-dev.txt` 並加入 `pytest`，或於 README 補充安裝指令

---

## 第三批（後續）— 生產化與效能

### 安全性

- [ ] 🔴 **無任何認證機制**
  - 影響：所有 API 端點公開可存取
  - 建議：後端加入 JWT middleware（`github.com/golang-jwt/jwt/v5`），前端存 token 於 `httpOnly cookie`

- [ ] 🟠 **CORS 設定過寬**
  - 檔案：`backend/cmd/server/main.go`
  - 問題：`Access-Control-Allow-Origin: *` 在生產環境不安全
  - 修正：限制為部署的前端 origin

- [ ] 🟠 **無 Rate Limiting**
  - 建議：加入 `github.com/ulule/limiter/v3` 限制每 IP 的請求頻率

### 效能

- [ ] 🟡 **100 台設備列表無虛擬滾動**
  - 檔案：`frontend/src/components/equipment/EquipmentGroupList.vue`
  - 問題：展開所有群組時 DOM 一次掛載 100 筆，拖慢渲染
  - 修正：安裝 `vue-virtual-scroller`，對展開的設備列表啟用虛擬滾動

- [ ] 🟡 **Python SPC 無快取**
  - 問題：相同資料序列每次都重新計算，高頻調用浪費算力
  - 修正：接入 Redis，以 `hash(values + ucl + lcl)` 為 key 快取分析結果（TTL 60s）

- [ ] 🟡 **資料庫連線池設定偏低**
  - 檔案：`backend/cmd/server/main.go`
  - 問題：`MaxOpenConns: 20`、`MaxIdleConns: 5` 比例不佳
  - 修正：調整為 `MaxOpenConns: 50`、`MaxIdleConns: 15`

### 可觀測性

- [ ] 🟡 **缺少結構化日誌**
  - 建議：後端引入 `go.uber.org/zap`，將 `log.Printf` 替換為結構化欄位輸出

- [ ] 🟡 **Wafer 表從未使用**
  - 檔案：`migrations/001_init.sql`
  - 問題：建立了 `wafer` 表但整個後端無任何讀寫
  - 決策：若不在 Phase 6 以內實作，建議移除或加上 TODO 說明

---

## 套件升級建議

### Go 後端

```bash
go get github.com/golang-jwt/jwt/v5            # JWT 認證
go get github.com/ulule/limiter/v3             # Rate limiting
go get go.uber.org/zap                         # 結構化日誌
go get github.com/go-playground/validator/v10  # 參數驗證
```

### Vue 前端

```bash
npm install zod            # API / WebSocket payload 執行時驗證
npm install @vueuse/core   # useDebounce、useLocalStorage 等實用 composable
npm install axios-retry    # 自動重試
npm install vue-virtual-scroller  # 虛擬滾動
```

### Python

```bash
pip install redis            # SPC 結果快取
pip install pydantic-settings  # 環境變數型別安全讀取
```

---

## 架構評分摘要

| 面向 | 評分 | 主要問題 |
|------|------|----------|
| 分層設計 | ⭐⭐⭐⭐⭐ | — |
| 設備狀態機 | ⭐⭐⭐⭐⭐ | — |
| 前端組件化 | ⭐⭐⭐⭐ | — |
| 錯誤處理 | ⭐⭐ | 缺 Transaction、重試、Recovery middleware |
| 安全性 | ⭐ | 無認證（Demo 可接受，正式環境必修） |
| 效能 | ⭐⭐⭐ | 索引不完整，無快取層 |
