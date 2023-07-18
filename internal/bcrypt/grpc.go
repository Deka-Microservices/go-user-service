package bcrypt

import (
	"context"
	"fmt"

	"github.com/deka-microservices/go-bcrypt-service/pkg/service"
	"github.com/deka-microservices/go-user-service/internal/consts"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type grpcBcryptClient struct {
	conn    *grpc.ClientConn
	service service.BcryptServiceClient
}

func NewGrpcBcryptServiceClient() (BcryptClient, error) {

	var options []grpc.DialOption
	options = append(options, grpc.WithTransportCredentials(insecure.NewCredentials()))

	address := viper.GetString(consts.CONFIG_BCRYPT_SERVER_ADDRESS)
	conn, err := grpc.Dial(address, options...)
	if err != nil {
		log.Error().Err(err).Msg("bcrypt service dial fail")
		return nil, fmt.Errorf("bcrypt service dial fail: %w", err)
	}

	service := service.NewBcryptServiceClient(conn)

	return &grpcBcryptClient{conn, service}, nil
}

func (client *grpcBcryptClient) HashPassword(ctx context.Context, password string) (string, error) {
	hash_resp, err := client.service.HashPassword(ctx, &service.HashRequest{Password: password})

	if err != nil {
		return "", err
	}

	return hash_resp.Hash, nil

}

func (client *grpcBcryptClient) CheckPassword(ctx context.Context, password string, hash string) (bool, error) {
	resp, err := client.service.CheckPassword(ctx, &service.CheckRequest{Password: password, Hash: hash})
	if err != nil {
		return false, err
	}

	return resp.Valid, nil
}

func (client *grpcBcryptClient) Close() {
	client.conn.Close()
}
