-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS public.usr_users (
	user_id uuid DEFAULT gen_random_uuid() NOT NULL,
	login varchar NOT NULL,
    hashed_password varchar NOT NULL,
	created_at timestamp without time zone DEFAULT current_timestamp NOT NULL,
	updated_at timestamp without time zone DEFAULT current_timestamp NOT NULL,
	deleted_at timestamp without time zone NULL,
	CONSTRAINT usr_users_pk PRIMARY KEY (user_id),
	CONSTRAINT usr_users_unique UNIQUE (login)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS public.usr_users;
-- +goose StatementEnd
