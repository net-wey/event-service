CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS venues (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name        VARCHAR(255) NOT NULL,
    address     VARCHAR(500) NOT NULL DEFAULT '',
    capacity    INT          NOT NULL DEFAULT 0,
    created_at  TIMESTAMP    NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS events (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title       VARCHAR(255) NOT NULL,
    description TEXT         NOT NULL DEFAULT '',
    location    VARCHAR(255) NOT NULL DEFAULT '',
    event_date  TIMESTAMP    NOT NULL,
    created_at  TIMESTAMP    NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMP    NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS participants (
    id            UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    event_id      UUID         NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    name          VARCHAR(255) NOT NULL,
    email         VARCHAR(255) NOT NULL,
    registered_at TIMESTAMP    NOT NULL DEFAULT NOW(),
    UNIQUE(event_id, email)
);

CREATE INDEX IF NOT EXISTS idx_participants_event_id ON participants(event_id);
