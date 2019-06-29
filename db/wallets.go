package db

const walletsTableName = "wallets"

type Wallet struct {
	ID   uint64 `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
	Kind string `db:"kind" json:"kind"`
}

func (w Wallet) TableName() string {
	return walletsTableName
}

func (d *DB) UserWallets(id uint64) ([]Wallet, error) {
	return nil, nil
}
