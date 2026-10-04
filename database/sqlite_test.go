package database

import (
	"path/filepath"
	"testing"
)

func TestPureGoSQLitePragmasAndJSON(t *testing.T) {
	if err := OpenDB(filepath.Join(t.TempDir(), "pure-go.db")); err != nil {
		t.Fatal(err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatal(err)
	}
	defer sqlDB.Close()
	var mode string
	if err := db.Raw("PRAGMA journal_mode").Scan(&mode).Error; err != nil || mode != "wal" {
		t.Fatalf("journal mode=%q: %v", mode, err)
	}
	for query, want := range map[string]int{"PRAGMA busy_timeout": 10000, "PRAGMA cache_size": -200, "SELECT json_extract('{\"value\":42}', '$.value')": 42} {
		var got int
		if err := db.Raw(query).Scan(&got).Error; err != nil || got != want {
			t.Fatalf("%s: got %d, want %d: %v", query, got, want, err)
		}
	}
}
