package wallet

import (
	"database/sql"

	"github.com/Pharmeum/pharmeum-payment-api/db"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

type Payer interface {
	Pay(sender, receiver, amount string) error
}

type PayerImp struct {
	client *channel.Client
	db     *db.DB
}

func NewPayer(client *channel.Client, db *db.DB) Payer {
	return &PayerImp{
		client: client,
		db:     db,
	}
}

var (
	ErrInvalidAmountOfTokens = errors.New("invalid amount of tokens")
)

func (p PayerImp) ValidateAmount(amount string) error {
	decamount, err := decimal.NewFromString(amount)
	if err != nil {
		return err
	}

	if decamount.IsNegative() || decamount.IsZero() {
		return ErrInvalidAmountOfTokens
	}

	return nil
}

func (p PayerImp) Pay(sender, receiver, amount string) error {
	//check amount(should not be negative or zero)
	if err := p.ValidateAmount(amount); err != nil {
		return errors.Wrap(err, "amount validation failed")
	}

	//check if receiver wallet exists
	if _, err := p.db.Wallet(receiver); err != nil {
		if err == sql.ErrNoRows {
			return errors.Wrap(err, "receiver wallet does not exists")
		}
		return errors.Wrap(err, "failed to get receiver wallet")
	}

	//send payment transaction
	chaincodeArgs := [][]byte{[]byte(sender), []byte(receiver), []byte(amount)}
	_, err := p.client.Execute(
		channel.Request{
			ChaincodeID: paymentChaincodeID,
			Fcn:         transferPayment,
			Args:        chaincodeArgs,
		})
	return err
}
