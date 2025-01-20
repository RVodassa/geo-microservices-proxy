package grpc_auth_service

import (
	"context"
	proto "github.com/RVodassa/geo-microservices-auth_service/proto/generated"
	"github.com/RVodassa/geo-microservices-proxy/internal/domain/entity"
	"google.golang.org/grpc"
)

type AuthServiceClient struct {
	conn   *grpc.ClientConn
	client proto.AuthServiceClient
}

func NewAuthServiceClient(conn *grpc.ClientConn) *AuthServiceClient {
	client := proto.NewAuthServiceClient(conn)

	return &AuthServiceClient{
		client: client,
		conn:   conn,
	}
}

func (a *AuthServiceClient) Login(ctx context.Context, auth *entity.LoginRequest) (jwtToken string, err error) {

	req := &proto.LoginRequest{
		Login:    auth.Login,
		Password: auth.Password,
	}

	resp, err := a.client.Login(ctx, req)
	if err != nil {
		return "", err
	}
	return resp.Token, nil
}

func (a *AuthServiceClient) Register(ctx context.Context, auth *entity.RegisterRequest) (id uint64, err error) {

	req := &proto.RegisterRequest{
		Login:    auth.Login,
		Password: auth.Password,
	}

	resp, err := a.client.Register(ctx, req)
	if err != nil {
		return 0, err
	}

	return resp.Id, nil
}

func (a *AuthServiceClient) CheckToken(ctx context.Context, token *entity.CheckTokenRequest) (bool, error) {

	req := &proto.CheckTokenRequest{
		Token: token.Token,
	}

	verify, err := a.client.CheckToken(ctx, req)
	if err != nil {
		return false, err
	}
	return verify.Status, nil
}
