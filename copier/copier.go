package copier

import (
	"context"
	"time"

	pb "github.com/MetaRPC/GoCopier/copier/copierpb"
)

type CopierService struct {
	Account    *CopierAccount
	grpcClient pb.CopierServiceClient
}

func NewCopierService(endpoint, userKey string) (*CopierService, error) {
	acc, err := NewCopierAccount(endpoint, userKey)
	if err != nil {
		return nil, err
	}
	return &CopierService{
		Account:    acc,
		grpcClient: pb.NewCopierServiceClient(acc.Conn),
	}, nil
}

func (s *CopierService) Start(ctx context.Context, req *StartRequest) (*StartReply, error) {
	userKey := req.UserKey
	if userKey == "" {
		userKey = s.Account.UserKey
	}
	mgrKey := req.ManagerKey
	if mgrKey == "" {
		mgrKey = s.Account.ManagerKey
	}

	callCtx, cancel := context.WithTimeout(s.Account.AuthContext(ctx), 15*time.Second)
	defer cancel()

	pbReq := &pb.StartRequest{
		UserKey:    userKey,
		ManagerKey: mgrKey,
		Master: &pb.Account{
			Type:     req.Master.Type,
			User:     req.Master.User,
			Password: req.Master.Password,
			Server:   req.Master.Server,
			Name:     req.Master.Name,
		},
		Slave: &pb.Account{
			Type:     req.Slave.Type,
			User:     req.Slave.User,
			Password: req.Slave.Password,
			Server:   req.Slave.Server,
			Name:     req.Slave.Name,
		},
		RiskType:           req.RiskType,
		RiskValue:          req.RiskValue,
		FixedMasterBalance: req.FixedMasterBalance,
		CopySl:             req.CopySl,
		CopyTp:             req.CopyTp,
		CopyPendingOrders:  req.CopyPendingOrders,
		ReverseCopy:        req.ReverseCopy,
	}

	reply, err := s.grpcClient.Start(callCtx, pbReq)
	if err != nil {
		return &StartReply{Ok: false, Error: err.Error()}, nil
	}
	return &StartReply{Ok: reply.Ok, CopierId: reply.CopierId, Error: reply.Error}, nil
}

func (s *CopierService) List(ctx context.Context) (*ListReply, error) {
	callCtx, cancel := context.WithTimeout(s.Account.AuthContext(ctx), 10*time.Second)
	defer cancel()

	reply, err := s.grpcClient.List(callCtx, &pb.ListRequest{UserKey: s.Account.UserKey})
	if err != nil {
		return &ListReply{Ok: false, Copiers: []*CopierSummary{}, Error: err.Error()}, nil
	}

	var copiers []*CopierSummary
	for _, c := range reply.Copiers {
		copiers = append(copiers, &CopierSummary{
			Id:           c.Id,
			MasterType:   c.MasterType,
			MasterUser:   c.MasterUser,
			MasterServer: c.MasterServer,
			SlaveType:    c.SlaveType,
			SlaveUser:    c.SlaveUser,
			SlaveServer:  c.SlaveServer,
			RiskType:     c.RiskType,
			RiskValue:    c.RiskValue,
			Paused:       c.Paused,
			PauseReason:  c.PauseReason,
		})
	}
	return &ListReply{Ok: reply.Ok, Copiers: copiers, Error: reply.Error}, nil
}

func (s *CopierService) Pause(ctx context.Context, copierId string, paused bool) (*SimpleReply, error) {
	callCtx, cancel := context.WithTimeout(s.Account.AuthContext(ctx), 10*time.Second)
	defer cancel()

	reply, err := s.grpcClient.Pause(callCtx, &pb.PauseRequest{UserKey: s.Account.UserKey, CopierId: copierId, Paused: paused})
	if err != nil {
		return &SimpleReply{Ok: false, Error: err.Error()}, nil
	}
	return &SimpleReply{Ok: reply.Ok, Error: reply.Error}, nil
}

func (s *CopierService) Remove(ctx context.Context, copierId string) (*SimpleReply, error) {
	callCtx, cancel := context.WithTimeout(s.Account.AuthContext(ctx), 10*time.Second)
	defer cancel()

	reply, err := s.grpcClient.Remove(callCtx, &pb.RemoveRequest{UserKey: s.Account.UserKey, CopierId: copierId})
	if err != nil {
		return &SimpleReply{Ok: false, Error: err.Error()}, nil
	}
	return &SimpleReply{Ok: reply.Ok, Error: reply.Error}, nil
}

func (s *CopierService) Close() error {
	return s.Account.Close()
}
