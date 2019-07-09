package wallet

import (
	"github.com/Pharmeum/pharmeum-payment-api/db"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
)

type Lister interface {
	List(userID uint64) ([]db.Wallet, error)
}

type ListerImp struct {
	client *channel.Client
	db     *db.DB
}

func NewLister(client *channel.Client, db *db.DB) Lister {
	return &ListerImp{
		client: client,
		db:     db,
	}
}

func (l ListerImp) List(userID uint64) ([]db.Wallet, error) {
	return l.db.UserWallets(userID)
}
