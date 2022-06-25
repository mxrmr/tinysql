// SPDX-License-Identifier: 0BSD

package tinysql

type Tx struct {
	db *DB
}

func beginTx(db *DB) (*Tx, error) {
	err := db.Exec(`BEGIN EXCLUSIVE TRANSACTION`)
	if err != nil {
		return nil, err
	}
	return &Tx{db}, nil
}

func (tx *Tx) Commit() error {
	return tx.db.Exec(`COMMIT TRANSACTION`)
}

func (tx *Tx) Rollback() error {
	return tx.db.Exec(`ROLLBACK TRANSACTION`)
}
