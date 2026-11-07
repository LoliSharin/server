package main

import (
	"context"
	"fmt"
	"testing"
	"time"
)

type MockProviders struct {
	Name    string
	Price   float64
	Latency time.Duration
}

func (m MockProviders) GetBid(ctx context.Context) (Bid, error) {
	time.Sleep(m.Latency)
	return Bid{ProviderName: m.Name, Price: m.Price}, nil
}


type MockGRPCClient struct {
	Name    string
	Price   float64
	Latency time.Duration
}

func (m MockGRPCClient) GetBid(ctx context.Context, in *BidRequest) (*BidResponse, error) {
	time.Sleep(m.Latency)
	return &BidResponse{ProviderName: m.Name, Price: m.Price}, nil
}

func TestProvidersMock(t *testing.T) {
	tests := []BidProvider{
		MockProviders{Name: "Yandex", Price: 123.12, Latency: 12 * time.Millisecond},
		MockProviders{Name: "VK", Price: 263.13, Latency: 20 * time.Millisecond},
		MockProviders{Name: "Max", Price: 139.56, Latency: 21 * time.Millisecond},
		MockProviders{Name: "Sber", Price: 203.11, Latency: 8 * time.Millisecond},
		MockProviders{Name: "Apple", Price: 637.2, Latency: 100 * time.Millisecond},
		MockProviders{Name: "Microsoft", Price: 612, Latency: 40 * time.Millisecond},
		&GRPCProvider{
			VideoID: "video_777",
			client:  MockGRPCClient{Name: "Google_gRPC", Price: 750.0, Latency: 15 * time.Millisecond},
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	bestBid, found := RunAuction(ctx, tests)
	if bestBid.ProviderName != "Google_gRPC" {
		t.Errorf("Тест сломался, должен выиграть Google_gRPC, но победил %v", bestBid)
	}
	fmt.Printf("Found: %v, Best Bid: %+v\n", found, bestBid)
}