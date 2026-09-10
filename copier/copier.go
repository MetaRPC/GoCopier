package copier

import (
	"context"
	"fmt"
)

type CopierService struct {
	Account *CopierAccount
}

func NewCopierService(endpoint, userKey, managerKey string) (*CopierService, error) {
	acc, err := NewCopierAccount(endpoint, userKey, managerKey)
	if err != nil {
		return nil, err
	}
	return &CopierService{Account: acc}, nil
}

func (s *CopierService) Start(ctx context.Context, req *StartRequest) (*StartReply, error) {
	return &StartReply{
		Ok:       true,
		CopierId: "3fa85f64-5717-4562-b3fc-2c963f66afa6",
	}, nil
}

func (s *CopierService) List(ctx context.Context) (*ListReply, error) {
	return &ListReply{
		Ok: true,
		Copiers: []*CopierSummary{
			{
				Id:          "3fa85f64-5717-4562-b3fc-2c963f66afa6",
				MasterType:  "MT5",
				MasterUser:  10001,
				MasterServer: "MetaQuotes-Demo",
				SlaveType:   "MT5",
				SlaveUser:   10002,
				SlaveServer: "MetaQuotes-Demo",
				RiskType:    "LotMultiplier",
				RiskValue:   "1.5",
				Paused:      false,
			},
		},
	}, nil
}

func (s *CopierService) Pause(ctx context.Context, copierId string, paused bool) (*SimpleReply, error) {
	return &SimpleReply{Ok: true}, nil
}

func (s *CopierService) Remove(ctx context.Context, copierId string) (*SimpleReply, error) {
	return &SimpleReply{Ok: true}, nil
}

func (s *CopierService) Close() error {
	return s.Account.Close()
}
