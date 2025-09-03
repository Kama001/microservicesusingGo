package server

import (
	"context"
	"log"

	protoc "github.com/Kama001/microservicesusingGo/8.grpc/currency_converter/protos/currency"
)

type Currency struct {
	l *log.Logger
	protoc.UnimplementedCurrencyServer
}

func NewCurrency(l *log.Logger) *Currency {
	return &Currency{l: l}
}

func (s *Currency) GetRate(ctx context.Context, rr *protoc.RateRequest) (*protoc.RateResponse, error) {
	s.l.Println("Handle GetRate", "base:", rr.GetBase(), "destination:", rr.GetDestination())
	return &protoc.RateResponse{Rate: 0.5}, nil
}
