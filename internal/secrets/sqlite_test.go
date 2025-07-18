package secrets

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestSQLiteSecurityConfiguration(t *testing.T) {
	// Create in-memory database to test configuration
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}
	defer db.Close()

	// Apply security configuration
	err = configureSQLiteSecurity(db)
	if err != nil {
		t.Fatalf("failed to configure SQLite security: %v", err)
	}

	// Test that secure_delete is enabled
	var secureDelete string
	err = db.QueryRow("PRAGMA secure_delete").Scan(&secureDelete)
	if err != nil {
		t.Fatalf("failed to query secure_delete: %v", err)
	}
	if secureDelete != "2" { // FAST mode returns "2"
		t.Errorf("expected secure_delete to be FAST (2), got %s", secureDelete)
	}

	// Test that foreign_keys is enabled
	var foreignKeys string
	err = db.QueryRow("PRAGMA foreign_keys").Scan(&foreignKeys)
	if err != nil {
		t.Fatalf("failed to query foreign_keys: %v", err)
	}
	if foreignKeys != "1" {
		t.Errorf("expected foreign_keys to be ON (1), got %s", foreignKeys)
	}

	// Test that journal_mode is set (may be "memory" for in-memory DBs)
	var journalMode string
	err = db.QueryRow("PRAGMA journal_mode").Scan(&journalMode)
	if err != nil {
		t.Fatalf("failed to query journal_mode: %v", err)
	}
	// In-memory databases use "memory" journal mode, file databases use "wal"
	if journalMode != "wal" && journalMode != "memory" {
		t.Errorf("expected journal_mode to be WAL or memory, got %s", journalMode)
	}

	// Test that synchronous mode is NORMAL
	var synchronous string
	err = db.QueryRow("PRAGMA synchronous").Scan(&synchronous)
	if err != nil {
		t.Fatalf("failed to query synchronous: %v", err)
	}
	if synchronous != "1" { // NORMAL mode returns "1"
		t.Errorf("expected synchronous to be NORMAL (1), got %s", synchronous)
	}
}

func TestSecureDeleteIntegration(t *testing.T) {
	// Create in-memory store
	store, err := newSqliteStore(":memory:", "test-project", NewNoOpEncryptionHandler())
	if err != nil {
		t.Fatalf("failed to create secret store: %v", err)
	}
	defer store.Close()

	// Add a secret
	err = store.WriteSecret("test-key", "sensitive-data")
	if err != nil {
		t.Fatalf("failed to write secret: %v", err)
	}

	// Verify secret exists
	value, err := store.ReadSecret("test-key")
	if err != nil {
		t.Fatalf("failed to read secret: %v", err)
	}
	if value != "sensitive-data" {
		t.Errorf("expected 'sensitive-data', got %s", value)
	}

	// Delete the secret (should use secure delete)
	err = store.DeleteSecret("test-key")
	if err != nil {
		t.Fatalf("failed to delete secret: %v", err)
	}

	// Verify secret is gone
	_, err = store.ReadSecret("test-key")
	if err == nil {
		t.Error("expected error when reading deleted secret")
	}
}
