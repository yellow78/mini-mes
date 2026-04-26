# 2026-04-26 — 測試環境建置與開發規範整理

## 一、前端測試修正（Vitest）

**問題根因與解法：**

| 問題 | 原因 | 解法 |
|------|------|------|
| vitest 啟動崩潰 | vitest 4.x 需 Node ≥ 20，環境為 v18 | `nvm install 20 && nvm use 20` |
| 測試打真實後端 | Store fetch 未 mock，回 ECONNREFUSED | `vi.mock()` 攔截 API 模組 |
| fetch 後資料為空 | async fetch 未 await 就斷言 | 測試加 `async/await` |
| 重置測試失敗 | `mockResolvedValue` 回傳同一物件參考，被 store mutate | 改用 `mockImplementation` 每次回傳新複本 |
| 呼叫不存在方法 | Phase 2 刪除 `startSimulation`，測試未同步 | 移除相關測試案例 |

---

## 二、分析引擎測試建立（pytest）

- 新增 `analytics/tests/test_control_chart.py`，28 個測試全數通過
- 測試分組：`TestBasic` / `TestRule1` / `TestRule2` / `TestRule3` / `TestMultipleRules`
- Rule 2 mock 技巧：用少數異常值把樣本均值拉偏，使目標 9 點落在同側

**待修（已記入 BACKLOG）：**
- `pytest` 未列入 `requirements.txt`
- Rule 2/3 應改用 `(UCL+LCL)/2` 作為中心線，目前用樣本均值

---

## 三、slash commands 建立

在 `.claude/commands/test/` 建立三個自訂指令：

| 指令 | 檔案 | 用途 |
|------|------|------|
| `/test:frontend` | `test/frontend.md` | 執行 Vitest 前端測試 |
| `/test:backend` | `test/backend.md` | 執行 Go 後端測試 |
| `/test:analytics` | `test/analytics.md` | 執行 pytest 分析引擎測試 |

指令會自動處理環境切換、回報通過/失敗數量與錯誤訊息。

---

## 四、開發規範整理（CLAUDE.md）

**分支命名：** `<type>-<scope>-<description>`，單字用 `_` 連接
```
test-frontend_analytics-store_and_spc
fix-backend-ws_reconnect
```

**Commit 訊息：** `<type>:[scope] <description>`，scope 可並列
```
test:[frontend][analytics] 修正前端 store 測試並新增 SPC 分析引擎測試
fix:[backend] 修正 WebSocket 斷線重連邏輯
```

---

## 五、AI 協同觀察

- 環境問題（Node 版本、套件缺失）需人工確認方向後 AI 執行
- 測試 mock 資料設計需理解業務邏輯，AI 可產出但需驗證
- 規範文件（CLAUDE.md）由人主導意圖，AI 負責格式整理與補全
