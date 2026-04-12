"""
Mini-MES 分析引擎 — FastAPI 服務
提供 SPC（統計製程管制）分析端點，供 Go 後端呼叫
"""
from fastapi import FastAPI
from pydantic import BaseModel, Field

from spc.control_chart import analyze

app = FastAPI(title="Mini-MES Analytics", version="1.0.0")


# ── 請求 / 回應 Schema ──────────────────────────────────────

class SpcAnalyzeRequest(BaseModel):
    equipment_id: int
    parameter: str                      # "temperature" | "pressure"
    values: list[float] = Field(..., min_length=1)
    ucl: float
    lcl: float


class SpcAnalyzeResponse(BaseModel):
    equipment_id: int
    parameter: str
    mean: float
    std: float
    ucl: float
    lcl: float
    is_alarm: bool
    violations: list[str]
    alarm_indices: list[int]
    severity: str                       # "CRITICAL" | "WARNING" | ""


# ── 端點 ────────────────────────────────────────────────────

@app.get("/health")
def health():
    """健康檢查（供 Docker healthcheck 使用）"""
    return {"status": "ok"}


@app.post("/spc/analyze", response_model=SpcAnalyzeResponse)
def spc_analyze(req: SpcAnalyzeRequest) -> SpcAnalyzeResponse:
    """
    對給定的量測序列執行 SPC 分析（Nelson Rules 1/2/3）
    Go 後端在模擬器每次產生新量測值時呼叫此端點
    """
    result = analyze(req.values, req.ucl, req.lcl)
    return SpcAnalyzeResponse(
        equipment_id=req.equipment_id,
        parameter=req.parameter,
        **result,
    )
