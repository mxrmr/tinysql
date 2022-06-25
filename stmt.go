// SPDX-License-Identifier: 0BSD

package tinysql

// #include <sqlite3.h>
// #include <stdlib.h>
// static const char *const_char_cast(const unsigned char *ptr) { return (const char *)ptr; }
import "C"
import (
	"unsafe"
)

type Stmt struct {
	ptr *C.sqlite3_stmt
}

func prepareStmt(db *DB, sql string) (*Stmt, error) {
	csql := C.CString(sql)
	defer C.free(unsafe.Pointer(csql))

	var stmt *C.sqlite3_stmt
	res := C.sqlite3_prepare_v2(db.ptr, csql, -1, &stmt, nil)
	if res != C.SQLITE_OK {
		return nil, makeError(res)
	}
	return &Stmt{stmt}, nil
}

func (s *Stmt) Close() {
	C.sqlite3_finalize(s.ptr)
}

func (s *Stmt) Exec(args ...any) error {
	if err := s.bindArgs(args...); err != nil {
		return err
	}
	if res := C.sqlite3_step(s.ptr); res != C.SQLITE_DONE {
		return makeError(res)
	}
	return nil
}

func (s *Stmt) Query(args ...any) (*Rows, error) {
	if err := s.bindArgs(args...); err != nil {
		return nil, err
	}
	return &Rows{s.ptr, nil}, nil
}

func (s *Stmt) bindArgs(args ...any) error {
	C.sqlite3_reset(s.ptr)
	for idx, arg := range args {
		var err C.int
		argNum := C.int(idx+1)
		switch typedArg := arg.(type) {
		case int:
			err = C.sqlite3_bind_int(s.ptr, argNum, C.int(typedArg))
		case string:
			cVal := C.CString(typedArg)
			err = C.sqlite3_bind_text(
				s.ptr,
				argNum,
				cVal,
				C.int(len(typedArg)),
				(*[0]byte)(C.free),
			)
		case []byte:
			cVal := C.CBytes(typedArg)
			err = C.sqlite3_bind_blob(
				s.ptr,
				argNum,
				cVal,
				C.int(len(typedArg)),
				(*[0]byte)(C.free),
			)
		default:
			panic("unsupported type")
		}
		if err != C.SQLITE_OK {
			return makeError(err)
		}
	}
	return nil
}

type Rows struct {
	ptr *C.sqlite3_stmt
	err error
}

func (rs *Rows) Next() bool {
	switch res := C.sqlite3_step(rs.ptr); res {
	case C.SQLITE_ROW:
		return true
	case C.SQLITE_DONE:
		return false
	default:
		rs.err = makeError(res)
		return false
	}
}

func (rs *Rows) Err() error {
	return rs.err
}

func (rs *Rows) Scan(dest ...any) error {
	for idx, dst := range dest {
		argIdx := C.int(idx)
		switch typedDst := dst.(type) {
		case *string:
			cVal := C.sqlite3_column_text(rs.ptr, argIdx)
			if cVal == nil {
				*typedDst = ""
			} else {
				*typedDst = C.GoString(C.const_char_cast(cVal))
			}
		case *int:
			cVal := C.sqlite3_column_int64(rs.ptr, argIdx)
			*typedDst = int(cVal)
		case *[]byte:
			cVal := C.sqlite3_column_blob(rs.ptr, argIdx)
			cLen := C.sqlite3_column_bytes(rs.ptr, argIdx)
			*typedDst = C.GoBytes(cVal, cLen)
		default:
			panic("unsupported type")
		}
	}
	return nil
}
