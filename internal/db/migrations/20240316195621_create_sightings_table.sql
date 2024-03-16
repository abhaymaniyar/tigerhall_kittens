-- +goose Up
-- +goose StatementBegin
CREATE TABLE sightings
(
    id                  VARCHAR(36) PRIMARY KEY,
    tiger_id            INTEGER                  NOT NULL,
    reported_by_user_id VARCHAR(36)              NOT NULL,
    image_url           VARCHAR(200)             DEFAULT NULL,
    sighted_at          TIMESTAMP WITH TIME ZONE NOT NULL,
    lat                 DECIMAL(9, 6),
    lon                 DECIMAL(9, 6),
    created_at          TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at          TIMESTAMP WITH TIME ZONE DEFAULT NULL,
    FOREIGN KEY (tiger_id) REFERENCES tigers (id) ON DELETE CASCADE,
    FOREIGN KEY (reported_by_user_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE INDEX idx_sightings_tiger_id ON sightings (tiger_id);
CREATE INDEX idx_sightings_user_id ON sightings (reported_by_user_id);
CREATE INDEX idx_sightings_lat ON sightings (lat);
CREATE INDEX idx_sightings_lon ON sightings (lon);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS sightings CASCADE;
-- +goose StatementEnd
