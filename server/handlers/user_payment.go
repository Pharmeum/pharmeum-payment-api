package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/Pharmeum/pharmeum-payment-api/services/wallet"

	"github.com/pkg/errors"

	"github.com/Pharmeum/pharmeum-payment-api/utils"

	"github.com/go-ozzo/ozzo-validation"

	"github.com/Pharmeum/pharmeum-payment-api/db"
	"github.com/sirupsen/logrus"
)

type paymentRequest struct {
	Receiver string `json:"receiver"`
	Sender   string `json:"sender"`
	Amount   string `json:"amount"`
}

func (p paymentRequest) Validate() error {
	return validation.Errors{
		"receiver": validation.Validate(&p.Receiver, validation.Required),
		"sender":   validation.Validate(&p.Sender, validation.Required),
	}.Filter()
}

type PaymentHandler struct {
	Log   *logrus.Entry
	DB    *db.DB
	Payer wallet.Payer
}

func (p PaymentHandler) isAllowed(address string, ownerID uint64) (bool, error) {
	err := p.DB.IsAllowed(address, ownerID)
	if err == nil {
		return true, nil
	}

	if err == sql.ErrNoRows {
		return false, nil
	}

	return false, err
}

func (p *PaymentHandler) Handle(w http.ResponseWriter, r *http.Request) {
	log := p.Log.WithField("handler", "user_payment")

	request := &paymentRequest{}
	if err := json.NewDecoder(r.Body).Decode(request); err != nil {
		log.WithError(err).Debug("failed to decode payment request")
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write(ErrResponse(http.StatusBadRequest, err))
		return
	}

	if err := request.Validate(); err != nil {
		log.WithError(err).Debug("payment request validation failed")
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write(ErrResponse(http.StatusBadRequest, err))
		return
	}

	senderID := utils.UserID(r.Context())
	if senderID == 0 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	allowed, err := p.isAllowed(request.Sender, senderID)
	if err != nil {
		log.WithError(err).Error("failed to check if user allowed to see content")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !allowed {
		log.WithFields(map[string]interface{}{
			"sender_id":      senderID,
			"sender_address": request.Sender,
		}).Debug("user is not allowed to see content")
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write(ErrResponse(http.StatusBadRequest, errors.New("user is not allowed to use payment operation from this wallet")))
		return
	}

	if err := p.Payer.Pay(
		request.Sender,
		request.Receiver,
		request.Amount); err != nil {
		switch errors.Cause(err) {
		case wallet.ErrInvalidAmountOfTokens:
			log.WithError(err).Debug("failed to send payment request")
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write(ErrResponse(http.StatusBadRequest, err))
			return
		case sql.ErrNoRows:
			log.WithError(err).
				WithField("receiver_wallet", request.Receiver).
				Debug("receiver wallet does not exist")
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write(ErrResponse(http.StatusBadRequest, err))
			return
		default:
			log.WithError(err).Error("failed to send payment transaction")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}
