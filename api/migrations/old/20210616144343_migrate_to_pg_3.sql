-- +goose Up
-- SQL in this section is executed when the migration is applied.
-- +goose StatementBegin


create schema if not exists blog;

CREATE TABLE IF NOT EXISTS blog.posts (
    id serial NOT NULL,
    created_at    timestamptz NOT NULL DEFAULT now(),
    updated_at            timestamptz NOT NULL DEFAULT now(),

    title varchar(255) not null,
    content varchar not null default '',
    image_urls varchar not null default '',
    tickers varchar not null default '',
    published_at            timestamptz NULL,
    posted_by int NOT NULL  REFERENCES core.users(id),

    PRIMARY KEY (id)  
);

DROP TRIGGER IF EXISTS set_timestamp on blog.posts;
create trigger set_timestamp before update
    on
    blog.posts for each row execute procedure trigger_set_timestamp();


CREATE TABLE IF NOT EXISTS blog.comments (
    id serial NOT NULL,
    created_at    timestamptz NOT NULL DEFAULT now(),
    updated_at            timestamptz NOT NULL DEFAULT now(),

    content varchar not null default '',
    published_at            timestamptz NULL,
    posted_by int NOT NULL  REFERENCES core.users(id),
    post_id int NOT NULL  REFERENCES blog.posts(id),

    PRIMARY KEY (id)  
);

DROP TRIGGER IF EXISTS set_timestamp on blog.comments;
create trigger set_timestamp before update
    on
    blog.comments for each row execute procedure trigger_set_timestamp();


CREATE TABLE IF NOT EXISTS blog.post_votes (
    id serial NOT NULL,
    created_at    timestamptz NOT NULL DEFAULT now(),
    updated_at            timestamptz NOT NULL DEFAULT now(),

    post_id int NOT NULL  REFERENCES blog.posts(id),
    user_id int NOT NULL  REFERENCES core.users(id),
    value int NOT NULL,

    PRIMARY KEY (id)  
);

DROP TRIGGER IF EXISTS set_timestamp on blog.post_votes;
create trigger set_timestamp before update
    on
    blog.post_votes for each row execute procedure trigger_set_timestamp();


create unique index if not exists post_votes_post_id_user_id_uidx on blog.post_votes (post_id,user_id);



CREATE TABLE IF NOT EXISTS blog.threads (
    id serial NOT NULL,
    created_at    timestamptz NOT NULL DEFAULT now(),
    updated_at            timestamptz NOT NULL DEFAULT now(),

    title varchar(255) not null,
    thread_type varchar(6) not null,

    PRIMARY KEY (id)  
);

DROP TRIGGER IF EXISTS set_timestamp on blog.threads;
create trigger set_timestamp before update
    on
    blog.threads for each row execute procedure trigger_set_timestamp();


CREATE TABLE IF NOT EXISTS blog.thread_messages (
    id serial NOT NULL,
    created_at    timestamptz NOT NULL DEFAULT now(),

    thread_id int NOT NULL  REFERENCES blog.threads(id),
    sender_id int NOT NULL  REFERENCES core.users(id),
    message_type varchar(10) not null,
    message varchar not null,
    timestamp bigint not null,

    PRIMARY KEY (id)  
);


CREATE TABLE IF NOT EXISTS blog.thread_participants (
    id serial NOT NULL,
    created_at    timestamptz NOT NULL DEFAULT now(),
    updated_at            timestamptz NOT NULL DEFAULT now(),

    thread_id int NOT NULL  REFERENCES blog.threads(id),
    user_id int NOT NULL  REFERENCES core.users(id),

    PRIMARY KEY (id)  
);

create unique index if not exists thread_participants_thread_id_user_id_uidx on blog.thread_participants (thread_id,user_id);

DROP TRIGGER IF EXISTS set_timestamp on blog.thread_participants;
create trigger set_timestamp before update
    on
    blog.thread_participants for each row execute procedure trigger_set_timestamp();


-- +goose StatementEnd
-- +goose Down
-- SQL in this section is executed when the migration is rolled back.