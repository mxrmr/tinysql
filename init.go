// SPDX-License-Identifier: 0BSD

package tinysql

// #cgo CFLAGS: -I/usr/local/include
// #cgo LDFLAGS: -L/usr/local/lib -lsqlite3
// #include <sqlite3.h>
import "C"
import (
	"fmt"
)

type SQLiteError struct {
	Code int
}

func (err *SQLiteError) Error() string {
	return fmt.Sprintf("SQLite error %d", err.Code)
}

type SQLiteExtendedError struct {
	Code   int
	ExCode int
	Msg    string
}

func (err *SQLiteExtendedError) Error() string {
	return fmt.Sprintf("SQLite error %d (%d): %v", err.Code, err.ExCode, err.Msg)
}

func makeError(code C.int) error {
	return &SQLiteError{int(code)}
}

func (db *DB) makeError(code C.int) error {
	exCode := int(C.sqlite3_extended_errcode(db.ptr))
	msg := C.GoString(C.sqlite3_errmsg(db.ptr))
	return &SQLiteExtendedError{int(code), exCode, msg}
}
