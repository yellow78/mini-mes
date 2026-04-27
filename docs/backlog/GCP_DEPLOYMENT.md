# Mini-MES GCP 部署計畫

> 建立日期：2026-04-27
> 狀態：待執行

---

## Context

Mini-MES 目前以 Docker Compose 在本機運行，需要部署至 GCP 以支援公開展示與遠端 demo。
目標：以最低成本、最少程式碼修改，將四個服務部署至 GCP，確保 WebSocket 即時更新正常運作。

---

## 架構選型：Cloud Run + Neon（PostgreSQL）

**選 Cloud Run 的理由**（vs GKE / Compute Engine）：

| 比較 | Cloud Run | GKE | Compute Engine |
|------|-----------|-----|----------------|
| 費用（閒置） | 接近 $0 | $50–150/月 | $20–40/月 |
| 設定時間 | 30–60 分鐘 | 2–4 小時 | 1–2 小時 |
| WebSocket | 支援（需調參） | 原生 | 原生 |
| HTTPS | 自動 | 需設 Ingress | 需自設 |

**選 Neon 取代 Cloud SQL 的理由**：
- 免費層：0.5 GB 儲存 + 每月 191.9 compute-hours，demo 完全夠用
- 計算資源自動縮放至零，閒置不收費
- 連線方式與標準 PostgreSQL 完全相同，後端程式碼零修改
- 無需 Cloud SQL Proxy 或 VPC 設定，只需換環境變數

**目標架構**：

```
Internet (HTTPS)
    ↓
Cloud Run: frontend (nginx)
    ├── /api/*  → Cloud Run: backend (Go + Gin)
    └── /ws     → Cloud Run: backend (WebSocket)
                      ├── ANALYTICS_URL → Cloud Run: analytics (Python, internal ingress)
                      └── DATABASE_URL  → Neon (PostgreSQL, 外部 TCP)
```

---

## 必要的程式碼修改（共 3 處）

### 1. `frontend/nginx.conf`（必改）

**問題**：`set $backend http://backend:8080;` 使用 Docker 內部 DNS，Cloud Run 不存在此 DNS。

**修改**：改用 nginx 官方 template 機制，啟動時從環境變數 `BACKEND_URL` 注入。

```nginx
server {
    listen 80;

    location / {
        root   /usr/share/nginx/html;
        index  index.html;
        try_files $uri $uri/ /index.html;
    }

    location /api/ {
        proxy_pass         ${BACKEND_URL};
        proxy_set_header   Host $http_host;
        proxy_set_header   X-Real-IP $remote_addr;
    }

    location /ws {
        proxy_pass         ${BACKEND_URL};
        proxy_http_version 1.1;
        proxy_set_header   Upgrade $http_upgrade;
        proxy_set_header   Connection "upgrade";
        proxy_set_header   Host $http_host;
        proxy_read_timeout 3600s;
    }
}
```

> 注意：`$http_host`、`$http_upgrade` 等 nginx 內建變數需保留，
> 只讓 `envsubst` 替換 `$BACKEND_URL`。
> 做法：在 Dockerfile 加入 `ENV NGINX_ENVSUBST_FILTER=BACKEND_URL`。

### 2. `frontend/Dockerfile`（必改）

```dockerfile
FROM node:20-alpine AS builder
WORKDIR /app
COPY package.json package-lock.json ./
RUN npm ci
COPY . .
RUN npm run build

FROM nginx:alpine
RUN apk add --no-cache gettext
COPY --from=builder /app/dist /usr/share/nginx/html
# 存為 .template，nginx:alpine entrypoint 自動執行 envsubst
COPY nginx.conf /etc/nginx/templates/default.conf.template
ENV NGINX_ENVSUBST_FILTER=BACKEND_URL
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
```

### 3. `backend/Dockerfile`（確認版本號）

第 1 行 `golang:1.25-alpine` 若映像不存在會 build 失敗，改為已發布的版本：

```dockerfile
FROM golang:1.23-alpine AS builder
```

---

## 不需修改的檔案（已相容 Cloud Run + Neon）

| 檔案 | 說明 |
|------|------|
| `frontend/src/composables/useWebSocket.ts` | 使用 `window.location.host` 動態產生 WS URL |
| `frontend/src/api/equipment.ts` | 使用相對路徑 `/api/v1/...`，由 nginx proxy 轉發 |
| `backend/cmd/server/main.go` | `lib/pq` 標準 TCP 連線，Neon 連線字串格式相同 |
| `backend/cmd/server/main.go` | `ANALYTICS_URL` 從環境變數讀取，不需修改 |

---

## 部署步驟

### 前置一：Neon 資料庫設定

1. 前往 [neon.tech](https://neon.tech) 註冊免費帳號
2. 建立新 Project（選擇離台灣最近的 region，如 `aws-ap-southeast-1`）
3. 建立 Database，名稱設為 `mes_production`
4. 複製連線字串，格式如下：
   ```
   postgresql://mes_prod:PASSWORD@ep-xxxx.ap-southeast-1.aws.neon.tech/mes_production?sslmode=require
   ```
5. 使用連線字串執行 migration：
   ```bash
   psql "postgresql://mes_prod:PASSWORD@ep-xxxx.ap-southeast-1.aws.neon.tech/mes_production?sslmode=require" \
     -f migrations/001_init.sql
   psql "postgresql://mes_prod:PASSWORD@ep-xxxx.ap-southeast-1.aws.neon.tech/mes_production?sslmode=require" \
     -f migrations/seed.sql
   ```

### 前置二：GCP 環境設定

```bash
export PROJECT_ID="mini-mes-demo"
export REGION="asia-east1"
export REGISTRY="${REGION}-docker.pkg.dev/${PROJECT_ID}/mini-mes"

# 只需要 Cloud Run + Artifact Registry，不需要 sqladmin
gcloud services enable \
  run.googleapis.com \
  artifactregistry.googleapis.com \
  secretmanager.googleapis.com

gcloud artifacts repositories create mini-mes \
  --repository-format=docker --location=${REGION}

gcloud auth configure-docker ${REGION}-docker.pkg.dev

gcloud iam service-accounts create mini-mes-runner \
  --display-name="Mini MES Cloud Run Runner"
```

### Step 1：Secret Manager 儲存 Neon 連線字串

```bash
# 將 Neon 連線字串存入 Secret Manager
echo -n "postgresql://mes_prod:PASSWORD@ep-xxxx.ap-southeast-1.aws.neon.tech/mes_production?sslmode=require" \
  | gcloud secrets create neon-database-url --data-file=-

gcloud secrets add-iam-policy-binding neon-database-url \
  --member="serviceAccount:mini-mes-runner@${PROJECT_ID}.iam.gserviceaccount.com" \
  --role="roles/secretmanager.secretAccessor"
```

### Step 2：Build & Push Images

```bash
docker build -t ${REGISTRY}/analytics:latest ./analytics
docker build -t ${REGISTRY}/backend:latest   ./backend
docker build -t ${REGISTRY}/frontend:latest  ./frontend

docker push ${REGISTRY}/analytics:latest
docker push ${REGISTRY}/backend:latest
docker push ${REGISTRY}/frontend:latest
```

### Step 3：部署至 Cloud Run（依序）

```bash
# 3a. Analytics（僅內部流量）
gcloud run deploy analytics \
  --image=${REGISTRY}/analytics:latest \
  --region=${REGION} \
  --service-account=mini-mes-runner@${PROJECT_ID}.iam.gserviceaccount.com \
  --port=8001 --memory=512Mi \
  --ingress=internal \
  --allow-unauthenticated

ANALYTICS_URL=$(gcloud run services describe analytics \
  --region=${REGION} --format="value(status.url)")

# 3b. Backend（連接 Neon，不需要 --add-cloudsql-instances）
# DATABASE_URL 格式讓 backend 直接用整串連線，
# 或維持原本拆開的 DB_HOST/DB_PORT/DB_NAME/DB_USER/DB_PASSWORD 環境變數
NEON_URL=$(gcloud secrets versions access latest --secret=neon-database-url)
# 拆解連線字串為個別環境變數（維持後端原有讀取方式）
export DB_HOST=$(echo $NEON_URL | sed 's/.*@//' | cut -d'/' -f1 | cut -d':' -f1)
export DB_PORT=5432
export DB_NAME=mes_production
export DB_USER=mes_prod
export DB_PASSWORD=$(echo $NEON_URL | sed 's/.*://' | cut -d'@' -f1)

gcloud run deploy backend \
  --image=${REGISTRY}/backend:latest \
  --region=${REGION} \
  --service-account=mini-mes-runner@${PROJECT_ID}.iam.gserviceaccount.com \
  --port=8080 --memory=512Mi \
  --min-instances=0 \
  --timeout=3600 \
  --allow-unauthenticated \
  --set-env-vars="DB_HOST=${DB_HOST},DB_PORT=5432,DB_NAME=mes_production,DB_USER=mes_prod,API_PORT=8080,ANALYTICS_URL=${ANALYTICS_URL}" \
  --set-secrets="DB_PASSWORD=neon-database-url:latest"

# 備註：若後端支援 DATABASE_URL 單一環境變數，可簡化為：
# --set-secrets="DATABASE_URL=neon-database-url:latest"

BACKEND_URL=$(gcloud run services describe backend \
  --region=${REGION} --format="value(status.url)")

# 3c. Frontend（注入 BACKEND_URL）
gcloud run deploy frontend \
  --image=${REGISTRY}/frontend:latest \
  --region=${REGION} \
  --service-account=mini-mes-runner@${PROJECT_ID}.iam.gserviceaccount.com \
  --port=80 --memory=256Mi \
  --allow-unauthenticated \
  --set-env-vars="BACKEND_URL=${BACKEND_URL}"

FRONTEND_URL=$(gcloud run services describe frontend \
  --region=${REGION} --format="value(status.url)")
echo "Demo URL: ${FRONTEND_URL}"
```

> **補充說明**：Backend 目前使用 `DB_HOST/DB_PORT/DB_NAME/DB_USER/DB_PASSWORD` 分開讀取。
> Neon 連線需加上 `?sslmode=require`，需確認 `connectDB()` 的 DSN 格式有帶 `sslmode=require`。
> 若沒有，可在 DB_HOST 加上 `sslmode` 參數或修改 `main.go` 的 DSN 字串。

---

## WebSocket 相容性

| 問題 | 解法 |
|------|------|
| Cloud Run 預設 timeout 300s | `--timeout=3600` |
| 冷啟動斷線（min-instances=0） | 前端 `useWebSocket.ts` 已有 5 秒自動重連，可接受 |
| 前端斷線重連 | `useWebSocket.ts` 已有 5 秒重連邏輯，不需修改 |

---

## 費用估算（Demo 用途）

| 服務 | 規格 | 月費 |
|------|------|------|
| Cloud Run: backend | 512Mi, min-instances=0，按請求計費 | ~$0–2 |
| Cloud Run: analytics / frontend | 512Mi / 256Mi，按請求計費 | ~$0–1 |
| Neon PostgreSQL | 免費層（0.5 GB） | **$0** |
| Artifact Registry | 少量 image 儲存 | ~$0.5 |
| **合計** | | **~$0.5–3/月** |

> 對比原方案（Cloud SQL）：~$12–20/月 → **降低 85%**

---

## 驗證步驟

```bash
# 1. 確認所有服務狀態
gcloud run services list --region=${REGION}

# 2. 確認 API 可達（冷啟動第一次約等 2–5 秒）
curl ${FRONTEND_URL}/api/v1/equipment | python3 -m json.tool

# 3. 開啟瀏覽器 → ${FRONTEND_URL}
#    DevTools → Network → WS → 確認 /ws 連線成功
#    Dashboard 應顯示設備列表，10 秒後出現 WebSocket 狀態更新

# 4. 確認 analytics 不對外暴露（應拒絕連線）
curl ${ANALYTICS_URL}/health

# 5. 確認 backend log 正常（確認 Neon 連線成功）
gcloud run services logs read backend --region=${REGION} --limit=50
```
