package server

import (
	"context"

	"github.com/Kama001/microservicesusingGo/8.grpc/currency_converter/data"
	protos "github.com/Kama001/microservicesusingGo/8.grpc/currency_converter/protos/currency"
	"k8s.io/klog/v2"
)

type Currency struct {
	l     *klog.Logger
	rates data.ExchangeRates
	protos.UnimplementedCurrencyServer
}

func NewCurrency(l *klog.Logger, r *data.ExchangeRates) *Currency {
	return &Currency{l: l, rates: *r}
}

func (s *Currency) GetRate(ctx context.Context, rr *protos.RateRequest) (*protos.RateResponse, error) {
	base := rr.GetBase()
	destination := rr.GetDestination()
	s.l.Info("Handle GetRate", "base:", base, "destination:", destination)

	// get the rate from the data
	rate, err := s.rates.GetRate(base, destination)
	if err != nil {
		s.l.Error(err, "Cannot get rate for given base and destination")
		return nil, err
	}

	return &protos.RateResponse{Rate: rate}, nil
}
