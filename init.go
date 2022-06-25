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

func makeError(code C.int) error {
	return &SQLiteError{int(code)}
}
