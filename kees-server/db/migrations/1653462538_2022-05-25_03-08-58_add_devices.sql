CREATE TABLE IF NOT EXISTS devices (
  id TEXT PRIMARY KEY,
  created_at TIMESTAMP,
  updated_at TIMESTAMP,
  name TEXT,
  version TEXT,
  controller TEXT,
  online INTEGER,
  last_heartbeat TIMESTAMP,
  token TEXT,
  capabilities TEXT
);

INSERT INTO migrations (name, migrated_at)
  VALUES ("1653462538_2022-05-25_03-08-58_add_devices.sql", strftime('%s', "now"));