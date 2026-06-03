CREATE TABLE cars (
    vin TEXT PRIMARY KEY,
    brand TEXT NOT NULL,
    year INT NOT NULL,
    engine_id TEXT,
    transmission_id TEXT
);
CREATE TABLE engines (
    id TEXT PRIMARY KEY,
    horsepower INT NOT NULL
);
CREATE TABLE transmissions (
    id TEXT PRIMARY KEY,
    type TEXT NOT NULL
);