package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/Pharmeum/pharmeum-payment-api/services/wallet"
	"github.com/sirupsen/logrus"

	"github.com/Pharmeum/pharmeum-payment-api/utils"
)

type UserWalletsHandler struct {
	log    *logrus.Entry
	lister wallet.Lister
}

func NewUserWalletsHandler(lister wallet.Lister, log *logrus.Entry) *UserWalletsHandler {
	return &UserWalletsHandler{
		log:    log,
		lister: lister,
	}
}

//UserWallets returns list of user wallets
func (h UserWalletsHandler) UserWallets(w http.ResponseWriter, r *http.Request) {
	log := h.log.WithField("handler", "user_wallets")

	userID := utils.UserID(r.Context())
	if userID == 0 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	wallets, err := h.lister.List(userID)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}

	response, err := json.Marshal(&wallets)
	if err != nil {
		log.WithError(err).Errorf("failed to serialize user(%d) wallets", userID)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(response)
}
