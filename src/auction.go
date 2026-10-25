package main

import (
	"context"
	"sync"
)

type Bid struct {
	ProviderName string
	Price        float64
}

type BidProvider interface {
	GetBid(ctx context.Context) (Bid, error)
}

func FindMaxBid(bids []Bid)(Bid, bool){
	if len(bids) == 0 {
		return Bid{}, false
	}
	MaxPrice := bids[0]
	found := true
	for _, p := range bids[1:]{
		if p.Price > MaxPrice.Price{
			MaxPrice = p
		}
	}	
	return MaxPrice, found
}

func RunAuction(ctx context.Context, providers []BidProvider) (Bid, bool) {
	result := make(chan Bid, len(providers))
	var wg sync.WaitGroup

	for _, provider := range providers {
		wg.Add(1)
		go func(p BidProvider) {
			defer wg.Done()
			if bid, err := p.GetBid(ctx); err == nil {
				result <- bid
			}
		}(provider)
	}

	go func() {
		wg.Wait()
		close(result)
	}()

	var bestBid Bid
	found := false

	for {
		select {
		case <-ctx.Done():
		
			return bestBid, found

		case bid, ok := <-result:
			if !ok {
			
				return bestBid, found
			}
			if !found || bid.Price > bestBid.Price {
				bestBid = bid
				found = true
			}
		}
	}
}