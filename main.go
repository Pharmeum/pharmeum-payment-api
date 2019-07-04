package app

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"

	"github.com/Pharmeum/pharmeum-payment-api/payment"

	"github.com/pkg/errors"

	"github.com/Pharmeum/pharmeum-payment-api/config"
	"github.com/Pharmeum/pharmeum-payment-api/server"

	"github.com/sirupsen/logrus"
)

type App struct {
	config          config.Config
	log             *logrus.Entry
	paymentUploader chan payment.Uploader
}

func New(config config.Config) *App {
	return &App{
		config:          config,
		log:             config.Log(),
		paymentUploader: make(chan payment.Uploader),
	}
}

func (a *App) Start() error {
	conf := a.config

	//start Blockchain payment handler in separate go-routine
	go a.Payment()

	httpConfiguration := conf.HTTP()

	url, err := httpConfiguration.URL()
	if err != nil {
		return err
	}

	router := server.Router(
		conf.Log(),
		url,
		conf.DB(),
		conf.JWT(),
		&a.paymentUploader,
	)

	serverHost := fmt.Sprintf("%s:%s", httpConfiguration.Host, httpConfiguration.Port)
	a.log.WithField("api", "start").
		Info(fmt.Sprintf("listenig addr =  %s, tls = %v", serverHost, httpConfiguration.SSL))

	httpServer := http.Server{
		Addr:           serverHost,
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		IdleTimeout:    120 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	switch httpConfiguration.SSL {
	case true:
		tlsConfig := &tls.Config{
			CipherSuites: []uint16{
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305, // Go 1.8 only
				tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,   // Go 1.8 only
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,

				// Best disabled, as they don't provide Forward Secrecy,
				// but might be necessary for some clients
				// tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
				// tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
			},
			PreferServerCipherSuites: true,
			MinVersion:               tls.VersionTLS12,
			CurvePreferences: []tls.CurveID{
				tls.CurveP256,
				tls.X25519, // Go 1.8 only
			},
			InsecureSkipVerify: true,
		}

		httpServer.TLSConfig = tlsConfig
		if err := httpServer.ListenAndServeTLS(httpConfiguration.ServerCertPath, httpConfiguration.ServerKeyPath); err != nil {
			return errors.Wrap(err, "failed to start https server")
		}

	default:
		if err := httpServer.ListenAndServe(); err != nil {
			return errors.Wrap(err, "failed to start http server")
		}
	}

	return nil
}

func (a App) Payment() chan payment.Uploader {
	channelClient := a.config.Channel()
	log := a.config.Log().WithField("payment", "uploader")

	const (
		paymentChaincodeID = "pharmeumccpayment"
		createWallet       = "create_wallet"
	)

	for {
		uploader := <-a.paymentUploader
		switch uploader.Operation {
		case payment.CreateWallet:
			if uploader.Wallet == nil {
				log.Error("invalid wallet value, wallet can't be nil")
				continue
			}

			bytes, err := json.Marshal(uploader.Wallet)
			if err != nil {
				fmt.Println("err", err)
				log.WithError(err).Errorf("failed to serialize wallet %s", uploader.Wallet.PublicKey)
				continue
			}

			chaincodeArgs := [][]byte{[]byte(uploader.Wallet.PublicKey), bytes}

			_, err = channelClient.Execute(
				channel.Request{ChaincodeID: paymentChaincodeID, Fcn: createWallet, Args: chaincodeArgs},
				channel.WithRetry(retry.DefaultChannelOpts),
			)
			if err != nil {
				log.WithError(err).Error("failed to create wallet in Blockchain")
				continue
			}
		}
	}
}
