package database

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/wanan9999/s-ui/database/model"
)

func TestReadOnlyProbeDoesNotCreateDatabase(t *testing.T) {
	name := filepath.Join(t.TempDir(), "missing", "s-ui.db")
	if err := OpenReadOnlyDB(name); err == nil {
		t.Fatal("probe opened a missing database")
	}
	if _, err := os.Stat(filepath.Dir(name)); !os.IsNotExist(err) {
		t.Fatalf("probe created a directory: %v", err)
	}
}

func TestReadOnlyProbeDuringWrite(t *testing.T) {
	name := filepath.Join(t.TempDir(), "s-ui.db")
	if err := InitDB(name); err != nil {
		t.Fatal(err)
	}
	writer := GetDB()
	pool, _ := writer.DB()
	defer pool.Close()
	if err := writer.Create(&model.Setting{Key: "webPort", Value: "2095"}).Error; err != nil {
		t.Fatal(err)
	}
	tx := writer.Begin()
	if tx.Error != nil {
		t.Fatal(tx.Error)
	}
	defer tx.Rollback()
	if err := tx.Model(&model.Setting{}).Where("key = ?", "webPort").Update("value", "9090").Error; err != nil {
		t.Fatal(err)
	}
	if err := OpenReadOnlyDB(name); err != nil {
		t.Fatal(err)
	}
	reader := GetDB()
	readPool, _ := reader.DB()
	defer readPool.Close()
	var setting model.Setting
	if err := reader.Where("key = ?", "webPort").First(&setting).Error; err != nil {
		t.Fatal(err)
	}
	if setting.Value != "2095" {
		t.Fatalf("probe saw uncommitted setting: %q", setting.Value)
	}
	if err := reader.Create(&model.Setting{Key: "probe", Value: "must not write"}).Error; err == nil {
		t.Fatal("probe connection allowed writes")
	}
	if err := tx.Commit().Error; err != nil {
		t.Fatal(err)
	}
	if err := reader.Where("key = ?", "webPort").First(&setting).Error; err != nil || setting.Value != "9090" {
		t.Fatalf("probe did not see committed settings: %q, %v", setting.Value, err)
	}
}
