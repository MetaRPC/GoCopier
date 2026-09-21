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
