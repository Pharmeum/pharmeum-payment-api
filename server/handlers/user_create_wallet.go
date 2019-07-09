package handlers

import (
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/Pharmeum/pharmeum-payment-api/services/wallet"

	"github.com/Pharmeum/pharmeum-payment-api/utils"
)

type CreateWalletHandler struct {
	log     *logrus.Entry
	creator wallet.Creator
}

func NewCreateWalletHandler(creator wallet.Creator, log *logrus.Entry) *CreateWalletHandler {
	return &CreateWalletHandler{
		log:     log,
		creator: creator,
	}
}

//UserCreateWallet creates user wallet with elliptic curve public Key
//Two optionals types: doctor and patient
func (c CreateWalletHandler) CreateWallet(w http.ResponseWriter, r *http.Request) {
	log := c.log.WithField("handler", "user_create_wallet")

	owner := utils.UserID(r.Context())
	if owner == 0 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if err := c.creator.Create(owner); err != nil {
		log.WithError(err).Error("failed to create wallet")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
