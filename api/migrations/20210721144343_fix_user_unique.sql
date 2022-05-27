-- +goose Up
-- SQL in this section is executed when the migration is applied.
-- +goose StatementBegin


ALTER TABLE core.users DROP CONSTRAINT users_username_key;
DROP INDEX IF EXISTS users_username_key;
ALTER TABLE core.users DROP CONSTRAINT users_email_key;
DROP INDEX IF EXISTS users_email_key;


create unique index if not exists users_username_key on core.users (username)
WHERE username != '';

create unique index if not exists users_email_key on core.users (email)
WHERE email != '';

-- +goose StatementEnd
-- +goose Down
-- SQL in this section is executed when the migration is rolled back.