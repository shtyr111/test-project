CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
    id TEXT PRIMARY KEY DEFAULT uuid_generate_v4(),
    name TEXT,
    age INTEGER,
    number SERIAL
);