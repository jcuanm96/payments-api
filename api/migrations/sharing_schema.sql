CREATE TABLE IF NOT EXISTS core.user_bio (
  user_id int NOT NULL REFERENCES core.users(id) UNIQUE,
  text_content VARCHAR(1000) NOT NULL DEFAULT '', 
  PRIMARY KEY (user_id)
);

create schema if not exists sharing;

CREATE TABLE IF NOT EXISTS sharing.message_links (
  id serial NOT NULL UNIQUE,
  link_suffix VARCHAR(200) NOT NULL DEFAULT '',
  sendbird_channel_id VARCHAR(200) NOT NULL DEFAULT '',
  message_id VARCHAR(200) NOT NULL DEFAULT '',
  PRIMARY KEY (link_suffix)
);

CREATE TABLE IF NOT EXISTS sharing.themes (
  id serial NOT NULL UNIQUE,
  theme_name VARCHAR(150) NOT NULL DEFAULT '' UNIQUE,
  theme_vama_logo_color VARCHAR(6) NOT NULL DEFAULT '000000',
  theme_icon_color VARCHAR(6) NOT NULL DEFAULT '000000',
  theme_top_gradient_color VARCHAR(6) NOT NULL DEFAULT '000000',
  theme_bottom_gradient_color VARCHAR(6) NOT NULL DEFAULT 'FFFFFF',
  theme_username_color VARCHAR(6) NOT NULL DEFAULT '000000',
  theme_bio_color VARCHAR(6) NOT NULL DEFAULT '000000',
  theme_row_color VARCHAR(6) NOT NULL DEFAULT 'FFFFFF',
  theme_row_text_color VARCHAR(6) NOT NULL DEFAULT '000000',
  PRIMARY KEY (theme_name)
);

CREATE TABLE IF NOT EXISTS sharing.bio_links (
  goat_id int NOT NULL REFERENCES core.users(id) UNIQUE,
  text_contents VARCHAR(1254) NOT NULL DEFAULT '', 
  links VARCHAR(1254) NOT NULL DEFAULT '',
  theme_id int NOT NULL REFERENCES sharing.themes(id) DEFAULT 1,
  PRIMARY KEY (goat_id)
);

