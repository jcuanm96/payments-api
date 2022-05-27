-- +goose Up
-- SQL in this section is executed when the migration is applied.
-- +goose StatementBegin


create schema if not exists subscription;

CREATE TABLE IF NOT EXISTS  subscription.types(
    id int NOT NULL,
    internal_name varchar(50) not null,
    full_name varchar(100) not null,

    PRIMARY KEY (id)  
);

INSERT INTO subscription.types (id, internal_name,full_name) VALUES(1, 'batch_of_messages', 'Batch of Messages') on conflict do nothing;
INSERT INTO subscription.types (id, internal_name,full_name) VALUES(2, 'monthly_subscription','Monthly Subscription') on conflict do nothing;


CREATE TABLE IF NOT EXISTS subscription.prices (
    id serial NOT NULL,
    created_at    timestamptz NOT NULL DEFAULT now(),
    updated_at            timestamptz NOT NULL DEFAULT now(),

    subscription_type_id int not null  REFERENCES subscription.types(id),
    seller_user_id int not null  REFERENCES core.users(id),

    is_active boolean NOT NULL DEFAULT false,
    amount float not null,

    PRIMARY KEY (id)  
);
create unique index if not exists prices_subscription_type_id_seller_user_id_uidx on subscription.prices (subscription_type_id,seller_user_id);

DROP TRIGGER IF EXISTS set_timestamp on subscription.prices;
create trigger set_timestamp before update
    on
    subscription.prices for each row execute procedure trigger_set_timestamp();

CREATE TABLE IF NOT EXISTS subscription.receipts (
    id serial NOT NULL,
    uuid varchar(36) not null UNIQUE,
    created_at    timestamptz NOT NULL DEFAULT now(),
    updated_at            timestamptz NOT NULL DEFAULT now(),

    subscription_price_id int not null  REFERENCES subscription.prices(id),
    seller_user_id int not null  REFERENCES core.users(id),
    buyer_user_id int not null  REFERENCES core.users(id),

    is_refunded boolean NOT NULL DEFAULT false,
    is_confirmed boolean NOT NULL DEFAULT false,
    amount float not null,

    PRIMARY KEY (id)  
);

DROP TRIGGER IF EXISTS set_timestamp on subscription.receipts;
create trigger set_timestamp before update
    on
    subscription.receipts for each row execute procedure trigger_set_timestamp();

CREATE TABLE IF NOT EXISTS subscription.packages (
    id serial NOT NULL,
    created_at    timestamptz NOT NULL DEFAULT now(),
    updated_at       timestamptz NOT NULL DEFAULT now(),

    subscription_type_id int  not null REFERENCES subscription.types(id),
    seller_user_id int  not null  REFERENCES core.users(id),
    buyer_user_id int  not null  REFERENCES core.users(id),

    is_active boolean NOT NULL DEFAULT false,

    PRIMARY KEY (id)  
);
create unique index if not exists packages_subscription_type_id_seller_user_id_buyer_user_id_uidx on subscription.packages (subscription_type_id,seller_user_id,buyer_user_id);

DROP TRIGGER IF EXISTS set_timestamp on subscription.packages;
create trigger set_timestamp before update
    on
    subscription.packages for each row execute procedure trigger_set_timestamp();



-- +goose StatementEnd
-- +goose Down
-- SQL in this section is executed when the migration is rolled back.