package grpc_geo_service

import (
	"context"
	"errors"
	"github.com/RVodassa/geo-microservices-proxy/internal/domain/entity"
	proto "github.com/RVodassa/geo-microservices-user_service/proto/generated"
	"google.golang.org/grpc"
	"log"
)

type UserServiceClient struct {
	conn   *grpc.ClientConn
	client proto.UserServiceClient
}

func NewUserServiceClient(conn *grpc.ClientConn) *UserServiceClient {
	// Создание клиента
	client := proto.NewUserServiceClient(conn)

	return &UserServiceClient{
		client: client,
		conn:   conn,
	}
}

func (u *UserServiceClient) ListProfiles(ctx context.Context, request *entity.ListRequest) (*entity.ListResponse, error) {
	req := &proto.ListRequest{
		Offset: request.Offset,
		Limit:  request.Limit,
	}

	list, err := u.client.List(ctx, req)
	if err != nil {
		log.Println("Ошибка при обращении к клиенту user")
		return nil, err
	}

	var profiles []*entity.ProfileResponse

	for _, user := range list.Users {
		profile := &entity.ProfileResponse{
			ID:    user.Id,
			Login: user.Login,
		}
		profiles = append(profiles, profile)
	}
	if len(profiles) == 0 {
		return &entity.ListResponse{}, errors.New("users not found")
	}

	response := &entity.ListResponse{
		Profiles: profiles,
		Count:    list.Total,
	}

	return response, nil
}
func (u *UserServiceClient) Profile(ctx context.Context, request *entity.ProfileRequest) (*entity.ProfileResponse, error) {
	req := &proto.ProfileRequest{
		Id: request.ID,
	}
	resp, err := u.client.Profile(ctx, req)
	if err != nil {
		return nil, err
	}

	response := &entity.ProfileResponse{
		ID:    resp.Id,
		Login: resp.Login,
	}
	return response, nil
}
