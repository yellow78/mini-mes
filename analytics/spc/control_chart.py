"""
SPC 管制圖分析模組 — Nelson Rules 實作
Nelson Rules 是半導體業常用的製程管制判斷準則
"""
from __future__ import annotations

import numpy as np
from typing import TypedDict


class AnalysisResult(TypedDict):
    mean: float
    std: float
    ucl: float
    lcl: float
    is_alarm: bool
    violations: list[str]
    alarm_indices: list[int]
    severity: str


def analyze(values: list[float], ucl: float, lcl: float) -> AnalysisResult:
    """
    依 Nelson Rules 判斷 SPC 是否異常

    Args:
        values: 量測值序列（時序由舊到新）
        ucl:    Upper Control Limit（管制上限，通常為均值 +3σ）
        lcl:    Lower Control Limit（管制下限，通常為均值 -3σ）

    Returns:
        AnalysisResult: 分析結果，含違反規則清單與嚴重程度
    """
    if not values:
        return AnalysisResult(
            mean=0.0, std=0.0, ucl=ucl, lcl=lcl,
            is_alarm=False, violations=[], alarm_indices=[], severity="",
        )

    arr = np.array(values, dtype=float)
    mean = float(arr.mean())
    std = float(arr.std(ddof=1)) if len(arr) > 1 else 0.0

    violations: list[str] = []
    alarm_set: set[int] = set()

    # --- Nelson Rule 1 ---
    # 任一點超出 UCL 或 LCL（3σ 管制界限）
    for i, v in enumerate(values):
        if v > ucl or v < lcl:
            direction = "UCL" if v > ucl else "LCL"
            violations.append(f"Rule1: point[{i}]={v:.2f} 超出{direction}={ucl if v > ucl else lcl:.2f}")
            alarm_set.add(i)

    # --- Nelson Rule 2 ---
    # 連續 9 點在均值同一側（製程偏移）
    if len(values) >= 9:
        for i in range(len(values) - 8):
            seg = values[i:i + 9]
            if all(v > mean for v in seg) or all(v < mean for v in seg):
                side = "上方" if seg[0] > mean else "下方"
                violations.append(f"Rule2: 連續9點在均值{side}（製程偏移）")
                alarm_set.update(range(i, i + 9))
                break

    # --- Nelson Rule 3 ---
    # 連續 6 點單調遞增或遞減（製程趨勢）
    if len(values) >= 6:
        for i in range(len(values) - 5):
            seg = values[i:i + 6]
            increasing = all(seg[j] < seg[j + 1] for j in range(5))
            decreasing = all(seg[j] > seg[j + 1] for j in range(5))
            if increasing or decreasing:
                direction = "上升" if increasing else "下降"
                violations.append(f"Rule3: 連續6點{direction}趨勢（漂移）")
                alarm_set.update(range(i, i + 6))
                break

    is_alarm = len(violations) > 0

    # Rule 1 為 CRITICAL（直接超標），Rule 2/3 為 WARNING（趨勢異常）
    severity = ""
    if is_alarm:
        has_rule1 = any(v.startswith("Rule1") for v in violations)
        severity = "CRITICAL" if has_rule1 else "WARNING"

    return AnalysisResult(
        mean=mean,
        std=std,
        ucl=ucl,
        lcl=lcl,
        is_alarm=is_alarm,
        violations=violations,
        alarm_indices=sorted(alarm_set),
        severity=severity,
    )
