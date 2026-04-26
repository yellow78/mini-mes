---
name: test-backend
description: 執行 Go 後端單元測試，可選傳入 package 路徑（例如 ./internal/service/...）
---

在 backend/ 目錄執行 Go 單元測試。

步驟：
1. 從專案根目錄進入 backend/
2. 執行測試指令：
   - 若有 $ARGUMENTS：`cd backend && go test $ARGUMENTS -v`
   - 若無 $ARGUMENTS（預設）：`cd backend && go test ./... -v`
3. 回報結果：
   - 列出通過 / 失敗的測試數量
   - 失敗時顯示測試名稱與錯誤訊息
   - 全部通過時顯示簡短成功摘要
