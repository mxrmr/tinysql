// SPDX-License-Identifier: 0BSD

package tinysql

import "testing"

func testDB() *DB {
	res, _ := Open(":memory:")
	return res
}

func TestOpen(t *testing.T) {
	db, err := Open(":memory:")
	if err != nil {
		t.Fatalf("Open() got err %v", err)
	}
	err = db.Close()
	if err != nil {
		t.Fatalf("Close() got err %v", err)
	}
}

func TestDBExec(t *testing.T) {
	db := testDB()
	defer db.Close()

	db.Exec(`CREATE TABLE t (value TEXT)`)
	db.Exec(`INSERT INTO t (value) VALUES ('Hello World')`)

	stmt, _ := db.Prepare("SELECT value FROM t")
	defer stmt.Close()

	rows, _ := stmt.Query()
	rows.Next()

	var got string
	rows.Scan(&got)

	if want := "Hello World"; got != want {
		t.Errorf("Exec() got %q, want %q", got, want)
	}
}
