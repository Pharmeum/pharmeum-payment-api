package handlers

import (
	"context"
	"net/http"
	"net/url"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"

	"github.com/Pharmeum/pharmeum-payment-api/db"
	"github.com/go-chi/jwtauth"
	"github.com/sirupsen/logrus"
)

type CtxKey int

const (
	logCtxKey = iota
	httpCtxKey
	channelClientCtxKey
	dbCtxKey
	jwtCtxKey
)

func CtxLog(entry *logrus.Entry) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, logCtxKey, entry)
	}
}

func Log(r *http.Request) *logrus.Entry {
	return r.Context().Value(logCtxKey).(*logrus.Entry)
}

func CtxHTTP(http *url.URL) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, httpCtxKey, http)
	}
}

func HTTP(r *http.Request) *url.URL {
	return r.Context().Value(httpCtxKey).(*url.URL)
}

func CtxChannelClient(client *channel.Client) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, channelClientCtxKey, client)
	}
}

func ChannelClient(r *http.Request) *channel.Client {
	return r.Context().Value(channelClientCtxKey).(*channel.Client)
}

func CtxDB(db *db.DB) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, dbCtxKey, db)
	}
}

func DB(r *http.Request) *db.DB {
	return r.Context().Value(dbCtxKey).(*db.DB)
}

func CtxJWT(entry *jwtauth.JWTAuth) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, jwtCtxKey, entry)
	}
}

func JWT(r *http.Request) *jwtauth.JWTAuth {
	return r.Context().Value(jwtCtxKey).(*jwtauth.JWTAuth)
}
