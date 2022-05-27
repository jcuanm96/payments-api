-- +goose Up
-- SQL in this section is executed when the migration is applied.

alter table core.users add column last_active timestamp NOT NULL DEFAULT now();

update core.users set last_active = last_seen;

drop TRIGGER users_updated_at;

CREATE TRIGGER users_updated_at BEFORE UPDATE ON core.users
FOR EACH ROW
BEGIN
    IF NEW.last_active = OLD.last_active THEN
      SET  NEW.updated_at  = now();
    END IF;
END

DROP VIEW IF EXISTS contacts;

CREATE VIEW contacts AS
SELECT
	users_contacts.id,
  users_contacts.user_id,
  users_contacts.contact_id,
  users.first_name,
  users.last_name,
  users.phone_number,
  users.country_code,
  users.email,
  users.user_type,
  users.profile_avatar,
  users.last_active,
  users_contacts.created_at,
  users_contacts.updated_at
FROM users_contacts JOIN users ON users_contacts.contact_id = users.id
WHERE users.deleted_at IS NULL;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.