package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/Pharmeum/pharmeum-payment-api/utils"
)

//UserWallets returns list of user wallets
func UserWallets(w http.ResponseWriter, r *http.Request) {
	log := Log(r).WithField("handler", "user_wallets")

	userID := utils.UserID(r.Context())
	if userID == 0 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	wallets, err := DB(r).UserWallets(userID)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}

	if len(wallets) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	response, err := json.Marshal(&wallets)
	if err != nil {
		log.WithError(err).Errorf("failed to serialize user(%d) wallets", userID)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(response)
}
