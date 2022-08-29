// SPDX-License-Identifier: 0BSD

package tinysql

// #include <sqlite3.h>
// #include <stdlib.h>
import "C"
import (
	"unsafe"
)

type DB struct {
	ptr *C.sqlite3
}

func Open(path string) (*DB, error) {
	cPath := C.CString(path)
	defer C.free(unsafe.Pointer(cPath))

	var db *C.sqlite3

	res := C.sqlite3_open_v2(
		cPath,
		&db,
		C.SQLITE_OPEN_READWRITE|C.SQLITE_OPEN_CREATE|C.SQLITE_OPEN_NOMUTEX,
		nil,
	)
	if res != C.SQLITE_OK {
		return nil, makeError(res)
	}

	return &DB{db}, nil
}

func (db *DB) Close() error {
	res := C.sqlite3_close_v2(db.ptr)
	if res != C.SQLITE_OK {
		return makeError(res)
	}
	return nil
}

func (db *DB) Exec(sql string) error {
	cSQL := C.CString(sql)
	defer C.free(unsafe.Pointer(cSQL))

	res := C.sqlite3_exec(db.ptr, cSQL, nil, nil, nil)
	if res != C.SQLITE_OK {
		return db.makeError(res)
	}
	return nil
}

func (db *DB) Begin() (*Tx, error) {
	return beginTx(db)
}

func (db *DB) Prepare(sql string) (*Stmt, error) {
	return prepareStmt(db, sql)
}
