-- Feed schema
create schema if not exists subscription;

-- Creator Subscription Tiers Table
CREATE TABLE IF NOT EXISTS subscription.tiers (
  id serial NOT NULL,
  goat_user_id int NOT NULL REFERENCES core.users(id),
  price_in_smallest_denom int NOT NULL DEFAULT 500, -- defaulting to 5 dollars
  currency varchar(3) NOT NULL DEFAULT 'usd',
  tier_name VARCHAR(20) NOT NULL DEFAULT 'TIERLESS',
  stripe_product_id VARCHAR(50) NOT NULL,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  deleted_at timestamptz NULL,
  PRIMARY KEY (id),
  UNIQUE (goat_user_id, tier_name)
);

DROP TRIGGER IF EXISTS set_timestamp on subscription.tiers;
create trigger set_timestamp before update
    on
    subscription.tiers for each row execute procedure trigger_set_timestamp();

-- User Subscriptions Table
CREATE TABLE IF NOT EXISTS subscription.user_subscriptions (
  id serial NOT NULL,
  current_period_end timestamptz NOT NULL DEFAULT now() + '1 month',
  user_id int NOT NULL REFERENCES core.users(id),
  stripe_subscription_id VARCHAR(50) NOT NULL DEFAULT '',
  goat_user_id int NOT NULL REFERENCES core.users(id),
  tier_id int NOT NULL REFERENCES subscription.tiers(id),
  is_renewing boolean NOT NULL DEFAULT TRUE,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  deleted_at timestamptz NULL,
  PRIMARY KEY (user_id, goat_user_id)
);

DROP TRIGGER IF EXISTS set_timestamp on subscription.user_subscriptions;
create trigger set_timestamp before update
    on
    subscription.user_subscriptions for each row execute procedure trigger_set_timestamp();

drop index if exists subscription.subscription_user_subscriptions_user_id_uidx;
create index subscription_user_subscriptions_user_id_uidx on subscription.user_subscriptions (user_id);

CREATE TABLE IF NOT EXISTS subscription.paid_group_chat_subscriptions (
  id serial NOT NULL,
  current_period_end timestamptz NOT NULL DEFAULT now() + '1 month',
  user_id int NOT NULL REFERENCES core.users(id),
  goat_user_id int NOT NULL REFERENCES core.users(id),
  stripe_subscription_id VARCHAR(50) NOT NULL DEFAULT '',
  sendbird_channel_id varchar(200) NOT NULL,
  is_renewing boolean NOT NULL DEFAULT TRUE,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY (user_id, sendbird_channel_id)
);

DROP TRIGGER IF EXISTS set_timestamp on subscription.paid_group_chat_subscriptions;
create trigger set_timestamp before update
    on
    subscription.paid_group_chat_subscriptions for each row execute procedure trigger_set_timestamp();