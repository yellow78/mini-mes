# Mini-MES Backend — Claude CLI 後端開發手冊

> 全域資訊（專案定位、Phase 順序、產業術語、WebSocket 事件格式）請參考根目錄 `../CLAUDE.md`

## Claude 互動規範
- 回答語言：繁體中文
- 程式碼註解：繁體中文
- 變數、函式、檔案名稱：英文

---

## 技術棧
- Go 1.22+
- Gin（HTTP 框架）
- sqlx（DB 存取）
- lib/pq（PostgreSQL driver）
- gorilla/websocket（WebSocket）
- go-redis/v9（Redis，Phase 3 加入）
- testify（測試斷言）
- testcontainers-go（repository 層整合測試）

---

## 目錄結構

```
backend/
├── cmd/
│   ├── server/
│   │   └── main.go              ← 程式進入點
│   └── seed/
│       └── main.go              ← 資料初始化工具
├── internal/
│   ├── handler/                 ← HTTP 路由層（解析 request、呼叫 service、回傳 JSON）
│   │   ├── equipment.go
│   │   ├── lot.go
│   │   └── alarm.go
│   ├── service/                 ← 業務邏輯層（狀態驗證、群組計算、派工邏輯）
│   │   ├── equipment.go
│   │   ├── equipment_test.go
│   │   ├── lot.go
│   │   ├── lot_test.go
│   │   └── dispatch.go
│   ├── repository/              ← DB 存取層（純 SQL，不含業務邏輯）
│   │   ├── equipment.go
│   │   ├── equipment_test.go    ← testcontainers 整合測試
│   │   ├── lot.go
│   │   └── interfaces.go        ← repository interface 定義（供 mock 使用）
│   └── model/                   ← 資料結構定義
│       ├── equipment.go
│       ├── lot.go
│       ├── wafer.go
│       ├── recipe.go
│       └── spc_record.go
├── pkg/
│   └── websocket/
│       └── hub.go               ← WebSocket 廣播中心
├── migrations/
│   ├── 001_init.sql
│   └── seed.sql
├── .env.example
├── go.mod
└── go.sum
```

---

## 分層架構說明

```
HTTP Request
    ↓
Handler 層        解析 request params / body，呼叫 service，回傳 JSON
    ↓
Service 層        業務邏輯（狀態驗證、計算、規則）
    ↓
Repository 層     DB 存取（SQL 查詢），透過 interface 解耦
    ↓
PostgreSQL
```

**原則：**
- Handler 不含業務邏輯，只做 HTTP 解析與回應
- Service 不直接操作 DB，只透過 repository interface
- Repository 不含業務邏輯，只做 CRUD
- 測試時 Service 注入 mock repository，不需要真實 DB

---

## 各服務 Port

| 服務       | Port  |
|------------|-------|
| Go API     | 8080  |
| PostgreSQL | 5432  |
| Redis      | 6379  |

---

## API 端點規範

所有路由前綴：`/api/v1/`

### Equipment
```
GET    /api/v1/equipment               設備列表（含狀態、當前 Lot）
GET    /api/v1/equipment/:id           單台設備詳細
PUT    /api/v1/equipment/:id/status    更新設備狀態
POST   /api/v1/equipment/:id/hold      Hold 設備（狀態改為 DOWN，Lot 改為 ON_HOLD）
```

### Lot
```
GET    /api/v1/lots                    Lot 列表（含 WIP 狀態）
POST   /api/v1/lots                    建立新 Lot
GET    /api/v1/lots/:id                單一 Lot 詳細
POST   /api/v1/lots/:id/dispatch       觸發派工（指派到可用設備）
```

### Alarm
```
GET    /api/v1/alarms                  告警列表（預設只回未確認）
PUT    /api/v1/alarms/:id/acknowledge  確認告警
```

### SPC
```
GET    /api/v1/spc/:equipment_id       SPC 歷史資料（預設最近 100 筆）
```

### WebSocket
```
WS     /ws                             前端連線，接收即時事件
```

### 統一回應格式
```json
// 成功（列表）
{ "data": [...], "total": 100 }

// 成功（單筆）
{ "data": {...} }

// 錯誤
{ "error": "錯誤說明" }
```

---

## Model 定義

### equipment.go
```go
type EquipmentStatus string

const (
    StatusIdle    EquipmentStatus = "IDLE"
    StatusRunning EquipmentStatus = "RUNNING"
    StatusDown    EquipmentStatus = "DOWN"
    StatusPM      EquipmentStatus = "PM"
)

// 合法的狀態轉換規則
var ValidTransitions = map[EquipmentStatus][]EquipmentStatus{
    StatusIdle:    {StatusRunning, StatusPM, StatusDown},
    StatusRunning: {StatusIdle, StatusDown},
    StatusDown:    {StatusPM},
    StatusPM:      {StatusIdle},
}

type Equipment struct {
    ID           int             `db:"id"             json:"id"`
    Name         string          `db:"name"           json:"name"`
    Type         string          `db:"type"           json:"type"`
    Status       EquipmentStatus `db:"status"         json:"status"`
    CurrentLotID *int            `db:"current_lot_id" json:"current_lot_id"`
    CurrentLot   *string         `db:"-"              json:"current_lot"`   // JOIN 欄位
    RecipeName   *string         `db:"-"              json:"recipe_name"`
    Utilization  float64         `db:"utilization"    json:"utilization"`
    Temperature  float64         `db:"temperature"    json:"temperature"`
    Pressure     float64         `db:"pressure"       json:"pressure"`
    UCLTemp      float64         `db:"ucl_temp"       json:"ucl_temp"`
    LCLTemp      float64         `db:"lcl_temp"       json:"lcl_temp"`
    UCLPressure  float64         `db:"ucl_pressure"   json:"ucl_pressure"`
    LCLPressure  float64         `db:"lcl_pressure"   json:"lcl_pressure"`
    IsAlarm      bool            `db:"is_alarm"       json:"is_alarm"`
    UpdatedAt    time.Time       `db:"updated_at"     json:"updated_at"`
}
```

### lot.go
```go
type LotStatus string

const (
    LotQueued    LotStatus = "QUEUED"
    LotRunning   LotStatus = "RUNNING"
    LotCompleted LotStatus = "COMPLETED"
    LotOnHold    LotStatus = "ON_HOLD"
)

type Lot struct {
    ID          int       `db:"id"           json:"id"`
    LotNumber   string    `db:"lot_number"   json:"lot_number"`
    Product     string    `db:"product"      json:"product"`
    RecipeID    int       `db:"recipe_id"    json:"recipe_id"`
    RecipeName  *string   `db:"-"            json:"recipe_name"`
    Priority    int       `db:"priority"     json:"priority"`
    Status      LotStatus `db:"status"       json:"status"`
    WaferCount  int       `db:"wafer_count"  json:"wafer_count"`
    CreatedAt   time.Time `db:"created_at"   json:"created_at"`
}
```

---

## Repository Interface（interfaces.go）

```go
// 所有 repository 必須實作 interface，讓 service 層可以 mock

type EquipmentRepository interface {
    FindAll(ctx context.Context) ([]model.Equipment, error)
    FindByID(ctx context.Context, id int) (*model.Equipment, error)
    UpdateStatus(ctx context.Context, id int, status model.EquipmentStatus) error
    FindIdleByType(ctx context.Context, equipType string) (*model.Equipment, error)
}

type LotRepository interface {
    FindAll(ctx context.Context) ([]model.Lot, error)
    FindByID(ctx context.Context, id int) (*model.Lot, error)
    Create(ctx context.Context, lot *model.Lot) error
    UpdateStatus(ctx context.Context, id int, status model.LotStatus) error
}
```

---

## 單元測試規範

### 測試框架
- **標準 `testing` 套件** + **testify**（斷言）
- **service 層**：注入 mock repository，純邏輯測試，不需要 DB
- **repository 層**：使用 testcontainers-go 起真實 PostgreSQL

### 安裝
```bash
go get github.com/stretchr/testify
go get github.com/testcontainers/testcontainers-go
go get github.com/testcontainers/testcontainers-go/modules/postgres
```

### Service 層測試範例（equipment_test.go）

```go
// 用 mock repository 測試業務邏輯，不依賴 DB

// MockEquipmentRepository 實作 EquipmentRepository interface
type MockEquipmentRepository struct {
    equipments []model.Equipment
    updateErr  error
}

func (m *MockEquipmentRepository) FindAll(ctx context.Context) ([]model.Equipment, error) {
    return m.equipments, nil
}
func (m *MockEquipmentRepository) UpdateStatus(ctx context.Context, id int, status model.EquipmentStatus) error {
    return m.updateErr
}
// ... 其他 interface 方法

func TestEquipmentService_UpdateStatus(t *testing.T) {

    t.Run("合法狀態轉換應成功", func(t *testing.T) {
        mock := &MockEquipmentRepository{
            equipments: []model.Equipment{
                {ID: 1, Status: model.StatusIdle},
            },
        }
        svc := NewEquipmentService(mock)

        err := svc.UpdateStatus(context.Background(), 1, model.StatusRunning)
        assert.NoError(t, err)
    })

    t.Run("非法狀態轉換應回傳錯誤", func(t *testing.T) {
        // DOWN → RUNNING 不合法，必須先 PM
        mock := &MockEquipmentRepository{
            equipments: []model.Equipment{
                {ID: 1, Status: model.StatusDown},
            },
        }
        svc := NewEquipmentService(mock)

        err := svc.UpdateStatus(context.Background(), 1, model.StatusRunning)
        assert.ErrorContains(t, err, "invalid status transition")
    })
}

func TestEquipmentService_GetGroups(t *testing.T) {

    t.Run("應依類型正確分群", func(t *testing.T) {
        mock := &MockEquipmentRepository{
            equipments: []model.Equipment{
                {ID: 1, Type: "CVD",  Status: model.StatusRunning},
                {ID: 2, Type: "CVD",  Status: model.StatusIdle},
                {ID: 3, Type: "Etch", Status: model.StatusRunning, IsAlarm: true},
            },
        }
        svc := NewEquipmentService(mock)

        groups, err := svc.GetGroups(context.Background())
        assert.NoError(t, err)
        assert.Len(t, groups, 2)

        cvd := findGroup(groups, "CVD")
        assert.Equal(t, 2, len(cvd.Equipments))
    })

    t.Run("Alarm 設備應排在群組最頂部", func(t *testing.T) {
        mock := &MockEquipmentRepository{
            equipments: []model.Equipment{
                {ID: 1, Type: "CVD", Name: "CVD-01", IsAlarm: false},
                {ID: 2, Type: "CVD", Name: "CVD-02", IsAlarm: true},
            },
        }
        svc := NewEquipmentService(mock)

        groups, _ := svc.GetGroups(context.Background())
        cvd := findGroup(groups, "CVD")
        assert.Equal(t, "CVD-02", cvd.Equipments[0].Name)
    })
}
```

### Repository 層測試範例（equipment_test.go）

```go
// 使用 testcontainers 起真實 PostgreSQL，測試 SQL 正確性

func setupTestDB(t *testing.T) *sqlx.DB {
    ctx := context.Background()
    container, err := postgres.Run(ctx,
        "postgres:16-alpine",
        postgres.WithDatabase("mes_test"),
        postgres.WithUsername("test"),
        postgres.WithPassword("test"),
        testcontainers.WithWaitStrategy(
            wait.ForLog("database system is ready").
                WithStartupTimeout(30*time.Second),
        ),
    )
    require.NoError(t, err)
    t.Cleanup(func() { container.Terminate(ctx) })

    connStr, _ := container.ConnectionString(ctx, "sslmode=disable")
    db, err := sqlx.Connect("postgres", connStr)
    require.NoError(t, err)

    // 執行 migration
    migration, _ := os.ReadFile("../../migrations/001_init.sql")
    db.MustExec(string(migration))
    return db
}

func TestEquipmentRepository_FindAll(t *testing.T) {
    db := setupTestDB(t)
    repo := NewEquipmentRepository(db)

    // 插入測試資料
    db.MustExec(`INSERT INTO equipment (name, type, status) VALUES ('CVD-01', 'CVD', 'IDLE')`)

    equipments, err := repo.FindAll(context.Background())
    assert.NoError(t, err)
    assert.Len(t, equipments, 1)
    assert.Equal(t, "CVD-01", equipments[0].Name)
}
```

### 測試執行指令
```bash
# 所有測試
go test ./...

# 只跑 service 層（快，不需要 Docker）
go test ./internal/service/...

# 只跑 repository 層（需要 Docker）
go test ./internal/repository/...

# 含覆蓋率報告
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

---

## Docker Compose 設定（Phase 2 版本）

```yaml
# docker-compose.yml
services:
  postgres:
    image: postgres:16-alpine
    environment:
      POSTGRES_DB: mes_dev
      POSTGRES_USER: mes
      POSTGRES_PASSWORD: mes_password
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./migrations/001_init.sql:/docker-entrypoint-initdb.d/001_init.sql
      - ./migrations/seed.sql:/docker-entrypoint-initdb.d/002_seed.sql
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U mes -d mes_dev"]
      interval: 5s
      timeout: 5s
      retries: 5

  backend:
    build: ./backend
    ports:
      - "8080:8080"
    environment:
      DB_HOST: postgres
      DB_PORT: 5432
      DB_NAME: mes_dev
      DB_USER: mes
      DB_PASSWORD: mes_password
    depends_on:
      postgres:
        condition: service_healthy

volumes:
  postgres_data:
```

---

## 環境變數（.env.example）

```env
DB_HOST=localhost
DB_PORT=5432
DB_NAME=mes_dev
DB_USER=mes
DB_PASSWORD=mes_password
API_PORT=8080
```

---

## Seed 資料規格（migrations/seed.sql）

插入資料量：
- Recipe：8 筆（每種設備類型 2 個 Recipe）
- Equipment：100 台（CVD×24、Etch×20、CMP×16、Diffusion×40）
- Lot：10 筆（3 筆 RUNNING、4 筆 QUEUED、3 筆 COMPLETED）

設備狀態分佈（模擬真實產線）：
- 約 65% RUNNING
- 約 20% IDLE
- 約 10% DOWN
- 約 5% PM

---

## 啟動指令

```bash
# 啟動 DB + 後端
docker compose up -d

# 查看後端 log
docker compose logs -f backend

# 只啟動 DB（本機跑 Go 開發用）
docker compose up -d postgres

# 本機跑 Go
cd backend
go run cmd/server/main.go

# 停止並清除資料
docker compose down -v
```

---

## 給 Claude CLI 的指令範例

```bash
# 進入 Phase 2 開發時
> 請閱讀 CLAUDE.md 和 backend/CLAUDE.md，
  現在要開始 Phase 2，先幫我建立 migrations/001_init.sql

> 依照 backend/CLAUDE.md 的 interface 規範，
  建立 internal/repository/interfaces.go

> 幫我實作 internal/service/equipment.go 的 UpdateStatus 方法，
  包含狀態轉換驗證邏輯，並同時建立對應的 equipment_test.go
```
