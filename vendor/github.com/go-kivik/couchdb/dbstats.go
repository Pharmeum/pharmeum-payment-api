package couchdb

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/go-kivik/couchdb/chttp"
	"github.com/go-kivik/kivik"
	"github.com/go-kivik/kivik/driver"
)

type dbStats struct {
	driver.DBStats
	Sizes struct {
		File     int64 `json:"file"`
		External int64 `json:"external"`
		Active   int64 `json:"active"`
	} `json:"sizes"`
	UpdateSeq json.RawMessage `json:"update_seq"` // nolint: govet
	rawBody   json.RawMessage
}

func (s *dbStats) UnmarshalJSON(p []byte) error {
	type dbStatsClone dbStats
	c := dbStatsClone(*s)
	if err := json.Unmarshal(p, &c); err != nil {
		return err
	}
	*s = dbStats(c)
	s.rawBody = p
	return nil
}

func (s *dbStats) driverStats() *driver.DBStats {
	stats := &s.DBStats
	if s.Sizes.File > 0 {
		stats.DiskSize = s.Sizes.File
	}
	if s.Sizes.External > 0 {
		stats.ExternalSize = s.Sizes.External
	}
	if s.Sizes.Active > 0 {
		stats.ActiveSize = s.Sizes.Active
	}
	stats.UpdateSeq = string(bytes.Trim(s.UpdateSeq, `"`))
	stats.RawResponse = s.rawBody
	return stats
}

func (d *db) Stats(ctx context.Context) (*driver.DBStats, error) {
	result := dbStats{}
	if _, err := d.Client.DoJSON(ctx, kivik.MethodGet, d.dbName, nil, &result); err != nil {
		return nil, err
	}
	return result.driverStats(), nil
}

type dbsInfoRequest struct {
	Keys []string `json:"keys"`
}

type dbsInfoResponse struct {
	Key    string  `json:"key"`
	DBInfo dbStats `json:"info"`
	Error  string  `json:"error"`
}

func (c *client) DBsStats(ctx context.Context, dbnames []string) ([]*driver.DBStats, error) {
	opts := &chttp.Options{
		Body: chttp.EncodeBody(dbsInfoRequest{Keys: dbnames}),
	}
	result := []dbsInfoResponse{}
	_, err := c.DoJSON(context.Background(), kivik.MethodPost, "/_dbs_info", opts, &result)
	if err != nil {
		return nil, err
	}
	stats := make([]*driver.DBStats, len(result))
	for i := range result {
		if result[i].Error == "" {
			stats[i] = result[i].DBInfo.driverStats()
		}
	}
	return stats, nil
}
