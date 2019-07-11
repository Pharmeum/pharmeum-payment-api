package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"

	"github.com/Pharmeum/pharmeum-payment-api/db"

	"github.com/Pharmeum/pharmeum-payment-api/utils"

	"github.com/Pharmeum/pharmeum-payment-api/services/wallet"
	"github.com/go-ozzo/ozzo-validation"
	"github.com/sirupsen/logrus"
)

type WalletBalanceRequest struct {
	Address string `json:"address"`
}

func (w WalletBalanceRequest) Validate() error {
	return validation.Errors{
		"address": validation.Validate(&w.Address, validation.Required),
	}.Filter()
}

type WalletBalanceResponse struct {
	Balance string `json:"balance"`
}

type WalletBalanceHandler struct {
	Log      *logrus.Entry
	DB       *db.DB
	Balancer wallet.Balancer
}

func (h WalletBalanceHandler) isAllowed(address string, ownerID uint64) (bool, error) {
	err := h.DB.IsAllowed(address, ownerID)
	if err == nil {
		return true, nil
	}

	if err == sql.ErrNoRows {
		return false, nil
	}

	return false, err
}

func (h WalletBalanceHandler) WalletBalance(w http.ResponseWriter, r *http.Request) {
	log := h.Log.WithField("handler", "wallet_balance")

	request := &WalletBalanceRequest{}
	if err := json.NewDecoder(r.Body).Decode(request); err != nil {
		log.WithError(err).Debug("failed to decode request")
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write(ErrResponse(http.StatusBadRequest, err))
		return
	}

	if err := request.Validate(); err != nil {
		log.WithError(err).Debug("invalid request arguments")
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write(ErrResponse(http.StatusBadRequest, err))
		return
	}

	ownerID := utils.UserID(r.Context())
	if ownerID == 0 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	allowed, err := h.isAllowed(request.Address, ownerID)
	if err != nil {
		log.WithError(err).Error("failed to check if user allowed to see content")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !allowed {
		log.WithFields(map[string]interface{}{
			"owner_id": ownerID,
			"address":  request.Address,
		}).Debug("user is not allowed to see content")
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write(ErrResponse(http.StatusBadRequest, errors.New("You are not allowed to see balance of requested wallet")))
		return
	}

	result, err := h.Balancer.Balance(request.Address)
	if err != nil {
		if err.Error() == "Not Found: missing" {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write(ErrResponse(http.StatusBadRequest, errors.New("wallet not exist")))
			return
		}

		log.WithFields(map[string]interface{}{
			"address": request.Address,
			"user_id": ownerID,
		}).WithError(err).Error("failed to get user balance")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(&WalletBalanceResponse{
		Balance: result,
	})
	if err != nil {
		log.WithError(err).Error("failed to serialize response")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(response)
}
