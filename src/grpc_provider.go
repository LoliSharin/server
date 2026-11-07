package main

import (
	"context"

	"google.golang.org/grpc"
)

type GRPCProvider struct{
	client BidderServiceClient
	VideoID string
}

type BidRequest struct {
	VideoID string
}
type BidResponse struct {
	ProviderName string
	Price float64
}

type BidderServiceClient interface{
	GetBid(ctx context.Context, in *BidRequest, opt ...grpc.PeerCallOption) (*BidResponse, error)
}

func (g *GRPCProvider)GetBid(ctx context.Context)(Bid, error){
	resp, err := g.client.GetBid(ctx, &BidRequest{VideoID: g.VideoID})
	if err != nil{
		return Bid{}, err
	}
	return Bid{resp.ProviderName, resp.Price}, nil
}