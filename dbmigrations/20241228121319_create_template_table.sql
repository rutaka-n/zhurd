-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS templates (
	id BIGSERIAL PRIMARY KEY,
	label_id BIGINT NOT NULL references labels(id),
	type TEXT NOT NULL,
	body BYTEA NOT NULL default ''
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS templates;
-- +goose StatementEnd

