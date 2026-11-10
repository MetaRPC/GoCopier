package copier

import (
	"context"
	"crypto/tls"
	"fmt"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

type CopierAccount struct {
	Endpoint   string
	UserKey    string
	ManagerKey string
	Conn       *grpc.ClientConn
}

func NewCopierAccount(endpoint, userKey string) (*CopierAccount, error) {
	if userKey == "" || strings.EqualFold(userKey, "TRIAL") || strings.IndexFunc(userKey, func(r rune) bool { return r < 33 || r > 126 }) >= 0 {
        return nil, fmt.Errorf("a real owner API key is required (printable ASCII, no whitespace)")
    }
	clean := strings.TrimPrefix(endpoint, "https://")
	clean = strings.TrimPrefix(clean, "http://")
	clean = strings.TrimSuffix(clean, "/")
	if !strings.Contains(clean, ":") {
		clean = clean + ":443"
	}
	creds := credentials.NewTLS(&tls.Config{})
	conn, err := grpc.Dial(clean, grpc.WithTransportCredentials(creds))
	if err != nil {
		return nil, err
	}
	return &CopierAccount{
		Endpoint:   clean,
		UserKey:    userKey,
		ManagerKey: userKey,
		Conn:       conn,
	}, nil
}

func (a *CopierAccount) AuthContext(ctx context.Context) context.Context {
	md, _ := metadata.FromOutgoingContext(ctx)
	md = md.Copy()
	for key, values := range metadata.Pairs(
		"apikey", a.UserKey,
		"x-metarpc-client-sdk", "GoCopier/1.0.0",
	) { md.Set(key, values...) }
	return metadata.NewOutgoingContext(ctx, md)
}

func (a *CopierAccount) Close() error {
	if a.Conn != nil {
		return a.Conn.Close()
	}
	return nil
}
