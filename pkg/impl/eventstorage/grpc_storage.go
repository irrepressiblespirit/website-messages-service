package eventstorage

import (
	"context"
	"fmt"
	"log"

	"github.com/irrepressiblespirit/website-messages-service/pkg/centrifugopb"
	"google.golang.org/grpc"
)

type GrpcEventStorage struct {
	Debug   bool
	grpcURL string
	conn    *grpc.ClientConn
	client  centrifugopb.CentrifugoApiClient
}

func NewGrpcEventStorage(grpcURL string) *GrpcEventStorage {
	return &GrpcEventStorage{
		Debug:   false,
		grpcURL: grpcURL,
	}
}

func (s *GrpcEventStorage) Start() error {
	log.Printf("Connecting to centrifugo %s\n", s.grpcURL)
	conn, err := grpc.Dial(s.grpcURL, grpc.WithInsecure())
	if err != nil {
		return fmt.Errorf("error EventStorageGrpcImpl.Start %w", err)
	}
	s.conn = conn
	s.client = centrifugopb.NewCentrifugoApiClient(conn)
	return nil
}

func (s *GrpcEventStorage) Close() {
	log.Println("Disonnect from centrifugo")
	s.conn.Close()
}

func (s *GrpcEventStorage) SendMessage(ctx context.Context, chanel string, body []byte) error {
	if s.Debug {
		log.Printf("%s <- %s", chanel, string(body))
	}
	_, err := s.client.Publish(ctx, &centrifugopb.PublishRequest{
		Channel: chanel,
		Data:    body,
	})
	if err != nil {
		return fmt.Errorf("error in EventStorageGrpcImpl.SendMessage: %w", err)
	}
	return nil
}
