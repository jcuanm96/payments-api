-- +goose Up
-- SQL in this section is executed when the migration is applied.
-- +goose StatementBegin

-- Core schema
create schema if not exists core;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS core.users (
    id serial NOT NULL,
    created_at    timestamptz NOT NULL DEFAULT now(),
    updated_at            timestamptz NOT NULL DEFAULT now(),
    stripe_id varchar(50) null UNIQUE,
    first_name varchar(50) not null default '',
    last_name varchar(50) not null default '',
    phone_number varchar(40) not null UNIQUE,
    country_code varchar(2) not null default '',
    username varchar(50) not null UNIQUE default '',
    email varchar(50) null UNIQUE,
    user_type varchar(10) not null,
    profile_avatar varchar null ,
    deleted_at            timestamptz NULL,
    stripe_account_id varchar(50) null UNIQUE,
    uuid UUID DEFAULT uuid_generate_v4 () UNIQUE,
    PRIMARY KEY (id)  
);

DROP TRIGGER IF EXISTS set_timestamp on core.users;
create trigger set_timestamp before update
    on
    core.users for each row execute procedure trigger_set_timestamp();

create unique index if not exists users_stripe_account_id_uidx on core.users (stripe_account_id);
create unique index if not exists users_stripe_id_uidx on core.users (stripe_id);
create unique index if not exists users_email_uidx on core.users (email);
create unique index if not exists users_phone_number_uidx on core.users (phone_number);
create unique index if not exists users_username_uidx on core.users (username);
create unique index if not exists users_uuid_uidx on core.users (uuid);


CREATE TABLE IF NOT EXISTS core.tokens (
    id serial NOT NULL,
    user_id int NOT NULL REFERENCES core.users(id),
    fcm_token varchar(500),
    created_at    timestamptz NOT NULL DEFAULT now(),
    updated_at    timestamptz NOT NULL DEFAULT now(),
    deleted_at    timestamptz NULL,
    PRIMARY KEY (user_id)
);

DROP TRIGGER IF EXISTS set_timestamp on core.tokens;
create trigger set_timestamp before update
    on
    core.tokens for each row execute procedure trigger_set_timestamp();

CREATE TABLE IF NOT EXISTS core.goat_invite_codes (
    id serial NOT NULL,
    created_at    timestamptz NOT NULL DEFAULT now(),
    updated_at    timestamptz NOT NULL DEFAULT now(),
    invite_code varchar(6) NOT NULL UNIQUE,
    used_by int NULL  REFERENCES core.users(id) UNIQUE,
    invited_by int NULL  REFERENCES core.users(id),
    PRIMARY KEY (id)  
);

DROP TRIGGER IF EXISTS set_timestamp on core.goat_invite_codes;
create trigger set_timestamp before update
    on
    core.goat_invite_codes for each row execute procedure trigger_set_timestamp();


CREATE TABLE IF NOT EXISTS core.users_contacts (
    id serial NOT NULL,
    created_at    timestamptz NOT NULL DEFAULT now(),
    updated_at            timestamptz NOT NULL DEFAULT now(),
    user_id int NOT NULL   REFERENCES core.users(id),
    contact_id int NOT NULL   REFERENCES core.users(id),
    PRIMARY KEY (user_id, contact_id)  
);
DROP TRIGGER IF EXISTS set_timestamp on core.users_contacts;
create trigger set_timestamp before update
    on
    core.users_contacts for each row execute procedure trigger_set_timestamp();

CREATE TABLE IF NOT EXISTS core.pending_contacts (
    id serial NOT NULL,
    created_at    timestamptz NOT NULL DEFAULT now(),
    updated_at            timestamptz NOT NULL DEFAULT now(),
    user_id int NOT NULL   REFERENCES core.users(id),
    signed_up_user_id int default null,
    phone_number varchar(40) not null,
    country_code varchar(2) not null default '',
    first_name varchar(50) not null default '',
    last_name varchar(50) not null default '',
    PRIMARY KEY (user_id, phone_number)
);

DROP TRIGGER IF EXISTS set_timestamp on core.pending_contacts;
create trigger set_timestamp before update
    on
    core.pending_contacts for each row execute procedure trigger_set_timestamp();

create index if not exists core_pending_contacts_signed_up_user_id_uidx on core.pending_contacts (signed_up_user_id);

CREATE TABLE IF NOT EXISTS core.user_blocks (
    id serial NOT NULL,
    created_at    timestamptz NOT NULL DEFAULT now(),
    updated_at            timestamptz NOT NULL DEFAULT now(),
    user_id int NOT NULL   REFERENCES core.users(id),
    blocked_user_id int NOT NULL   REFERENCES core.users(id),
    PRIMARY KEY (user_id, blocked_user_id)  
);

DROP TRIGGER IF EXISTS set_timestamp on core.user_blocks;
create trigger set_timestamp before update
    on
    core.user_blocks for each row execute procedure trigger_set_timestamp();


CREATE TABLE IF NOT EXISTS core.user_reports (
    id serial NOT NULL,
    created_at    timestamptz NOT NULL DEFAULT now(),
    updated_at            timestamptz NOT NULL DEFAULT now(),
    reporter_user_id int NOT NULL   REFERENCES core.users(id),
    reported_user_id int NOT NULL   REFERENCES core.users(id),
    description varchar(1000) NULL DEFAULT '',
    PRIMARY KEY (id)  
);

DROP TRIGGER IF EXISTS set_timestamp on core.user_reports;
create trigger set_timestamp before update
    on
    core.user_reports for each row execute procedure trigger_set_timestamp();

create index if not exists core_user_reports_reporter_user_id_uidx on core.user_reports (reporter_user_id);
create index if not exists core_user_reports_reported_user_id_uidx on core.user_reports (reported_user_id);


CREATE TABLE IF NOT EXISTS core.banned_chat_users (
    id serial NOT NULL,
    banned_user_id int NOT NULL   REFERENCES core.users(id),
    user_id int NOT NULL   REFERENCES core.users(id),
    sendbird_channel_id varchar(200) NOT NULL,
    created_at    timestamptz NOT NULL DEFAULT now(),
    updated_at    timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (sendbird_channel_id, banned_user_id)  
);

DROP TRIGGER IF EXISTS set_timestamp on core.banned_chat_users;
create trigger set_timestamp before update
    on
    core.banned_chat_users for each row execute procedure trigger_set_timestamp();


-- Wallet schema
create schema if not exists wallet;

CREATE TABLE IF NOT EXISTS wallet.payout_periods (
    id serial NOT NULL,
    start_ts bigint,
    end_ts bigint,
    created_at    timestamptz NOT NULL DEFAULT now(),
    updated_at    timestamptz NOT NULL DEFAULT now(),
    deleted_at    timestamptz NULL,
    PRIMARY KEY (id)
);

DROP TRIGGER IF EXISTS set_timestamp on wallet.payout_periods;
create trigger set_timestamp before update
    on
    wallet.payout_periods for each row execute procedure trigger_set_timestamp();

-- Initializing wallet.payout_periods table
do $$
	DECLARE INITIAL_PAYOUT_DATE timestamptz = to_date('2022-01-26','YYYY-MM-DD');
	DECLARE NEXT_PAYOUT_START_TS bigint;
	DECLARE NEXT_PAYOUT_END_TS bigint;
	begin
	IF EXISTS (
		SELECT 1
		FROM   information_schema.tables 
		WHERE  table_schema = 'wallet'
		AND    table_name = 'payout_periods'
	) 
	THEN 
		for i in 0..10000 loop
			NEXT_PAYOUT_START_TS = extract(epoch FROM INITIAL_PAYOUT_DATE + interval '1 days' * 7 * 2 * i);
			NEXT_PAYOUT_END_TS = extract(epoch FROM INITIAL_PAYOUT_DATE + interval '1 days' * 7 * 2 * i + interval '14 days');
			INSERT INTO wallet.payout_periods(start_ts, end_ts) 
			VALUES(NEXT_PAYOUT_START_TS, NEXT_PAYOUT_END_TS);
		end loop;
	end if;
	end;
$$;

CREATE TABLE IF NOT EXISTS wallet.balances (
    id serial NOT NULL UNIQUE,
    provider_user_id int NOT NULL REFERENCES core.users(id),
    available_balance bigint NOT NULL DEFAULT 0,
    last_payout_ts timestamptz NOT NULL DEFAULT now(),
    last_paid_payout_period_id int NULL REFERENCES wallet.payout_periods(id),
    currency varchar(3) NOT NULL default 'usd',
    created_at    timestamptz NOT NULL DEFAULT now(),
    updated_at    timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY(provider_user_id, currency)
);

DROP TRIGGER IF EXISTS set_timestamp on wallet.balances;
create trigger set_timestamp before update
    on
    wallet.balances for each row execute procedure trigger_set_timestamp();


CREATE TABLE IF NOT EXISTS wallet.pending_transactions (
    id serial NOT NULL UNIQUE,
    provider_user_id int REFERENCES core.users(id) NOT NULL,
    customer_user_id int REFERENCES core.users(id) NOT NULL,
    payment_intent_id varchar(50),
    stripe_created_ts bigint,
    amount bigint NOT NULL,
    currency varchar(3) NOT NULL DEFAULT 'usd',
    version varchar(50) NOT NULL,
    created_at    timestamptz NOT NULL DEFAULT now(),
    updated_at    timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (id)
);


DROP TRIGGER IF EXISTS set_timestamp on wallet.pending_transactions;
create trigger set_timestamp before update
    on
    wallet.pending_transactions for each row execute procedure trigger_set_timestamp();

CREATE TABLE IF NOT EXISTS wallet.ledger (
    id serial NOT NULL UNIQUE,
    provider_user_id int REFERENCES core.users(id) NOT NULL,
    customer_user_id int REFERENCES core.users(id),
    stripe_transaction_id varchar(50),
    source_type varchar(30),
    stripe_created_ts bigint,
    stripe_fee bigint,
    vama_fee bigint,
    amount bigint NOT NULL,
    currency varchar(3) NOT NULL DEFAULT 'usd',
    pay_period_id int,
    version varchar(50) NOT NULL,
    created_at    timestamptz NOT NULL DEFAULT now(),
    updated_at    timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (id)
);


DROP TRIGGER IF EXISTS set_timestamp on wallet.ledger;
create trigger set_timestamp before update
    on
    wallet.ledger for each row execute procedure trigger_set_timestamp();

CREATE TABLE IF NOT EXISTS wallet.bank_info (
    id serial NOT NULL UNIQUE,
    user_id int REFERENCES core.users(id) NOT NULL,
    bank_name varchar(200) NOT NULL,
    account_number varchar(200) NOT NULL,
    routing_number varchar(200) NOT NULL,
    account_type varchar(100) NOT NULL,
    account_holder_name varchar(200) NOT NULL,
    account_holder_type varchar(100) NOT NULL,
    currency varchar(3) NOT NULL DEFAULT 'usd',
    country varchar(200) NOT NULL,
    created_at    timestamptz NOT NULL DEFAULT now(),
    updated_at    timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (id)
);

DROP TRIGGER IF EXISTS set_timestamp on wallet.bank_info;
create trigger set_timestamp before update
    on
    wallet.bank_info for each row execute procedure trigger_set_timestamp();

create unique index if not exists wallet_bank_info_user_id_uidx on wallet.bank_info (user_id);

CREATE TABLE IF NOT EXISTS wallet.billing_addresses (
    id serial NOT NULL,
    user_id int REFERENCES core.users(id) NOT NULL,
    street_1 varchar(200) NOT NULL,
    street_2 varchar(200),
    city varchar(200) NOT NULL,
    state varchar(200) NOT NULL,
    postal_code varchar(50) NOT NULL,
    country varchar(200) NOT NULL,
    created_at    timestamptz NOT NULL DEFAULT now(),
    updated_at    timestamptz NOT NULL DEFAULT now(),
    deleted_at    timestamptz NULL,
    PRIMARY KEY (id)
);

DROP TRIGGER IF EXISTS set_timestamp on wallet.billing_addresses;
create trigger set_timestamp before update
    on
    wallet.billing_addresses for each row execute procedure trigger_set_timestamp();

CREATE TABLE IF NOT EXISTS wallet.payout_periods (
    id serial NOT NULL,
    start_ts bigint,
    end_ts bigint,
    created_at    timestamptz NOT NULL DEFAULT now(),
    updated_at    timestamptz NOT NULL DEFAULT now(),
    deleted_at    timestamptz NULL,
    PRIMARY KEY (id)
);

DROP TRIGGER IF EXISTS set_timestamp on wallet.payout_periods;
create trigger set_timestamp before update
    on
    wallet.payout_periods for each row execute procedure trigger_set_timestamp();

-- Initializing wallet.payout_periods table
do $$
	DECLARE INITIAL_PAYOUT_DATE timestamptz = to_date('2022-02-09','YYYY-MM-DD');
	DECLARE NEXT_PAYOUT_START_TS bigint;
	DECLARE NEXT_PAYOUT_END_TS bigint;
	begin
	IF EXISTS (
		SELECT 1
		FROM   information_schema.tables 
		WHERE  table_schema = 'wallet'
		AND    table_name = 'payout_periods'
	) 
	THEN 
		for i in 0..10000 loop
			NEXT_PAYOUT_START_TS = extract(epoch FROM INITIAL_PAYOUT_DATE + interval '1 days' * 7 * 2 * i);
			NEXT_PAYOUT_END_TS = extract(epoch FROM INITIAL_PAYOUT_DATE + interval '1 days' * 7 * 2 * i + interval '14 days');
			INSERT INTO wallet.payout_periods(start_ts, end_ts) 
			VALUES(NEXT_PAYOUT_START_TS, NEXT_PAYOUT_END_TS);
		end loop;
	end if;
	end;
$$;

-- Product schema
create schema if not exists product;

CREATE TABLE IF NOT EXISTS product.goat_chats (
    id serial NOT NULL,
    price_in_smallest_denom int NOT NULL DEFAULT 300, -- defaulting to 3 dollars
    currency varchar(3) NOT NULL DEFAULT 'usd',
    goat_user_id int NOT NULL REFERENCES core.users(id),
    created_at    timestamptz NOT NULL DEFAULT now(),
    updated_at    timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (id)  
);

create unique index if not exists product_goat_chats_goat_user_id_uidx on product.goat_chats (goat_user_id);

DROP TRIGGER IF EXISTS set_timestamp on product.goat_chats;
create trigger set_timestamp before update
    on
    product.goat_chats for each row execute procedure trigger_set_timestamp();

CREATE TABLE IF NOT EXISTS product.paid_group_chats (
    id serial NOT NULL,
    price_in_smallest_denom int NOT NULL,
    currency varchar(3) NOT NULL DEFAULT 'usd',
    goat_user_id int NOT NULL REFERENCES core.users(id),
    sendbird_channel_id varchar(200) NOT NULL,
    stripe_product_id varchar(200) NOT NULL,
    link_suffix varchar(400) DEFAULT NULL UNIQUE,
    member_limit int NOT NULL DEFAULT 100,
    is_member_limit_enabled BOOLEAN NOT NULL DEFAULT FALSE,
    metadata json default '{}',
    created_at    timestamptz NOT NULL DEFAULT now(),
    updated_at    timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (sendbird_channel_id)
);

DROP TRIGGER IF EXISTS set_timestamp on product.paid_group_chats;
create trigger set_timestamp before update
    on
    product.paid_group_chats for each row execute procedure trigger_set_timestamp();

create unique index if not exists product_paid_group_chats_link_suffix on product.paid_group_chats(link_suffix);

CREATE TABLE IF NOT EXISTS product.free_group_chats (
    id serial NOT NULL,
    creator_user_id int NOT NULL REFERENCES core.users(id),
    sendbird_channel_id varchar(200) NOT NULL,
    link_suffix varchar(400) DEFAULT NULL UNIQUE,
    member_limit int NOT NULL DEFAULT 100,
    is_member_limit_enabled BOOLEAN NOT NULL DEFAULT FALSE,
    metadata json default '{}',
    created_at    timestamptz NOT NULL DEFAULT now(),
    updated_at    timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (sendbird_channel_id)
);

DROP TRIGGER IF EXISTS set_timestamp on product.free_group_chats;
create trigger set_timestamp before update
    on
    product.free_group_chats for each row execute procedure trigger_set_timestamp();

create unique index if not exists product_free_group_chats_link_suffix on product.free_group_chats(link_suffix);


CREATE TABLE IF NOT EXISTS product.free_group_creators (
     id serial NOT NULL,
     sendbird_channel_id varchar(200) NOT NULL,
     creator_user_id int NOT NULL REFERENCES core.users(id),
     created_at    timestamptz NOT NULL DEFAULT now(),
     updated_at    timestamptz NOT NULL DEFAULT now(),
     PRIMARY KEY (sendbird_channel_id, creator_user_id)
);

DROP TRIGGER IF EXISTS set_timestamp on product.free_group_creators;
create trigger set_timestamp before update
    on
    product.free_group_creators for each row execute procedure trigger_set_timestamp();

-- +goose StatementEnd
-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
