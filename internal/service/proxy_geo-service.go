package service

import (
	"context"
	"errors"
	"github.com/RVodassa/geo-microservices-proxy/internal/domain/entity"
	"github.com/RVodassa/geo-microservices-proxy/internal/infrastructure/cache"
	grpc_auth_service "github.com/RVodassa/geo-microservices-proxy/internal/infrastructure/grpc_auth-service"
	grpc_geo_service "github.com/RVodassa/geo-microservices-proxy/internal/infrastructure/grpc_geo-service"
	grpc_user_service "github.com/RVodassa/geo-microservices-proxy/internal/infrastructure/grpc_user-service"
	"github.com/go-redis/redis/v8"
	"log"
)

type ProxyGeoServiceProvider interface {
	Register(ctx context.Context, auth *entity.RegisterRequest) (uint64, error)
	Login(ctx context.Context, auth *entity.LoginRequest) (string, error)
	Geocode(ctx context.Context, geo *entity.GeocodeRequest) (entity.ResponseBody, error)
	Search(ctx context.Context, geo *entity.SearchRequest) (entity.ResponseBody, error)
	CheckToken(ctx context.Context, token *entity.CheckTokenRequest) (bool, error)
	ListProfiles(ctx context.Context, request *entity.ListRequest) (*entity.ListResponse, error)
	Profile(ctx context.Context, request *entity.ProfileRequest) (*entity.ProfileResponse, error)
}

type ProxyGeoService struct {
	cache       cache.CacheServiceProvider
	authService *grpc_auth_service.AuthServiceClient
	geoService  *grpc_geo_service.GeoServiceClient
	userService *grpc_user_service.UserServiceClient
}

func NewProxyGeoService(
	cache cache.CacheServiceProvider,
	authService *grpc_auth_service.AuthServiceClient,
	geoService *grpc_geo_service.GeoServiceClient,
	userService *grpc_user_service.UserServiceClient,
) *ProxyGeoService {

	return &ProxyGeoService{
		cache:       cache,
		authService: authService,
		geoService:  geoService,
		userService: userService,
	}
}

func (p *ProxyGeoService) Register(ctx context.Context, auth *entity.RegisterRequest) (uint64, error) {
	if auth.Password == "" || auth.Login == "" {
		log.Println("Password and Login are required fields")
		return 0, errors.New("password and Login are required fields")
	}

	id, err := p.authService.Register(ctx, auth)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (p *ProxyGeoService) Login(ctx context.Context, auth *entity.LoginRequest) (string, error) {
	if auth.Password == "" || auth.Login == "" {
		log.Println("Password and Login are required fields")
		return "", errors.New("password and Login are required fields")
	}

	tokenJWT, err := p.authService.Login(ctx, auth)
	if err != nil {
		return "", err
	}

	return tokenJWT, nil
}

func (p *ProxyGeoService) Geocode(ctx context.Context, geo *entity.GeocodeRequest) (entity.ResponseBody, error) {
	var resp entity.ResponseBody

	if geo.Lat == "" || geo.Lng == "" {
		log.Println("Lat and Lng are required fields")
		return entity.ResponseBody{}, errors.New("отсутствуют необходимые данные")
	}

	key := geo.Lat + "," + geo.Lng
	resp, err := p.cache.Get(ctx, key)

	if err != nil {
		if errors.Is(err, redis.Nil) {
			resp, err = p.geoService.GeoCode(ctx, geo)
			if err != nil {
				log.Printf("Ошибка при получении данных из базы данных: %v", err)
				return entity.ResponseBody{}, err
			}
			log.Printf("Данные для %s взяты из базы данных", key)

			err = p.cache.Set(ctx, key, resp.Addresses)
			if err != nil {
				log.Printf("Ошибка при сохранении данных в кэш: %v", err)
			}
			return resp, nil
		}

		log.Printf("Ошибка при получении данных из кэша: %v", err)
		return entity.ResponseBody{}, err
	}

	log.Printf("Данные для %s взяты из кэша", key)
	return resp, nil
}

func (p *ProxyGeoService) Search(ctx context.Context, geo *entity.SearchRequest) (entity.ResponseBody, error) {
	var resp entity.ResponseBody

	if geo.Query == "" {
		log.Println("query are required fields")
		return entity.ResponseBody{}, errors.New("отсутствуют необходимые данные")
	}

	key := geo.Query
	resp, err := p.cache.Get(ctx, key)

	if err != nil {
		if errors.Is(err, redis.Nil) {

			resp, err = p.geoService.Search(ctx, geo)
			if err != nil {
				return entity.ResponseBody{}, err
			}
			err = p.cache.Set(ctx, key, resp.Addresses)
			if err != nil {
				log.Printf("Ошибка при сохранении данных в кэш: %v", err)
			}

			log.Printf("Данные для %s взяты из базы данных", key)
			return resp, nil
		}

		log.Printf("Ошибка при получении данных из кэша: %v", err)
		return entity.ResponseBody{}, err
	}

	log.Printf("Данные для %s взяты из кэша", key)
	return resp, nil
}

func (p *ProxyGeoService) CheckToken(ctx context.Context, token *entity.CheckTokenRequest) (bool, error) {
	if token.Token == "" {
		log.Println("token is required")
		return false, errors.New("token is required")
	}
	status, err := p.authService.CheckToken(ctx, token)
	if err != nil {
		log.Println(err)
		return false, err
	}
	return status, nil
}

func (p *ProxyGeoService) ListProfiles(ctx context.Context, request *entity.ListRequest) (*entity.ListResponse, error) {
	list, err := p.userService.ListProfiles(ctx, request)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return list, nil
}

func (p *ProxyGeoService) Profile(ctx context.Context, request *entity.ProfileRequest) (*entity.ProfileResponse, error) {

	resp, err := p.userService.Profile(ctx, request)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return resp, nil
}
