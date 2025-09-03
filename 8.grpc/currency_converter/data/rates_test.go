package data

import (
	"fmt"
	"testing"

	"k8s.io/klog/v2"
)

func TestGetRates(t *testing.T) {
	l := klog.NewKlogr()
	rates := NewExchangeRates(&l)
	err := rates.getRates()
	if err != nil {
		t.Errorf("Failed to get rates: %v", err)
	}
	fmt.Printf("%+v", rates)
}
