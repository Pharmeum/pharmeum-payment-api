package couchdb

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/go-kivik/couchdb/chttp"
	"github.com/go-kivik/kivik"
	"github.com/go-kivik/kivik/driver"
)

func (c *client) AllDBs(ctx context.Context, opts map[string]interface{}) ([]string, error) {
	query, err := optionsToParams(opts)
	if err != nil {
		return nil, err
	}
	var allDBs []string
	_, err = c.DoJSON(ctx, kivik.MethodGet, "/_all_dbs", &chttp.Options{Query: query}, &allDBs)
	return allDBs, err
}

func (c *client) DBExists(ctx context.Context, dbName string, _ map[string]interface{}) (bool, error) {
	if dbName == "" {
		return false, missingArg("dbName")
	}
	_, err := c.DoError(ctx, kivik.MethodHead, dbName, nil)
	if kivik.StatusCode(err) == kivik.StatusNotFound {
		return false, nil
	}
	return err == nil, err
}

func (c *client) CreateDB(ctx context.Context, dbName string, opts map[string]interface{}) error {
	if dbName == "" {
		return missingArg("dbName")
	}
	query, err := optionsToParams(opts)
	if err != nil {
		return err
	}
	_, err = c.DoError(ctx, kivik.MethodPut, dbName, &chttp.Options{Query: query})
	return err
}

func (c *client) DestroyDB(ctx context.Context, dbName string, _ map[string]interface{}) error {
	if dbName == "" {
		return missingArg("dbName")
	}
	_, err := c.DoError(ctx, kivik.MethodDelete, dbName, nil)
	return err
}

func (c *client) DBUpdates(ctx context.Context) (updates driver.DBUpdates, err error) {
	resp, err := c.DoReq(ctx, kivik.MethodGet, "/_db_updates?feed=continuous&since=now", nil)
	if err != nil {
		return nil, err
	}
	if err := chttp.ResponseError(resp); err != nil {
		return nil, err
	}
	return newUpdates(ctx, resp.Body), nil
}

type couchUpdates struct {
	*iter
}

var _ driver.DBUpdates = &couchUpdates{}

type updatesParser struct{}

var _ parser = &updatesParser{}

func (p *updatesParser) decodeItem(i interface{}, dec *json.Decoder) error {
	return dec.Decode(i)
}

func newUpdates(ctx context.Context, body io.ReadCloser) *couchUpdates {
	return &couchUpdates{
		iter: newIter(ctx, nil, "", body, &updatesParser{}),
	}
}

func (u *couchUpdates) Next(update *driver.DBUpdate) error {
	return u.iter.next(update)
}

// Ping queries the /_up endpoint, and returns true if there are no errors, or
// if a 400 (Bad Request) is returned, and the Server: header indicates a server
// version prior to 2.x.
func (c *client) Ping(ctx context.Context) (bool, error) {
	resp, err := c.DoError(ctx, kivik.MethodHead, "/_up", nil)
	if kivik.StatusCode(err) == kivik.StatusBadRequest {
		return strings.HasPrefix(resp.Header.Get("Server"), "CouchDB/1."), nil
	}
	if kivik.StatusCode(err) == http.StatusNotFound {
		return false, nil
	}
	return err == nil, err
}
