-- +goose Up
CREATE TABLE link (
   uuid UUID PRIMARY KEY,
   original_url TEXT NOT NULL,
   short_url TEXT UNIQUE NOT NULL,
   created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

-- +goose Down
DROP TABLE link;
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
