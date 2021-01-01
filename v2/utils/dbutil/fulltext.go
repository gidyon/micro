package dbutil

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
)

// FullTextIndex is full text search index
const FullTextIndex = "fts_search_index"

// CreateFullTextIndex creates a full-text index
func CreateFullTextIndex(db *gorm.DB, tableName string, columns ...string) error {
	if ok := db.Migrator().HasIndex(tableName, FullTextIndex); ok {
		return nil
	}
	sqlQuery := fmt.Sprintf(
		"CREATE FULLTEXT INDEX %s ON %s(%s)",
		FullTextIndex, tableName, strings.Join(columns, ","),
	)
	err := db.Table(tableName).Exec(sqlQuery).Error
	return err
}

// DropFullTextIndex drops a full-text index
func DropFullTextIndex(db *gorm.DB, tableName string) error {
	if ok := db.Migrator().HasIndex(tableName, FullTextIndex); !ok {
		return nil
	}
	sqlQuery := fmt.Sprintf("ALTER TABLE %s	DROP INDEX %s", tableName, FullTextIndex)
	err := db.Table(tableName).Exec(sqlQuery).Error
	return err
}
