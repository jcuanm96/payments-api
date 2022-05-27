create schema if not exists push;

CREATE TABLE IF NOT EXISTS push.settings (
    id                   serial NOT NULL,
    user_id              int NOT NULL REFERENCES core.users(id),
    pending_balance      varchar(10) NOT NULL DEFAULT 'UNSET',
    created_at           timestamptz NOT NULL DEFAULT now(),
    updated_at           timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (user_id)
);

DROP TRIGGER IF EXISTS set_timestamp on push.settings;
create trigger set_timestamp before update
    on 
    push.settings for each row execute procedure trigger_set_timestamp();