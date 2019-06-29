package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Pharmeum/pharmeum-payment-api/utils"

	"github.com/Pharmeum/pharmeum-payment-api/db"

	"github.com/stellar/go/keypair"

	validation "github.com/go-ozzo/ozzo-validation"
)

type UserWalletsRequest struct {
	Name string `json:"name"`
	Kind string `json:"kind"`
}

const (
	patientWalletKind = "patient"
	doctorWalletKind  = "doctor"
)

func (u UserWalletsRequest) Validate() error {
	var err error
	switch u.Kind {
	case patientWalletKind, doctorWalletKind:
		//everything is ok
	default:
		err = ErrInvalidWalletKind
	}

	return validation.Errors{
		"name": validation.Validate(u.Name, validation.Required),
		"kind": err,
	}.Filter()
}

//UserCreateWallet creates user wallet with Stellar public Key
//Two optionals types: doctor and patient

func UserCreateWallet(w http.ResponseWriter, r *http.Request) {
	log := Log(r).WithField("handler", "user_create_wallet")

	var userWalletRequest UserWalletsRequest
	if err := json.NewDecoder(r.Body).Decode(&userWalletRequest); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(ErrResponse(http.StatusBadRequest, err))
		return
	}

	if err := userWalletRequest.Validate(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(ErrResponse(http.StatusBadRequest, err))
		return
	}

	owner := utils.UserID(r.Context())
	if owner == 0 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	//generate unique public key
	kp, err := keypair.Random()
	if err != nil {
		log.WithError(err).Error("failed to generate random keypair")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	wallet := &db.Wallet{
		PublicKey: kp.Address(),
		Name:      userWalletRequest.Name,
		Kind:      userWalletRequest.Kind,
		OwnerID:   owner,
	}

	//TODO: missed GRPC connection to Hyperledger Fabric
	//TODO: send Wallet request to separate go-routine to speedup time to respond
	if err := DB(r).CreateWallet(wallet); err != nil {
		log.WithError(err).Error("failed to create wallet")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
