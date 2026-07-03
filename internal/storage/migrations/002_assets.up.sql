CREATE TABLE IF NOT EXISTS assets (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE,
    site TEXT NOT NULL DEFAULT 'factory-1',
    area TEXT NOT NULL DEFAULT '',
    line TEXT NOT NULL DEFAULT '',
    type TEXT NOT NULL DEFAULT 'machine',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_assets_site ON assets(site);
CREATE INDEX IF NOT EXISTS idx_assets_name ON assets(name);

CREATE TABLE IF NOT EXISTS asset_sensors (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    asset_id INTEGER NOT NULL,
    sensor_id TEXT NOT NULL,
    sensor_type TEXT NOT NULL DEFAULT 'generic',
    topic TEXT NOT NULL DEFAULT '',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (asset_id) REFERENCES assets(id) ON DELETE CASCADE,
    UNIQUE(asset_id, sensor_id)
);

CREATE INDEX IF NOT EXISTS idx_asset_sensors_asset_id ON asset_sensors(asset_id);
CREATE INDEX IF NOT EXISTS idx_asset_sensors_sensor_id ON asset_sensors(sensor_id);
