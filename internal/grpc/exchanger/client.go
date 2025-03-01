package exchanger

import (
	"context"
	pb "github.com/galkin09/proto-exchange/exchange"
	"google.golang.org/grpc"
	"os"
)

type ExchangerClient struct {
	pb.ExchangeServiceClient
}

func NewExchangerClient() *ExchangerClient {
	addr := os.Getenv("GRPC_ADDR")

	conn, err := grpc.NewClient(addr)
	if err != nil {
		panic(err)
	}

	exchClient := pb.NewExchangeServiceClient(conn)

	return &ExchangerClient{exchClient}
}

func (e *ExchangerClient) GetExchangeRates(ctx context.Context, in *pb.Empty) (*pb.ExchangeRatesResponse, error) {
	return e.ExchangeServiceClient.GetExchangeRates(ctx, in)
}

func (e *ExchangerClient) GetExchangeRateForCurrency(ctx context.Context, in *pb.CurrencyRequest) (*pb.ExchangeRateResponse, error) {
	return e.ExchangeServiceClient.GetExchangeRateForCurrency(ctx, in)
}
