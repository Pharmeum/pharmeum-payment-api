package db

import "github.com/go-ozzo/ozzo-dbx"

const walletsTableName = "wallets"

type Wallet struct {
	PublicKey string `db:"public_key" json:"public_key"`
	OwnerID   uint64 `db:"owner_id" json:"-"`
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

func (d *DB) IsAllowed(address string, ownerID uint64) (err error) {
	wallet := &Wallet{}
	err = d.db.Select().
		Where(dbx.HashExp{"public_key": address, "owner_id": ownerID}).
		One(wallet)
	return
}

func (d *DB) Wallet(address string) (*Wallet, error) {
	wallet := &Wallet{}
	err := d.db.Select().
		Where(dbx.HashExp{"public_key": address}).
		One(wallet)
	return wallet, err
}
