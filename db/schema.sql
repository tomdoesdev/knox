CREATE TABLE vault_info (
    version INTEGER
);

CREATE TABLE secrets (
    id INTEGER PRIMARY KEY,
    namespace TEXT NOT NULL DEFAULT 'global',
    key TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT ''
)