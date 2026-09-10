package copier

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type CopierAccount struct {
	Endpoint   string
	UserKey    string
	ManagerKey string
	Conn       *grpc.ClientConn
}

func NewCopierAccount(endpoint, userKey, managerKey string) (*CopierAccount, error) {
	if managerKey == "" {
		managerKey = userKey
	}
	return &CopierAccount{
		Endpoint:   endpoint,
		UserKey:    userKey,
		ManagerKey: managerKey,
	}, nil
}

func (a *CopierAccount) AuthContext(ctx context.Context) context.Context {
	md := metadata.Pairs(
		"authorization", "Bearer "+a.UserKey,
		"x-metarpc-manager", a.ManagerKey,
		"x-metarpc-client-sdk", "GoCopier/1.0.0",
	)
	return metadata.NewOutgoingContext(ctx, md)
}

func (a *CopierAccount) Close() error {
	if a.Conn != nil {
		return a.Conn.Close()
	}
	return nil
}
