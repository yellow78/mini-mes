"""SPC 管制圖 analyze() 函式單元測試 — Nelson Rules 1 / 2 / 3"""
import sys
import os

sys.path.insert(0, os.path.dirname(os.path.dirname(os.path.abspath(__file__))))

import pytest
from spc.control_chart import analyze

# ─── 共用常數 ────────────────────────────────────────────────
UCL = 100.0
LCL = 0.0
WIDE_UCL = 1000.0
WIDE_LCL = -1000.0


# ─── 空輸入 / 基礎統計 ───────────────────────────────────────

class TestBasic:
    def test_empty_values_no_alarm(self):
        r = analyze([], UCL, LCL)
        assert r['is_alarm'] is False
        assert r['violations'] == []
        assert r['alarm_indices'] == []
        assert r['severity'] == ''

    def test_empty_values_mean_std_zero(self):
        r = analyze([], UCL, LCL)
        assert r['mean'] == 0.0
        assert r['std'] == 0.0

    def test_empty_values_preserves_ucl_lcl(self):
        r = analyze([], UCL, LCL)
        assert r['ucl'] == UCL
        assert r['lcl'] == LCL

    def test_single_value_within_limits_no_alarm(self):
        r = analyze([50.0], UCL, LCL)
        assert r['is_alarm'] is False
        assert r['std'] == 0.0

    def test_mean_correct(self):
        r = analyze([1.0, 2.0, 3.0], WIDE_UCL, WIDE_LCL)
        assert r['mean'] == pytest.approx(2.0)

    def test_std_uses_ddof1(self):
        # ddof=1: sqrt(((1-2)^2 + (2-2)^2 + (3-2)^2) / 2) = 1.0
        r = analyze([1.0, 2.0, 3.0], WIDE_UCL, WIDE_LCL)
        assert r['std'] == pytest.approx(1.0)

    def test_no_alarm_severity_is_empty_string(self):
        r = analyze([50.0, 51.0, 49.0], UCL, LCL)
        assert r['severity'] == ''


# ─── Nelson Rule 1：最新值超出管制限 ────────────────────────

class TestRule1:
    def test_latest_above_ucl_is_critical(self):
        r = analyze([50.0, 50.0, 50.0, 110.0], UCL, LCL)
        assert r['is_alarm'] is True
        assert r['severity'] == 'CRITICAL'
        assert any('Rule1' in v for v in r['violations'])

    def test_latest_below_lcl_is_critical(self):
        r = analyze([50.0, 50.0, 50.0, -5.0], UCL, LCL)
        assert r['is_alarm'] is True
        assert r['severity'] == 'CRITICAL'
        assert any('Rule1' in v for v in r['violations'])

    def test_alarm_index_is_last_element(self):
        values = [50.0, 50.0, 50.0, 110.0]
        r = analyze(values, UCL, LCL)
        assert len(values) - 1 in r['alarm_indices']

    def test_exactly_ucl_does_not_trigger(self):
        # 邊界值：嚴格 > 才觸發
        r = analyze([50.0, 50.0, 100.0], UCL, LCL)
        assert r['is_alarm'] is False

    def test_exactly_lcl_does_not_trigger(self):
        r = analyze([50.0, 50.0, 0.0], UCL, LCL)
        assert r['is_alarm'] is False

    def test_violation_message_contains_ucl_keyword(self):
        r = analyze([50.0, 110.0], UCL, LCL)
        assert 'UCL' in r['violations'][0]

    def test_violation_message_contains_lcl_keyword(self):
        r = analyze([50.0, -5.0], UCL, LCL)
        assert 'LCL' in r['violations'][0]


# ─── Nelson Rule 2：連續 9 點在均值同一側 ───────────────────

class TestRule2:
    # 前 3 個低值把均值往下拉，後 9 個高值全在均值上方
    # mean = (3*3 + 9*10) / 12 ≈ 8.25，後 9 點 (10) > 8.25
    _ABOVE = [3.0] * 3 + [10.0] * 9

    # 前 3 個高值拉高均值，後 9 個低值全在均值下方
    # mean ≈ (30 + 27) / 12 ≈ 4.75，後 9 點 (3) < 4.75
    _BELOW = [10.0] * 3 + [3.0] * 9

    def test_9_consecutive_above_mean_triggers(self):
        r = analyze(self._ABOVE, 20.0, -5.0)
        assert r['is_alarm'] is True
        assert any('Rule2' in v for v in r['violations'])

    def test_9_consecutive_above_mean_is_warning(self):
        r = analyze(self._ABOVE, 20.0, -5.0)
        assert r['severity'] == 'WARNING'

    def test_9_consecutive_below_mean_triggers(self):
        r = analyze(self._BELOW, 20.0, -5.0)
        assert r['is_alarm'] is True
        assert any('Rule2' in v for v in r['violations'])

    def test_8_consecutive_does_not_trigger(self):
        # 4 個低值 + 8 個高值：最長同側序列只有 8
        values = [3.0] * 4 + [10.0] * 8
        r = analyze(values, 20.0, -5.0)
        assert not any('Rule2' in v for v in r['violations'])

    def test_alarm_indices_include_9_violation_points(self):
        r = analyze(self._ABOVE, 20.0, -5.0)
        assert len(r['alarm_indices']) >= 9


# ─── Nelson Rule 3：連續 6 點單調趨勢 ───────────────────────

class TestRule3:
    _INCREASING = [1.0, 2.0, 3.0, 4.0, 5.0, 6.0]
    _DECREASING = [6.0, 5.0, 4.0, 3.0, 2.0, 1.0]

    def test_6_consecutive_increasing_triggers(self):
        r = analyze(self._INCREASING, WIDE_UCL, WIDE_LCL)
        assert r['is_alarm'] is True
        assert any('Rule3' in v for v in r['violations'])

    def test_6_consecutive_increasing_is_warning(self):
        r = analyze(self._INCREASING, WIDE_UCL, WIDE_LCL)
        assert r['severity'] == 'WARNING'

    def test_6_consecutive_decreasing_triggers(self):
        r = analyze(self._DECREASING, WIDE_UCL, WIDE_LCL)
        assert r['is_alarm'] is True
        assert any('Rule3' in v for v in r['violations'])

    def test_5_consecutive_does_not_trigger(self):
        r = analyze([1.0, 2.0, 3.0, 4.0, 5.0], WIDE_UCL, WIDE_LCL)
        assert not any('Rule3' in v for v in r['violations'])

    def test_alarm_indices_are_all_6_points(self):
        r = analyze(self._INCREASING, WIDE_UCL, WIDE_LCL)
        assert r['alarm_indices'] == [0, 1, 2, 3, 4, 5]

    def test_plateau_does_not_trigger(self):
        # 相鄰點相等，非嚴格遞增/遞減
        r = analyze([5.0] * 6, WIDE_UCL, WIDE_LCL)
        assert not any('Rule3' in v for v in r['violations'])


# ─── 多規則同時觸發 ──────────────────────────────────────────

class TestMultipleRules:
    def test_rule1_plus_rule2_severity_is_critical(self):
        # 前 9 個低值 (5) 在均值下方，最後一個 (100) 超出 UCL=50
        # mean = (9*5 + 100) / 10 = 14.5；5 < 14.5 觸發 Rule2
        # 100 > 50 觸發 Rule1 → CRITICAL 優先
        values = [5.0] * 9 + [100.0]
        r = analyze(values, 50.0, -10.0)
        assert r['severity'] == 'CRITICAL'
        assert any('Rule1' in v for v in r['violations'])
        assert any('Rule2' in v for v in r['violations'])

    def test_rule1_plus_rule3_severity_is_critical(self):
        # 前 5 點遞增，最後一點超出 UCL → Rule3 + Rule1 → CRITICAL
        values = [1.0, 2.0, 3.0, 4.0, 5.0, 200.0]
        r = analyze(values, 10.0, -10.0)
        assert r['severity'] == 'CRITICAL'
        assert any('Rule1' in v for v in r['violations'])
        assert any('Rule3' in v for v in r['violations'])

    def test_multiple_violations_all_recorded(self):
        values = [5.0] * 9 + [100.0]
        r = analyze(values, 50.0, -10.0)
        assert len(r['violations']) >= 2
