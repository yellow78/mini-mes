---
name: test-analytics
description: 執行 Python 分析引擎單元測試（pytest），可選傳入測試路徑或關鍵字
---

在 analytics/ 目錄執行 Python 單元測試。

步驟：
1. 確認 pytest 已安裝（若無則提示：`pip install pytest`）
2. 從專案根目錄進入 analytics/
3. 執行測試指令：
   - 若有 $ARGUMENTS：`cd analytics && python -m pytest $ARGUMENTS -v --ignore=.venv`
   - 若無 $ARGUMENTS（預設）：`cd analytics && python -m pytest -v --ignore=.venv`
4. 回報結果：
   - 列出通過 / 失敗的測試數量與檔案
   - 失敗時顯示測試名稱與錯誤訊息
   - 若尚無測試檔案，提示使用者在 analytics/ 下建立 test_*.py 檔案
   - 全部通過時顯示簡短成功摘要
