package tinysql

import (
	"strings"
	"testing"
)

func TestMakeError(t *testing.T) {
	db := testDB()
	defer db.Close()

	err := db.Exec(`SELECT not_a_column FROM not_a_table`)
	if got, want := err.Error(), "no such table"; !strings.Contains(got, want) {
		t.Errorf("Exec() got %q, want %q", got, want)
	}
}
