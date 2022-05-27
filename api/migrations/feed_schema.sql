-- Feed schema
create schema if not exists feed;

-- Creator Chats Messages Table
CREATE TABLE IF NOT EXISTS feed.goat_chat_messages (
  id SERIAL NOT NULL,
  sendbird_channel_id VARCHAR(100) NOT NULL,
  conversation_start_ts BIGINT DEFAULT NULL,
  conversation_end_ts BIGINT DEFAULT NULL,
  is_public BOOLEAN NOT NULL DEFAULT TRUE,
  customer_user_id int NOT NULL REFERENCES core.users(id),
  provider_user_id int NOT NULL REFERENCES core.users(id),
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  deleted_at timestamptz NULL,
  PRIMARY KEY (id)
);

drop index if exists feed.feed_goat_chat_messages_sendbird_channel_id_uidx;
create index feed_goat_chat_messages_sendbird_channel_id_uidx on feed.goat_chat_messages (sendbird_channel_id);

DROP TRIGGER IF EXISTS set_timestamp on feed.goat_chat_messages;
create trigger set_timestamp before update
    on
    feed.goat_chat_messages for each row execute procedure trigger_set_timestamp();

-- Posts Table
CREATE TABLE IF NOT EXISTS feed.posts (
  id SERIAL NOT NULL,
  user_id int NOT NULL REFERENCES core.users(id),
  goat_chat_msgs_id int REFERENCES feed.goat_chat_messages(id) DEFAULT NULL,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  deleted_at timestamptz NULL,
  num_upvotes int NOT NULL DEFAULT 0,
  num_downvotes int NOT NULL DEFAULT 0,
  text_content VARCHAR(40000) NOT NULL DEFAULT '',
  link_suffix VARCHAR(200) UNIQUE,
  PRIMARY KEY (id)
);

DROP TRIGGER IF EXISTS set_timestamp on feed.posts;
create trigger set_timestamp before update
    on
    feed.posts for each row execute procedure trigger_set_timestamp();

drop index if exists feed.feed_posts_user_id_uidx;
create index feed_posts_user_id_uidx on feed.posts (user_id);

drop index if exists feed.feed_posts_link_suffix_uidx;
create unique index feed_posts_link_suffix_uidx on feed.posts(link_suffix);

-- Post Images Table
CREATE TABLE IF NOT EXISTS feed.post_images (
  id SERIAL NOT NULL,
  post_id int NOT NULL REFERENCES feed.posts(id) UNIQUE,
  img_url varchar(200) NOT NULL DEFAULT '',
  img_width int NOT NULL DEFAULT 0,
  img_height int NOT NULL DEFAULT 0,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  deleted_at timestamptz NULL,
  PRIMARY KEY (id)
);

drop index if exists feed.feed_post_images_post_id_uidx;
create index feed_post_images_post_id_uidx on feed.post_images (post_id);

DROP TRIGGER IF EXISTS set_timestamp on feed.post_images;
create trigger set_timestamp before update
    on
    feed.post_images for each row execute procedure trigger_set_timestamp();

-- Post Reactions Table
CREATE TABLE IF NOT EXISTS feed.post_reactions (
  post_id int NOT NULL REFERENCES feed.posts(id),
  user_id int NOT NULL REFERENCES core.users(id),
  type VARCHAR(20) NOT NULL,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  deleted_at timestamptz NULL,
  PRIMARY KEY (post_id, user_id)
);

drop index if exists feed.feed_post_reactions_post_id_uidx;
create index feed_post_reactions_post_id_uidx on feed.post_reactions (post_id);

DROP TRIGGER IF EXISTS set_timestamp on feed.post_reactions;
create trigger set_timestamp before update
    on
    feed.post_reactions for each row execute procedure trigger_set_timestamp();

-- Comments Table
CREATE TABLE IF NOT EXISTS feed.post_comments (
  id SERIAL NOT NULL,
  user_id int NOT NULL REFERENCES core.users(id),
  post_id int NOT NULL REFERENCES feed.posts(id),
  text_content VARCHAR(10000) NOT NULL,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  deleted_at timestamptz NULL,
  PRIMARY KEY (id)
);

DROP TRIGGER IF EXISTS set_timestamp on feed.post_comments;
create trigger set_timestamp before update
    on
    feed.post_comments for each row execute procedure trigger_set_timestamp();

drop index if exists feed.feed_comment_post_id_uidx;
create index feed_comment_post_id_uidx on feed.post_comments (post_id);

-- Follows Table
CREATE TABLE IF NOT EXISTS feed.follows (
    id serial NOT NULL,
    created_at    timestamptz NOT NULL DEFAULT now(),
    updated_at            timestamptz NOT NULL DEFAULT now(),
    user_id int NOT NULL   REFERENCES core.users(id),
    goat_user_id int NOT NULL   REFERENCES core.users(id),
    PRIMARY KEY (user_id, goat_user_id)
);

DROP TRIGGER IF EXISTS set_timestamp on feed.follows;
create trigger set_timestamp before update
    on
    feed.follows for each row execute procedure trigger_set_timestamp();


CREATE TABLE IF NOT EXISTS feed.post_notifications (
    id serial NOT NULL,
    created_at    timestamptz NOT NULL DEFAULT now(),
    updated_at            timestamptz NOT NULL DEFAULT now(),
    user_id int NOT NULL   REFERENCES core.users(id),
    goat_user_id int NOT NULL   REFERENCES core.users(id),
    PRIMARY KEY (user_id, goat_user_id)
);

DROP TRIGGER IF EXISTS set_timestamp on feed.post_notifications;
create trigger set_timestamp before update
    on
    feed.post_notifications for each row execute procedure trigger_set_timestamp();


-- User home feed view
CREATE OR REPLACE VIEW feed.feed_posts_view AS 
  SELECT 
  -- Post data
    posts.id AS id,
    posts.user_id AS user_id,
    posts.created_at AS post_created_at,
    posts.deleted_at AS post_deleted_at,
    posts.num_upvotes AS post_num_upvotes,
    posts.num_downvotes AS post_num_downvotes,
    posts.text_content AS post_text_content,
    posts.link_suffix AS post_link_suffix,
    -- Creator Chat data
    COALESCE(goat_chat_messages.sendbird_channel_id, '') AS goat_chat_messages_sendbird_id,
    COALESCE(goat_chat_messages.customer_user_id, 0) AS goat_chat_messages_customer_id,
    COALESCE(goat_chat_messages.conversation_start_ts, 0) AS goat_chat_messages_convo_start_ts,
    COALESCE(goat_chat_messages.conversation_end_ts, 0) AS goat_chat_messages_convo_end_ts,
    -- Comments data
    (SELECT COUNT(*) 
      FROM feed.post_comments as comments
      WHERE comments.post_id = posts.id AND comments.deleted_at IS NULL
    ) AS num_comments,
    -- Image data
    COALESCE(post_images.img_url, '')   AS post_img_url,
    COALESCE(post_images.img_width, 0)  AS post_img_width,
    COALESCE(post_images.img_height, 0) AS post_img_height
  FROM feed.posts AS posts
  LEFT JOIN feed.goat_chat_messages AS goat_chat_messages ON posts.goat_chat_msgs_id = goat_chat_messages.id
  LEFT JOIN feed.post_images AS post_images ON posts.id = post_images.post_id
  WHERE posts.deleted_at IS NULL;
