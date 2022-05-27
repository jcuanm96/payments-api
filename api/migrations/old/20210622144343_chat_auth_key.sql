-- +goose Up
-- SQL in this section is executed when the migration is applied.
-- +goose StatementBegin

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

alter table core.users add column if not exists chat_key uuid DEFAULT uuid_generate_v4 () UNIQUE;


-- +goose StatementEnd
-- +goose Down
-- SQL in this section is executed when the migration is rolled back.