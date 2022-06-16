CREATE TABLE IF NOT EXISTS commands (
  id TEXT PRIMARY KEY,
  created_at TIMESTAMP,
  updated_at TIMESTAMP,

  operation TEXT,
  status TEXT DEFAULT "new",
  metadata TEXT,

  client TEXT,
  device_id TEXT
);


INSERT INTO migrations (name, migrated_at)
  VALUES ("1655347075_2022-06-15_22-37-55_add_commands.sql", strftime('%s', "now"));