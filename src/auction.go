package main

import (
	"context"
)

type bidAction struct{
	providerName string
	price float64
}
type bidProvider interface{
	Name() string
	GetBid(ctx context.Context) (bidAction, error) 
}

func RunAction (ctx context.Context) (bidAction, error){

	

}