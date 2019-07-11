package config

import (
	"fmt"

	"github.com/caarlos0/env"

	_ "github.com/go-kivik/couchdb" // The CouchDB driver
	"github.com/go-kivik/kivik"
)

type CouchDB struct {
	Host     string `env:"PHARMEUM_COUCH_HOST,required"`
	Port     uint32 `env:"PHARMEUM_COUCH_PORT,required"`
	Username string `env:"PHARMEUM_COUCH_USERNAME,required"`
	Password string `env:"PHARMEUM_COUCH_PASSWORD,required"`
	TLS      bool   `env:"PHARMEUM_COUCH_TLS"`
}

func (c CouchDB) URL() (couchDBURL string) {
	switch c.TLS {
	case true:
		couchDBURL = "https://"
	case false:
		couchDBURL = "http://"
	}

	couchDBURL = fmt.Sprintf("%s%s:%s@%s:%d", couchDBURL, c.Username, c.Password, c.Host, c.Port)
	return
}

func (c *ConfigImpl) CouchClient() *kivik.Client {
	if c.couchDBClient != nil {
		return c.couchDBClient
	}

	c.Lock()
	defer c.Unlock()

	couchDB := &CouchDB{}
	if err := env.Parse(couchDB); err != nil {
		panic(err)
	}

	client, err := kivik.New("couch", couchDB.URL())
	if err != nil {
		panic(err)
	}
	c.couchDBClient = client

	return c.couchDBClient
}
