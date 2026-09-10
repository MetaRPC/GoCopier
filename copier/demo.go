package copier

import (
	"context"
	"math/rand"
	"time"
)

type DemoAccountClient struct {
	Endpoint string
}

func NewDemoAccountClient(endpoint string) (*DemoAccountClient, error) {
	return &DemoAccountClient{Endpoint: endpoint}, nil
}

func (d *DemoAccountClient) OpenDemoAccount(ctx context.Context, req *GuiDemoOpenAccountRequest) (*GuiDemoOpenAccountReply, error) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	login := uint64(r.Intn(899999) + 100000)
	return &GuiDemoOpenAccountReply{
		ResultCode: 0,
		Login:      login,
		Password:   "Demo1234!",
		Investor:   "Inv1234!",
		Server:     req.Server,
	}, nil
}
