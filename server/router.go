package server

import (
	"net/url"
	"time"

	"github.com/Pharmeum/pharmeum-payment-api/services/wallet"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"

	"github.com/go-chi/jwtauth"

	"github.com/Pharmeum/pharmeum-payment-api/db"
	"github.com/Pharmeum/pharmeum-payment-api/server/handlers"
	"github.com/Pharmeum/pharmeum-payment-api/server/middlewares"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/sirupsen/logrus"
)

const durationThreshold = time.Second * 10

func Router(
	log *logrus.Entry,
	http *url.URL,
	db *db.DB,
	jwtAuth *jwtauth.JWTAuth,
	client *channel.Client,
) chi.Router {
	router := chi.NewRouter()

	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*", "https://localhost:3000"},
		AllowedMethods:   []string{"*", "GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*", "Accept", "Authorization", "Content-Type", "X-CSRF-Token", "x-auth"},
		ExposedHeaders:   []string{"*", "Link"},
		AllowCredentials: true,
		MaxAge:           300,
	})

	router.Use(
		cors.Handler,
		middlewares.Logger(log, durationThreshold),
		middlewares.Ctx(
			handlers.CtxLog(log),
			handlers.CtxHTTP(http),
			handlers.CtxDB(db),
			handlers.CtxJWT(jwtAuth),
		),
	)

	router.Route("/user", func(router chi.Router) {
		router.Use(
			jwtauth.Verifier(jwtAuth),
			jwtauth.Authenticator,
		)
		{
			handler := handlers.NewUserWalletsHandler(wallet.NewLister(client, db), log)
			router.Get("/wallets", handler.UserWallets)
		}
		{
			createWalletHandler := handlers.NewCreateWalletHandler(
				wallet.NewCreator(client, db),
				log,
			)
			router.Post("/create_wallet", createWalletHandler.CreateWallet)
		}
	})

	return router
}
