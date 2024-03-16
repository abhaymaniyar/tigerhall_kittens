-- +goose Up
-- +goose StatementBegin
CREATE TABLE tigers
(
    id                  SERIAL PRIMARY KEY,
    name                VARCHAR(255) NOT NULL,
    date_of_birth       DATE,
    last_seen_timestamp TIMESTAMP WITH TIME ZONE,
    last_seen_lat       DECIMAL(9, 6),
    last_seen_lon       DECIMAL(9, 6),
    created_at          TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at          TIMESTAMP WITH TIME ZONE DEFAULT NULL,
    deleted_at          TIMESTAMP WITH TIME ZONE DEFAULT NULL
);

CREATE INDEX idx_tigers_last_seen_lat ON tigers (last_seen_lat);
CREATE INDEX idx_tigers_last_seen_lon ON tigers (last_seen_lon);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS tigers CASCADE;
-- +goose StatementEnd
