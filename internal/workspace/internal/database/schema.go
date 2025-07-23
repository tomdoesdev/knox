package database

const tablesSchema = `
 CREATE TABLE linked_vaults (
      id INTEGER PRIMARY KEY,
      alias TEXT NOT NULL UNIQUE,
      path TEXT NOT NULL UNIQUE,
      created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
  );

  CREATE TABLE workspace_settings (
      key TEXT PRIMARY KEY,
      value TEXT NOT NULL,
      category TEXT NOT NULL, -- 'meta' or 'config'
      updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
  );
`
