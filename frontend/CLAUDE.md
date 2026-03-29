# Mini-MES Frontend — Claude CLI 前端開發手冊

> 全域資訊（專案定位、Phase 順序、產業術語、Port 對照）請參考根目錄 `../CLAUDE.md`

## Claude 互動規範
- 回答語言：繁體中文
- 程式碼註解：繁體中文
- 變數、函式、檔案名稱：英文

---

## 技術棧
- Vue 3 + TypeScript + `<script setup>` 語法
- Element Plus（UI 元件庫）
- Pinia（狀態管理）
- Vue Router 4
- Axios（API 呼叫）
- ECharts + vue-echarts（SPC 趨勢圖）
- Vite（開發工具）
- **Vitest + Vue Test Utils**（單元測試）

---

## 目錄結構

```
frontend/src/
├── __tests__/
│   ├── stores/
│   │   ├── equipment.test.ts        ← Store 邏輯測試
│   │   └── alarm.test.ts
│   ├── components/
│   │   ├── equipment/
│   │   │   ├── EquipmentGroup.test.ts
│   │   │   └── EquipmentRow.test.ts
│   │   └── alarm/
│   │       └── AlarmList.test.ts
│   ├── composables/
│   │   └── useWebSocket.test.ts
│   └── api/
│       └── equipment.test.ts        ← API 層 mock 測試
├── types/
│   └── mes.ts                   ← 所有型別集中定義，不在元件內重複
├── stores/
│   ├── equipment.ts             ← 設備狀態、群組計算
│   ├── alarm.ts                 ← 告警列表
│   └── lot.ts                   ← Lot / WIP 狀態
├── api/
│   ├── equipment.ts             ← axios 封裝，元件不直接呼叫 axios
│   ├── lot.ts
│   └── alarm.ts
├── composables/
│   └── useWebSocket.ts          ← WS 連線管理，統一訂閱事件
├── components/
│   ├── layout/
│   │   ├── AppHeader.vue        ← 頂部導覽、時鐘、Alarm 快捷
│   │   └── AppSidebar.vue       ← 左側 icon 導覽
│   ├── equipment/
│   │   ├── EquipmentGroupList.vue   ← 群組折疊主容器
│   │   ├── EquipmentGroup.vue       ← 單一群組（含折疊 Header + 列表）
│   │   ├── EquipmentRow.vue         ← 單台設備列（表格行）
│   │   └── EquipmentDrawer.vue      ← 點入後右側詳細抽屜
│   ├── alarm/
│   │   └── AlarmList.vue
│   └── spc/
│       └── SpcMiniChart.vue         ← Drawer 內的迷你趨勢圖
├── views/
│   ├── DashboardView.vue        ← Equipment 主頁
│   ├── LotView.vue              ← Lot / WIP 看板
│   └── SpcView.vue              ← SPC 完整趨勢圖頁
└── router/
    └── index.ts
```

---

## 命名規範

| 類型        | 規範                          | 範例                        |
|-------------|-------------------------------|-----------------------------|
| 元件        | PascalCase                    | `EquipmentGroup.vue`        |
| Composable  | useXxx                        | `useWebSocket`              |
| Store       | useXxxStore                   | `useEquipmentStore`         |
| API 函式    | 動詞 + 名詞                   | `fetchEquipments`, `holdEquipment` |
| 型別        | PascalCase interface          | `Equipment`, `AlarmEvent`   |
| CSS class   | kebab-case                    | `eq-row`, `alarm-tag`       |

---

## Element Plus 使用規範

- 優先使用 Element Plus 元件，不重複造輪子
- 深色主題：在 `main.ts` 引入 `element-plus/theme-chalk/dark/css-vars.css`
- CSS Token 覆蓋統一在 `src/style.css`，不在元件內 hardcode 顏色
- 常用元件對應：
  - 群組折疊 → `el-collapse` 或自製（推薦自製，彈性更高）
  - 設備 Drawer → `el-drawer`（direction="rtl"）
  - 搜尋框 → `el-input` + prefix icon
  - 狀態 Badge → `el-tag`（type: success / info / danger / warning）
  - SPC 趨勢圖 → ECharts（非 Element Plus，另外引入）
  - 篩選 Pill → `el-radio-group` + `el-radio-button`

---

## 主題色彩定義（src/style.css）

```css
:root {
  /* Element Plus Token 覆蓋 */
  --el-bg-color: #0f172a;
  --el-bg-color-page: #0f172a;
  --el-bg-color-overlay: #1e293b;
  --el-fill-color-blank: #1e293b;
  --el-border-color: #334155;
  --el-border-color-light: #1e293b;
  --el-text-color-primary: #e2e8f0;
  --el-text-color-regular: #94a3b8;
  --el-color-primary: #3b82f6;
}

/* MES 設備狀態專用色 */
:root {
  --mes-running: #22c55e;
  --mes-idle: #94a3b8;
  --mes-down: #ef4444;
  --mes-pm: #f59e0b;
  --mes-alarm: #ef4444;
  --mes-surface: #1e293b;
  --mes-surface-deep: #0f172a;
  --mes-border: #334155;
}
```

---

## Equipment 頁面設計規範（已確認）

### 核心問題
設備高達 100 台，卡片式排列資訊密度不足。
**解法：群組折疊 + 列表行**

### 整體佈局
```
Header（頂部）
  ├── Logo + 導覽列
  ├── Alarm 快捷按鈕（紅色，顯示數量）
  └── 即時時鐘

KPI 列
  └── 稼動率 / Running / Idle / Down / PM / Alarm 數量

工具列
  ├── 搜尋框（設備名稱 / Lot 編號即時 filter）
  ├── 類型篩選 Pill（全部 / CVD / Etch / CMP / Diffusion）
  └── 「僅看 Alarm」快篩 Pill

群組折疊區（EquipmentGroupList）
  ├── CVD 群組（EquipmentGroup）
  ├── Etch 群組
  ├── CMP 群組
  └── Diffusion 群組
```

### 群組 Header（收合狀態）顯示
- 設備類型名稱（CVD / Etch…）
- 各狀態小圓點數量（Running綠 / Idle灰 / Down紅 / PM黃）
- 稼動率橫條
- 若有 Alarm：Header 邊框變紅 + 顯示「N Alarm」紅色標籤（閃爍）

### 群組展開後 表格欄位
| 欄位       | 寬度   | 說明                          |
|------------|--------|-------------------------------|
| 設備名     | 80px   | font-mono，如 CVD-03          |
| 狀態       | 100px  | el-tag badge                  |
| Current Lot| flex   | lot_number，無則顯示 —        |
| 溫度       | 90px   | °C，超標顯示紅色 + ↑          |
| 壓力       | 90px   | mTorr，超標顯示紅色 + ↑       |
| 稼動率     | 70px   | 橫條 + 百分比                 |

### Alarm 設備排序規則
1. Alarm 設備永遠置頂（不管其他欄位排序）
2. Alarm 列：左側 2px 紅色邊框
3. 超標參數：紅色文字 + 上箭頭符號 ↑

### EquipmentDrawer 內容
點擊設備列後，右側滑出 `el-drawer`（寬度 320px）：
- 設備名稱 + 類型 + 狀態 Badge
- Current Lot + Recipe 名稱
- 製程參數卡片：溫度（含 UCL/LCL）、壓力（含 UCL/LCL）
- 稼動率橫條
- SpcMiniChart（最近 20 筆，顯示 UCL/LCL 管制線）
- 操作按鈕：Hold 設備（紅色）/ 查看完整 SPC（藍色）

### 篩選行為
- 「僅看 Alarm」：隱藏無 Alarm 的群組，有 Alarm 群組自動展開
- 類型篩選：只顯示對應類型的群組
- 搜尋：跨群組 filter，符合的設備所在群組自動展開，不符合的群組隱藏

---

## Pinia Store 設計

### useEquipmentStore
```typescript
// 核心 state
equipments: Equipment[]

// getters
groupedEquipments: EquipmentGroup[]   // 依 type 分組，alarm 置頂
filteredGroups: EquipmentGroup[]      // 套用搜尋 + 篩選後的結果
overallUtilization: number
statusCount: { running, idle, down, pm, alarm }

// actions
fetchEquipments()                     // 呼叫 GET /api/v1/equipment
updateEquipmentStatus(id, status)     // WebSocket 事件觸發時呼叫
setFilter(type, alarmOnly, keyword)   // 更新篩選條件
```

### useAlarmStore
```typescript
alarms: AlarmEvent[]
unacknowledgedCount: number

fetchAlarms()
acknowledgeAlarm(id)
addAlarm(alarm)                       // WebSocket spc_alarm 事件觸發
```

---

## WebSocket 整合

```typescript
// composables/useWebSocket.ts
// 訂閱事件，分發到對應 store
ws.onmessage = (msg: WSMessage) => {
  switch (msg.event) {
    case 'equipment_status_changed':
      equipmentStore.updateEquipmentStatus(...)
      break
    case 'spc_alarm':
      alarmStore.addAlarm(...)
      break
    case 'lot_dispatched':
      lotStore.updateLotStatus(...)
      break
  }
}
```

---

## Mock 資料規範（Phase 1 使用）

- Mock 資料統一放在各 store 的初始 state
- 模擬資料要包含：至少 2 台 Alarm 設備、各類型設備各 4 台以上
- 每 2 秒微幅波動溫度 / 壓力 / 稼動率（`startSimulation()` 方法）
- Phase 2 完成後，將 mock 資料換成 `api/` 呼叫，`startSimulation()` 方法刪除

---

## 單元測試規範

### 安裝
```bash
npm install -D vitest @vue/test-utils jsdom @vitest/coverage-v8
```

### vitest.config.ts
```typescript
import { defineConfig } from 'vitest/config'
import vue from '@vitejs/plugin-vue'
import { fileURLToPath } from 'url'

export default defineConfig({
  plugins: [vue()],
  test: {
    environment: 'jsdom',
    globals: true,
    coverage: {
      provider: 'v8',
      reporter: ['text', 'html'],
      include: ['src/stores/**', 'src/composables/**', 'src/api/**'],
      thresholds: { lines: 70 },   // 目標覆蓋率 70%
    },
  },
  resolve: {
    alias: { '@': fileURLToPath(new URL('./src', import.meta.url)) },
  },
})
```

---

### 測試重點一：Pinia Store 邏輯

Store 是最重要的測試目標，業務邏輯集中在此。

```typescript
// __tests__/stores/equipment.test.ts
import { setActivePinia, createPinia } from 'pinia'
import { useEquipmentStore } from '@/stores/equipment'
import { describe, it, expect, beforeEach } from 'vitest'
import type { Equipment } from '@/types/mes'

const mockEquipment = (override: Partial<Equipment> = {}): Equipment => ({
  id: 1, name: 'CVD-01', type: 'CVD', status: 'RUNNING',
  currentLotId: 1, currentLot: 'LOT-001', recipeName: 'CVD-HK-450',
  utilization: 87, temperature: 650, pressure: 120,
  ucl_temp: 680, lcl_temp: 600, ucl_pressure: 150, lcl_pressure: 80,
  isAlarm: false, updatedAt: new Date().toISOString(),
  ...override,
})

describe('useEquipmentStore', () => {
  beforeEach(() => { setActivePinia(createPinia()) })

  // Alarm 置頂排序
  it('Alarm 設備應排在群組最前面', () => {
    const store = useEquipmentStore()
    store.equipments = [
      mockEquipment({ id: 1, name: 'CVD-01', isAlarm: false }),
      mockEquipment({ id: 2, name: 'CVD-02', isAlarm: true }),
      mockEquipment({ id: 3, name: 'CVD-03', isAlarm: false }),
    ]
    const group = store.groupedEquipments.find(g => g.type === 'CVD')
    expect(group?.equipments[0].name).toBe('CVD-02')
  })

  // 群組計算
  it('groupedEquipments 應正確依 type 分組', () => {
    const store = useEquipmentStore()
    store.equipments = [
      mockEquipment({ id: 1, type: 'CVD' }),
      mockEquipment({ id: 2, type: 'Etch' }),
      mockEquipment({ id: 3, type: 'CVD' }),
    ]
    expect(store.groupedEquipments).toHaveLength(2)
    expect(store.groupedEquipments.find(g => g.type === 'CVD')?.equipments).toHaveLength(2)
  })

  // statusCount 計算
  it('statusCount 應正確統計各狀態數量', () => {
    const store = useEquipmentStore()
    store.equipments = [
      mockEquipment({ id: 1, status: 'RUNNING' }),
      mockEquipment({ id: 2, status: 'RUNNING' }),
      mockEquipment({ id: 3, status: 'DOWN' }),
      mockEquipment({ id: 4, status: 'IDLE' }),
      mockEquipment({ id: 5, status: 'PM' }),
    ]
    expect(store.statusCount.running).toBe(2)
    expect(store.statusCount.down).toBe(1)
  })

  // 篩選邏輯
  it('alarmOnly 篩選應只回傳有 Alarm 的群組', () => {
    const store = useEquipmentStore()
    store.equipments = [
      mockEquipment({ id: 1, type: 'CVD', isAlarm: true }),
      mockEquipment({ id: 2, type: 'Etch', isAlarm: false }),
    ]
    store.setFilter({ alarmOnly: true })
    expect(store.filteredGroups).toHaveLength(1)
    expect(store.filteredGroups[0].type).toBe('CVD')
  })

  // updateEquipmentStatus（WebSocket 事件觸發）
  it('updateEquipmentStatus 應正確更新指定設備狀態', () => {
    const store = useEquipmentStore()
    store.equipments = [mockEquipment({ id: 1, status: 'RUNNING' })]
    store.updateEquipmentStatus(1, 'DOWN')
    expect(store.equipments[0].status).toBe('DOWN')
  })

  // 稼動率計算
  it('overallUtilization 應回傳所有設備平均稼動率', () => {
    const store = useEquipmentStore()
    store.equipments = [
      mockEquipment({ id: 1, utilization: 80 }),
      mockEquipment({ id: 2, utilization: 60 }),
    ]
    expect(store.overallUtilization).toBe(70)
  })
})
```

---

### 測試重點二：元件渲染與互動

```typescript
// __tests__/components/equipment/EquipmentRow.test.ts
import { mount } from '@vue/test-utils'
import { createTestingPinia } from '@pinia/testing'
import EquipmentRow from '@/components/equipment/EquipmentRow.vue'
import { describe, it, expect, vi } from 'vitest'

const mockEq = {
  id: 1, name: 'CVD-03', type: 'CVD', status: 'RUNNING',
  currentLot: 'LOT-001', isAlarm: false,
  temperature: 650, pressure: 120,
  ucl_temp: 680, lcl_temp: 600,
  ucl_pressure: 150, lcl_pressure: 80,
  utilization: 87,
}

describe('EquipmentRow', () => {
  // 基本渲染
  it('應正確渲染設備名稱', () => {
    const wrapper = mount(EquipmentRow, {
      props: { equipment: mockEq },
      global: { plugins: [createTestingPinia()] },
    })
    expect(wrapper.text()).toContain('CVD-03')
  })

  // Alarm 樣式
  it('isAlarm 為 true 時應套用紅色邊框 class', () => {
    const wrapper = mount(EquipmentRow, {
      props: { equipment: { ...mockEq, isAlarm: true } },
      global: { plugins: [createTestingPinia()] },
    })
    expect(wrapper.classes()).toContain('alarm-row')
  })

  // 超標參數標示
  it('溫度超過 UCL 時應顯示紅色文字與 ↑ 符號', () => {
    const wrapper = mount(EquipmentRow, {
      props: { equipment: { ...mockEq, temperature: 695, ucl_temp: 680 } },
      global: { plugins: [createTestingPinia()] },
    })
    const tempCell = wrapper.find('[data-testid="temperature"]')
    expect(tempCell.classes()).toContain('over-limit')
    expect(tempCell.text()).toContain('↑')
  })

  // 點擊事件
  it('點擊列時應 emit select 事件並帶入設備資料', async () => {
    const wrapper = mount(EquipmentRow, {
      props: { equipment: mockEq },
      global: { plugins: [createTestingPinia()] },
    })
    await wrapper.trigger('click')
    expect(wrapper.emitted('select')?.[0]).toEqual([mockEq])
  })
})
```

---

### 測試重點三：Composable（useWebSocket）

```typescript
// __tests__/composables/useWebSocket.test.ts
import { describe, it, expect, vi, beforeEach } from 'vitest'
import { createPinia, setActivePinia } from 'pinia'
import { useWebSocket } from '@/composables/useWebSocket'
import { useEquipmentStore } from '@/stores/equipment'
import { useAlarmStore } from '@/stores/alarm'

// Mock WebSocket
class MockWebSocket {
  onmessage: ((e: MessageEvent) => void) | null = null
  send = vi.fn()
  close = vi.fn()
  simulate(data: object) {
    this.onmessage?.({ data: JSON.stringify(data) } as MessageEvent)
  }
}

describe('useWebSocket', () => {
  let mockWs: MockWebSocket

  beforeEach(() => {
    setActivePinia(createPinia())
    mockWs = new MockWebSocket()
    vi.stubGlobal('WebSocket', vi.fn(() => mockWs))
  })

  it('equipment_status_changed 事件應更新 equipmentStore', () => {
    const equipmentStore = useEquipmentStore()
    equipmentStore.equipments = [
      { id: 3, status: 'RUNNING' } as any,
    ]
    useWebSocket()
    mockWs.simulate({ event: 'equipment_status_changed', payload: { equipment_id: 3, status: 'DOWN' } })
    expect(equipmentStore.equipments[0].status).toBe('DOWN')
  })

  it('spc_alarm 事件應新增告警到 alarmStore', () => {
    const alarmStore = useAlarmStore()
    useWebSocket()
    mockWs.simulate({
      event: 'spc_alarm',
      payload: { equipment_id: 3, parameter: 'temperature', value: 692, ucl: 680 },
    })
    expect(alarmStore.alarms).toHaveLength(1)
    expect(alarmStore.alarms[0].parameter).toBe('temperature')
  })
})
```

---

### 測試重點四：API 層 Mock

```typescript
// __tests__/api/equipment.test.ts
import { describe, it, expect, vi, beforeEach } from 'vitest'
import axios from 'axios'
import { fetchEquipments, holdEquipment } from '@/api/equipment'

vi.mock('axios')
const mockedAxios = vi.mocked(axios)

describe('equipment API', () => {
  beforeEach(() => { vi.clearAllMocks() })

  it('fetchEquipments 應呼叫正確端點並回傳資料', async () => {
    const mockData = [{ id: 1, name: 'CVD-01', status: 'RUNNING' }]
    mockedAxios.get = vi.fn().mockResolvedValue({ data: { data: mockData } })

    const result = await fetchEquipments()

    expect(mockedAxios.get).toHaveBeenCalledWith('/api/v1/equipment')
    expect(result).toEqual(mockData)
  })

  it('holdEquipment 應呼叫 POST /api/v1/equipment/:id/hold', async () => {
    mockedAxios.post = vi.fn().mockResolvedValue({ data: { message: 'ok' } })

    await holdEquipment(3)

    expect(mockedAxios.post).toHaveBeenCalledWith('/api/v1/equipment/3/hold')
  })

  it('fetchEquipments API 失敗時應拋出錯誤', async () => {
    mockedAxios.get = vi.fn().mockRejectedValue(new Error('Network Error'))
    await expect(fetchEquipments()).rejects.toThrow('Network Error')
  })
})
```

---

### 測試優先級原則（給 Claude CLI 判斷用）

當不確定哪裡需要測試時，依以下優先序決定：

**P0 必測（有業務邏輯的地方）**
- Pinia Store 的 getter 和 action（排序、篩選、計算邏輯）
- Composable 的事件分發邏輯（useWebSocket）
- 有條件判斷的工具函式（超標判定、狀態轉換）

**P1 選測（有互動行為的元件）**
- 使用者操作會觸發 emit 或 store 變化的元件
- 有條件渲染（v-if）或動態 class 的關鍵元件
- API 層的端點正確性與錯誤處理

**P2 不必測（純樣板）**
- 只負責排版的純展示元件（沒有邏輯的 template）
- router/index.ts 路由設定
- 型別定義檔 types/mes.ts
- style.css / 主題設定

**原則：測行為，不測實作細節。** 測試應該描述「這個功能做了什麼」，而不是「內部怎麼實作的」。如果重構內部實作但行為不變，測試不應該壞掉。

---

### 測試命名與組織規則

- 測試檔命名：`對應原始檔名.test.ts`，放在 `__tests__/` 同層結構下
- `describe` 用**元件或 store 名稱**
- `it` 描述句用**中文**，說明預期行為（誰 / 做什麼 / 結果）
- 每個測試只測一件事，不在同一個 `it` 裡塞多個 `expect`

```typescript
// 好的命名
it('Alarm 設備應排在群組最前面')
it('溫度超過 UCL 時應顯示紅色文字與 ↑ 符號')
it('點擊列時應 emit select 事件並帶入設備資料')

// 避免的命名
it('test equipment')
it('should work')
```

---

## 啟動指令

```bash
# 安裝依賴
cd frontend
npm install

# 開發模式
npm run dev           # http://localhost:5173

# 型別檢查
npm run type-check

# 單元測試
npm run test          # 單次執行
npm run test:watch    # 監聽模式（開發時使用）
npm run test:coverage # 產出覆蓋率報告

# 建置
npm run build
```

### package.json scripts 對應設定
```json
{
  "scripts": {
    "test": "vitest run",
    "test:watch": "vitest",
    "test:coverage": "vitest run --coverage"
  }
}
```