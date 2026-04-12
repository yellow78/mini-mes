// Package spcclient 提供呼叫 Python SPC 分析引擎的 HTTP client
package spcclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// AnalyzeRequest SPC 分析請求（傳給 Python）
type AnalyzeRequest struct {
	EquipmentID int       `json:"equipment_id"`
	Parameter   string    `json:"parameter"`
	Values      []float64 `json:"values"`
	UCL         float64   `json:"ucl"`
	LCL         float64   `json:"lcl"`
}

// AnalyzeResponse Python SPC 分析結果
type AnalyzeResponse struct {
	EquipmentID  int      `json:"equipment_id"`
	Parameter    string   `json:"parameter"`
	Mean         float64  `json:"mean"`
	Std          float64  `json:"std"`
	UCL          float64  `json:"ucl"`
	LCL          float64  `json:"lcl"`
	IsAlarm      bool     `json:"is_alarm"`
	Violations   []string `json:"violations"`
	AlarmIndices []int    `json:"alarm_indices"`
	Severity     string   `json:"severity"` // "CRITICAL" | "WARNING" | ""
}

// Client Python 分析引擎 HTTP client
type Client struct {
	baseURL    string
	httpClient *http.Client
}

// NewClient 建立 SPC client，baseURL 如 "http://analytics:8001"
func NewClient(baseURL string) *Client {
	return &Client{
		baseURL:    baseURL,
		httpClient: &http.Client{Timeout: 5 * time.Second},
	}
}

// Analyze 呼叫 Python POST /spc/analyze
// 若 Python 服務不可用，回傳 nil, error（呼叫端需自行降級處理）
func (c *Client) Analyze(ctx context.Context, req AnalyzeRequest) (*AnalyzeResponse, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost,
		c.baseURL+"/spc/analyze", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("SPC 引擎無回應: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("SPC 引擎回傳 %d", resp.StatusCode)
	}

	var result AnalyzeResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}
