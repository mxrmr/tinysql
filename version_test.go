// SPDX-License-Identifier: 0BSD

package tinysql

import (
	"fmt"
	"testing"
)

func TestVersion(t *testing.T) {
	res := Version()
	if len(res) == 0 {
		t.Errorf("Version() = %q, want non-empty", res)
	}
	fmt.Println("Version:", res)
}
