-- +goose Up
-- SQL in this section is executed when the migration is applied.

create TABLE subscription_types(
  id int unsigned  PRIMARY KEY,
  internal_name varchar(50) not null,
  full_name varchar(100) not null,
);


alter table subscription.types add column full_name varchar(100) not null;

INSERT INTO subscription_types (id, internal_name,full_name) VALUES(1, 'batch_of_messages', 'Batch of Messages');
INSERT INTO subscription_types (id, internal_name,full_name) VALUES(2, 'monthly_subscription','Monthly Subscription');


CREATE TABLE subscription_prices (
  id int unsigned AUTO_INCREMENT PRIMARY KEY,
  created_at timestamp NOT NULL DEFAULT now(),
  updated_at timestamp NOT NULL DEFAULT now() ON UPDATE now(),

  subscription_type_id int unsigned not null, FOREIGN KEY (subscription_type_id) REFERENCES subscription_types (id),
  seller_user_id int unsigned not null, FOREIGN KEY (seller_user_id) REFERENCES users (id),
  UNIQUE KEY(subscription_type_id,seller_user_id),

  is_active boolean NOT NULL DEFAULT false,
  amount DOUBLE(10,2) not null
);

CREATE TABLE subscription_receipts (
  id int unsigned AUTO_INCREMENT PRIMARY KEY,
  uuid varchar(36) not null UNIQUE,
  created_at timestamp NOT NULL DEFAULT now(),
  updated_at timestamp NOT NULL DEFAULT now() ON UPDATE now(),

  subscription_price_id int unsigned not null, FOREIGN KEY (subscription_price_id) REFERENCES subscription_prices (id),
  seller_user_id int unsigned not null, FOREIGN KEY (seller_user_id) REFERENCES users (id),
  buyer_user_id int unsigned not null, FOREIGN KEY (buyer_user_id) REFERENCES users (id),

  is_refunded boolean NOT NULL DEFAULT false,
  is_confirmed boolean NOT NULL DEFAULT false,
  amount DOUBLE(10,2) not null
);

CREATE TABLE subscription_packages (
  id int unsigned AUTO_INCREMENT PRIMARY KEY,
  created_at timestamp NOT NULL DEFAULT now(),
  updated_at timestamp NOT NULL DEFAULT now() ON UPDATE now(),

  subscription_type_id int unsigned not null, FOREIGN KEY (subscription_type_id) REFERENCES subscription_types (id),
  seller_user_id int unsigned not null, FOREIGN KEY (seller_user_id) REFERENCES users (id),
  buyer_user_id int unsigned not null, FOREIGN KEY (buyer_user_id) REFERENCES users (id),

  is_active boolean NOT NULL DEFAULT false,

  UNIQUE KEY(subscription_type_id,seller_user_id,buyer_user_id)
);

CREATE TABLE wallets (
  id int unsigned AUTO_INCREMENT PRIMARY KEY,
  created_at timestamp NOT NULL DEFAULT now(),
  updated_at timestamp NOT NULL DEFAULT now() ON UPDATE now(),

  owner_user_id int unsigned not null UNIQUE, FOREIGN KEY (owner_user_id) REFERENCES users (id),

  balance DOUBLE(10,2) not null
);

CREATE TABLE wallet_history (
  id int unsigned AUTO_INCREMENT PRIMARY KEY,
  idempotent_key varchar(100) not null UNIQUE,
  created_at timestamp NOT NULL DEFAULT now(),

  wallet_id int unsigned not null UNIQUE, FOREIGN KEY (wallet_id) REFERENCES wallets (id),

  direction  varchar(5) not null,
  amount DOUBLE(10,2) not null
);


-- +goose Down
-- SQL in this section is executed when the migration is rolled back.