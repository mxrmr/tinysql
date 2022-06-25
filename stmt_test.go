// SPDX-License-Identifier: 0BSD

package tinysql

import (
	"bytes"
	"testing"
)

func TestStmtExec(t *testing.T) {
	db := testDB()
	defer db.Close()

	db.Exec(`CREATE TABLE t (id INTEGER, key TEXT, value BLOB)`)

	sampleData := []struct {
		id    int
		key   string
		value []byte
	}{
		{42, "ABC", []byte{1, 2, 3}},
		{43, "XYZ", []byte{4, 2}},
		{44, "MNO", []byte{}},
	}

	func() {
		stmt, _ := db.Prepare(`INSERT INTO t (id, key, value) VALUES (?, ?, ?)`)
		defer stmt.Close()

		for _, d := range sampleData {
			if err := stmt.Exec(d.id, d.key, d.value); err != nil {
				t.Fatalf("Exec(%q) got error %v", d, err)
			}
		}
	}()

	func() {
		stmt, _ := db.Prepare(`SELECT id, key, value FROM t WHERE id < ?`)
		defer stmt.Close()

		rows, _ := stmt.Query(sampleData[2].id)

		for i := 0; i < 2; i++ {
			if !rows.Next() {
				t.Fatalf("Next() #%v got false, want true", i)
			}
			var id int
			var key string
			var value []byte
			if err := rows.Scan(&id, &key, &value); err != nil {
				t.Fatalf("Scan() #%v got err %v", i, err)
			}
			if got, want := id, sampleData[i].id; got != want {
				t.Errorf("id got %v, want %v", got, want)
			}
			if got, want := key, sampleData[i].key; got != want {
				t.Errorf("key got %v, want %v", got, want)
			}
			if got, want := value, sampleData[i].value; !bytes.Equal(got, want) {
				t.Errorf("value got %v, want %v", got, want)
			}
		}
		if rows.Next() {
			t.Errorf("Scan() has too many results")
		}
	}()
}

func TestQueryBytes(t *testing.T) {
	testQuery(t, bytes.Equal, []struct {
		sql string
		res []byte
	}{
		{`SELECT X''`, []byte{}},
		{`SELECT X'00'`, []byte{0}},
		{`SELECT X'0001'`, []byte{0, 1}},
		{`SELECT NULL`, []byte{}},
	})
}

func TestQueryString(t *testing.T) {
	testQuery(t, func(a, b string) bool { return a == b }, []struct {
		sql string
		res string
	}{
		{`SELECT ''`, ""},
		{`SELECT 'ABC'`, "ABC"},
		{`SELECT NULL`, ""},
	})
}

func TestQueryInt(t *testing.T) {
	testQuery(t, func(a, b int) bool { return a == b }, []struct {
		sql string
		res int
	}{
		{`SELECT 0`, 0},
		{`SELECT 1`, 1},
		{`SELECT 123`, 123},
		{`SELECT NULL`, 0},
	})
}

func testQuery[T any](t *testing.T, cmp func(a, b T) bool, testCases []struct {
	sql string
	res T
}) {
	db := testDB()
	defer db.Close()

	for _, tc := range testCases {
		stmt, _ := db.Prepare(tc.sql)
		defer stmt.Close()

		rows, _ := stmt.Query()
		if !rows.Next() {
			t.Fatalf("Query(%q).Next() got false, want true", tc.sql)
		}

		var got T
		if err := rows.Scan(&got); err != nil {
			t.Fatalf("Query(%q) got error %v", tc.sql, err)
		}

		if want := tc.res; !cmp(got, want) {
			t.Errorf("Query(%q) got %v, want %v", tc.sql, got, want)
		}
	}
}
