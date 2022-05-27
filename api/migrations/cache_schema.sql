-- +goose Up
-- SQL in this section is executed when the migration is applied.
-- +goose StatementBegin


CREATE TABLE IF NOT EXISTS cache (
    key     varchar(100) not null,
    prefix  varchar not null,
    expired_at  bigint not null ,
    model  varchar not null,
    user_id int null DEFAULT NULL,

    PRIMARY KEY (key,prefix)  
);


-- +goose StatementEnd
-- +goose Down
-- SQL in this section is executed when the migration is rolled back.