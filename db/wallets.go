package db

import "github.com/go-ozzo/ozzo-dbx"

const walletsTableName = "wallets"

type Wallet struct {
	PublicKey string `db:"public_key" json:"public_key"`
	Name      string `db:"name" json:"name"`
	Kind      string `db:"kind" json:"kind"`
	OwnerID   uint64 `db:"owner_id" json:"owner_id"`
}

func (w Wallet) TableName() string {
	return walletsTableName
}

func (d *DB) UserWallets(ownerID uint64) ([]Wallet, error) {
	var wallets []Wallet
	err := d.db.Select().Where(dbx.HashExp{"owner_id": ownerID}).All(&wallets)
	return wallets, err
}

func (d *DB) CreateWallet(wallet *Wallet) error {
	return d.db.Model(wallet).Insert()
}
