-- +goose Up
-- +goose StatementBegin
CREATE TABLE trips (
                       id text NOT NULL PRIMARY KEY,
                       user_id text NOT NULL,
                       offer_id text,
                       status text NOT NULL
);
-- +goose StatementEnd

-- +goose Down
DROP TABLE trips;