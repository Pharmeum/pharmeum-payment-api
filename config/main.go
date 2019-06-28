package config

import (
	"sync"

	"github.com/Pharmeum/pharmeum-payment-api/db"

	"github.com/go-chi/jwtauth"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"

	"github.com/sirupsen/logrus"
)

type Config interface {
	HTTP() *HTTP
	Log() *logrus.Entry
	DB() *db.DB
	JWT() *jwtauth.JWTAuth
}

type ConfigImpl struct {
	sync.Mutex

	//internal objects
	http          *HTTP
	log           *logrus.Entry
	channelClient *channel.Client
	db            *db.DB
	jwt           *jwtauth.JWTAuth
}

func New() Config {
	return &ConfigImpl{
		Mutex: sync.Mutex{},
	}
}
