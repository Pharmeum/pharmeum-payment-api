package payment

import "github.com/Pharmeum/pharmeum-payment-api/db"

type OperationType uint8

const (
	CreateWallet OperationType = iota + 1
)

type Uploader struct {
	Operation OperationType
	Wallet    *db.Wallet
}
