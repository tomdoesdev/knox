package database

const tablesSchema = `
 CREATE TABLE linked_vaults (
      id INTEGER PRIMARY KEY,
      alias TEXT NOT NULL UNIQUE,
      path TEXT NOT NULL UNIQUE,
      created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
  );

  CREATE TABLE linked_projects (
      id INTEGER PRIMARY KEY,
      name TEXT NOT NULL UNIQUE,
      vault_id INTEGER NOT NULL,
      project_name TEXT, -- actual name in vault if different
      created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
      FOREIGN KEY (vault_id) REFERENCES linked_vaults(id) ON DELETE CASCADE
  );

  CREATE TABLE workspace_settings (
      key TEXT PRIMARY KEY,
      value TEXT NOT NULL,
      category TEXT NOT NULL, -- 'meta' or 'config'
      updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
  );
`
