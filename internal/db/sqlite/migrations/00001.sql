CREATE TABLE IF NOT EXISTS reports (
   id INTEGER PRIMARY KEY AUTOINCREMENT,
   date TEXT,
   severity_count TEXT
);

CREATE TABLE schema_version (
    version INTEGER NOT NULL
);

CREATE TABLE alerts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    fingerprint TEXT UNIQUE NOT NULL,
    started INTEGER NOT NULL,
    ended INTEGER,
    severity INTEGER NOT NULL,
    instance TEXT NOT NULL,
    name TEXT NOT NULL,
    data JSONB NOT NULL,
    CHECK (ended >= started)
);

-- Indexes for optimization
CREATE INDEX idx_alerts_started ON alerts(started);
CREATE INDEX idx_alerts_ended ON alerts(ended);
CREATE INDEX idx_alerts_severity ON alerts(severity);
CREATE INDEX idx_alerts_instance ON alerts(instance);
CREATE INDEX idx_alerts_name ON alerts(name);
