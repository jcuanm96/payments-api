-- +goose Up
-- SQL in this section is executed when the migration is applied.
-- +goose StatementBegin

CREATE OR REPLACE FUNCTION public.trigger_set_timestamp()
RETURNS trigger
LANGUAGE plpgsql
AS $function$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$function$
;
-- +goose StatementEnd
-- +goose Down
-- SQL in this section is executed when the migration is rolled back.