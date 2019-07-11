package wallet

import (
	"context"

	"github.com/go-kivik/kivik"
)

const paymentDBName = "pharmeum-channel_pharmeumccpayment"

type Balancer interface {
	Balance(address string) (string, error)
}

type BalancerImp struct {
	client *kivik.Client
}

func NewBalancer(client *kivik.Client) Balancer {
	return &BalancerImp{client: client}
}

type wallet struct {
	Balance string `json:"balance"`
}

func (b BalancerImp) Balance(address string) (string, error) {
	ctx := context.TODO()
	db := b.client.DB(ctx, paymentDBName)
	row := db.Get(ctx, address)

	wallet := &wallet{}
	if err := row.ScanDoc(wallet); err != nil {
		return "", err
	}

	return wallet.Balance, nil
}
