-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS labels (
	id BIGSERIAL PRIMARY KEY,
	name TEXT NOT NULL,
	comment TEXT NOT NULL default ''
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS labels;
-- +goose StatementEnd
