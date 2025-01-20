package grpc_geo_service

import (
	"context"
	proto "github.com/RVodassa/geo-microservices-geo_service/proto/generated"
	"github.com/RVodassa/geo-microservices-proxy/internal/domain/entity"
	"google.golang.org/grpc"
	"log"
)

type GeoServiceClient struct {
	conn   *grpc.ClientConn
	client proto.GeoServiceClient
}

func NewGeoServiceClient(conn *grpc.ClientConn) *GeoServiceClient {
	client := proto.NewGeoServiceClient(conn)

	return &GeoServiceClient{
		client: client,
		conn:   conn,
	}
}

func (g *GeoServiceClient) Search(ctx context.Context, search *entity.SearchRequest) (entity.ResponseBody, error) {
	var addresses []*entity.Address

	req := &proto.SearchRequest{
		Input: search.Query,
	}

	result, err := g.client.Search(ctx, req)
	if err != nil {
		return entity.ResponseBody{}, err
	}

	for _, address := range result.Addresses {
		addr := entity.Address{
			City:   address.City,
			Street: address.Street,
			House:  address.House,
			Lat:    address.Lat,
			Lon:    address.Lon,
		}
		addresses = append(addresses, &addr)
	}
	response := entity.ResponseBody{
		Addresses: addresses,
	}

	return response, nil
}

func (g *GeoServiceClient) GeoCode(ctx context.Context, geocode *entity.GeocodeRequest) (entity.ResponseBody, error) {

	req := &proto.GeoCodeRequest{
		Lat: geocode.Lat,
		Lng: geocode.Lng,
	}

	result, err := g.client.GeoCode(ctx, req)
	if err != nil {
		log.Printf("ошибка при grpc запросе в geo-service")
		return entity.ResponseBody{}, err
	}

	var addresses []*entity.Address
	for _, address := range result.Addresses {
		addr := entity.Address{
			City:   address.City,
			Street: address.Street,
			House:  address.House,
			Lat:    address.Lat,
			Lon:    address.Lon,
		}
		addresses = append(addresses, &addr)
	}
	response := entity.ResponseBody{
		Addresses: addresses,
	}

	return response, nil
}
