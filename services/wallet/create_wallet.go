package wallet

import (
	"github.com/Pharmeum/pharmeum-payment-api/db"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/pkg/errors"
)

const (
	paymentChaincodeID = "pharmeumccpayment"
	createWallet       = "create_wallet"
	transferPayment    = "transfer_payment"
)

type Creator interface {
	Create(owner uint64) error
}

type creatorImp struct {
	client *channel.Client
	db     *db.DB
}

func NewCreator(client *channel.Client, db *db.DB) Creator {
	return &creatorImp{
		client: client,
		db:     db,
	}
}

//CreateWallet create wallet in Hyperledger Fabric and PostgreSQL
func (c creatorImp) Create(owner uint64) error {
	//generate unique public key
	walletAddress, err := newAddress()
	if err != nil {
		return errors.Wrap(err, "failed to generate new user address")
	}

	err = c.sendTransaction(walletAddress)
	if err != nil {
		return errors.Wrap(err, "transaction failed")
	}

	wallet := &db.Wallet{
		PublicKey: walletAddress,
		OwnerID:   owner,
	}

	if err = c.db.CreateWallet(wallet); err != nil {
		return errors.Wrap(err, "wallet insertion failed")
	}

	return nil
}

func (c creatorImp) sendTransaction(walletAddress string) error {
	chaincodeArgs := [][]byte{[]byte(walletAddress)}
	_, err := c.client.Execute(
		channel.Request{
			ChaincodeID: paymentChaincodeID,
			Fcn:         createWallet,
			Args:        chaincodeArgs,
		})
	return err
}
