-- +goose Up
CREATE UNIQUE INDEX idx_original_url_unique ON link (original_url);
-- +goose StatementBegin
SELECT 'Added unique index to original_url';
-- +goose StatementEnd

-- +goose Down
DROP INDEX idx_original_url_unique;
-- +goose StatementBegin
SELECT 'Removed unique index from original_url';
-- +goose StatementEnd
