-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS printers (
	id BIGSERIAL PRIMARY KEY,
	addr TEXT NOT NULL,
	type TEXT NOT NULL,
	comment TEXT NOT NULL default ''
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS printers;
-- +goose StatementEnd
