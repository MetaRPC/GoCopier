package copier

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type DemoAccountClient struct {
	Endpoint string
	client   *http.Client
}

func NewDemoAccountClient(endpoint string) (*DemoAccountClient, error) {
	clean := strings.TrimPrefix(endpoint, "http://")
	clean = strings.TrimPrefix(clean, "https://")
	clean = strings.TrimSuffix(clean, ":443")
	clean = strings.TrimSuffix(clean, "/")
	return &DemoAccountClient{
		Endpoint: "https://" + clean,
		client:   &http.Client{Timeout: 120 * time.Second},
	}, nil
}

func (d *DemoAccountClient) OpenDemoAccount(ctx context.Context, req *GuiDemoOpenAccountRequest, optApiKey ...string) (*GuiDemoOpenAccountReply, error) {
	apiKey := "TRIAL"
	if len(optApiKey) > 0 && optApiKey[0] != "" {
		apiKey = optApiKey[0]
	}

	reqUrl := fmt.Sprintf("%s/DemoAccount/Open?server=%s", d.Endpoint, url.QueryEscape(req.Server))
	httpReq, err := http.NewRequestWithContext(ctx, "GET", reqUrl, nil)
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("APIKey", apiKey)
	httpReq.Header.Set("User-Agent", "GoCopier/1.0.0")

	resp, err := d.client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("DemoAccount/Open failed with HTTP %d: %s", resp.StatusCode, string(body))
	}

	var raw struct {
		ResultCode int    `json:"resultCode"`
		Login      string `json:"login"`
		Password   string `json:"password"`
		Investor   string `json:"investor"`
		Server     string `json:"server"`
	}
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, err
	}

	loginNum, _ := strconv.ParseUint(raw.Login, 10, 64)
	return &GuiDemoOpenAccountReply{
		ResultCode: int32(raw.ResultCode),
		Login:      loginNum,
		Password:   raw.Password,
		Investor:   raw.Investor,
		Server:     raw.Server,
	}, nil
}

func (d *DemoAccountClient) ConnectEx(ctx context.Context, user uint64, password, server string, optApiKey ...string) (*ConnectExReply, error) {
	apiKey := "TRIAL"
	if len(optApiKey) > 0 && optApiKey[0] != "" {
		apiKey = optApiKey[0]
	}

	params := url.Values{}
	params.Set("user", strconv.FormatUint(user, 10))
	params.Set("password", password)
	params.Set("mtClusterName", server)

	reqUrl := fmt.Sprintf("%s/ConnectEx?%s", d.Endpoint, params.Encode())
	httpReq, err := http.NewRequestWithContext(ctx, "GET", reqUrl, nil)
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("APIKey", apiKey)
	httpReq.Header.Set("User-Agent", "GoCopier/1.0.0")

	resp, err := d.client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ConnectEx failed with HTTP %d: %s", resp.StatusCode, string(body))
	}

	var raw struct {
		Data struct {
			TerminalInstanceGuid string `json:"terminalInstanceGuid"`
			TerminalType         string `json:"terminalType"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, err
	}

	return &ConnectExReply{
		TerminalInstanceGuid: raw.Data.TerminalInstanceGuid,
		TerminalType:         raw.Data.TerminalType,
	}, nil
}

func (d *DemoAccountClient) Disconnect(ctx context.Context, terminalId string, optApiKey ...string) (*DisconnectReply, error) {
	apiKey := "TRIAL"
	if len(optApiKey) > 0 && optApiKey[0] != "" {
		apiKey = optApiKey[0]
	}

	reqUrl := fmt.Sprintf("%s/Disconnect", d.Endpoint)
	httpReq, err := http.NewRequestWithContext(ctx, "GET", reqUrl, nil)
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("APIKey", apiKey)
	httpReq.Header.Set("id", terminalId)
	httpReq.Header.Set("User-Agent", "GoCopier/1.0.0")

	resp, err := d.client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Disconnect failed with HTTP %d: %s", resp.StatusCode, string(body))
	}

	var raw struct {
		Data struct {
			UniqueIdentifier    string `json:"uniqueIdentifier"`
			FullLifeTimeSeconds int    `json:"fullLifeTimeSeconds"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, err
	}

	return &DisconnectReply{
		UniqueIdentifier:    raw.Data.UniqueIdentifier,
		FullLifeTimeSeconds: raw.Data.FullLifeTimeSeconds,
	}, nil
}

type PositionInfo struct {
	Ticket uint64  `json:"ticket"`
	Symbol string  `json:"symbol"`
	Volume float64 `json:"volume"`
	Type   string  `json:"type"`
}

func ToHyphenGuid(guid string) string {
	clean := strings.TrimPrefix(guid, "mt5_live_")
	clean = strings.ReplaceAll(clean, "-", "")
	if len(clean) < 32 {
		return guid
	}
	return fmt.Sprintf("%s-%s-%s-%s-%s", clean[0:8], clean[8:12], clean[12:16], clean[16:20], clean[20:32])
}

func (d *DemoAccountClient) OrderSend(ctx context.Context, terminalId, symbol, operation string, volume float64, optApiKey ...string) (uint64, error) {
	apiKey := "TRIAL"
	if len(optApiKey) > 0 && optApiKey[0] != "" {
		apiKey = optApiKey[0]
	}

	params := url.Values{}
	params.Set("id", terminalId)
	params.Set("symbol", symbol)
	params.Set("operation", operation)
	params.Set("volume", fmt.Sprintf("%.2f", volume))

	reqUrl := fmt.Sprintf("%s/OrderSend?%s", d.Endpoint, params.Encode())
	httpReq, err := http.NewRequestWithContext(ctx, "GET", reqUrl, nil)
	if err != nil {
		return 0, err
	}
	httpReq.Header.Set("APIKey", apiKey)
	httpReq.Header.Set("id", terminalId)
	httpReq.Header.Set("User-Agent", "GoCopier/1.0.0")

	resp, err := d.client.Do(httpReq)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("OrderSend failed with HTTP %d: %s", resp.StatusCode, string(body))
	}

	var raw struct {
		Data struct {
			Order  uint64 `json:"order"`
			Ticket uint64 `json:"ticket"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &raw); err != nil {
		return 0, err
	}
	if raw.Data.Order != 0 {
		return raw.Data.Order, nil
	}
	return raw.Data.Ticket, nil
}

func (d *DemoAccountClient) OpenedOrders(ctx context.Context, terminalId string, optApiKey ...string) ([]*PositionInfo, error) {
	apiKey := "TRIAL"
	if len(optApiKey) > 0 && optApiKey[0] != "" {
		apiKey = optApiKey[0]
	}

	reqUrl := fmt.Sprintf("%s/OpenedOrders?id=%s", d.Endpoint, url.QueryEscape(terminalId))
	httpReq, err := http.NewRequestWithContext(ctx, "GET", reqUrl, nil)
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("APIKey", apiKey)
	httpReq.Header.Set("id", terminalId)
	httpReq.Header.Set("User-Agent", "GoCopier/1.0.0")

	resp, err := d.client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("OpenedOrders failed with HTTP %d: %s", resp.StatusCode, string(body))
	}

	var raw struct {
		Data struct {
			PositionInfos []*PositionInfo `json:"positionInfos"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, err
	}

	return raw.Data.PositionInfos, nil
}

func (d *DemoAccountClient) OrderClose(ctx context.Context, terminalId string, ticket uint64, optApiKey ...string) (string, error) {
	apiKey := "TRIAL"
	if len(optApiKey) > 0 && optApiKey[0] != "" {
		apiKey = optApiKey[0]
	}

	params := url.Values{}
	params.Set("id", terminalId)
	params.Set("ticket", strconv.FormatUint(ticket, 10))
	params.Set("volume", "0")
	params.Set("slippage", "20")

	reqUrl := fmt.Sprintf("%s/OrderClose?%s", d.Endpoint, params.Encode())
	httpReq, err := http.NewRequestWithContext(ctx, "GET", reqUrl, nil)
	if err != nil {
		return "", err
	}
	httpReq.Header.Set("APIKey", apiKey)
	httpReq.Header.Set("id", terminalId)
	httpReq.Header.Set("User-Agent", "GoCopier/1.0.0")

	resp, err := d.client.Do(httpReq)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("OrderClose failed with HTTP %d: %s", resp.StatusCode, string(body))
	}

	return string(body), nil
}

