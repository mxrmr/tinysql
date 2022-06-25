// SPDX-License-Identifier: 0BSD

package tinysql

// #include <sqlite3.h>
import "C"

func Version() string {
	return C.GoString(C.sqlite3_libversion())
}
