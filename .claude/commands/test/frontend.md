---
name: test-frontend
description: 執行 Vue 前端單元測試（Vitest），可選傳入檔案路徑或關鍵字過濾
---

在 frontend/ 目錄執行 Vitest 單元測試。

步驟：
1. 從專案根目錄進入 frontend/
2. 執行測試指令：
   - 若有 $ARGUMENTS：`cd frontend && npx vitest run $ARGUMENTS`
   - 若無 $ARGUMENTS（預設）：`cd frontend && npx vitest run`
3. 回報結果：
   - 列出通過 / 失敗的測試數量與檔案
   - 失敗時顯示測試名稱與錯誤訊息
   - 全部通過時顯示簡短成功摘要
