-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS public.bll_balance (
	balance_id uuid DEFAULT gen_random_uuid() NOT NULL,
    user_id uuid NOT NULL,
	current numeric DEFAULT 0 NOT NULL,
	withdrawn numeric DEFAULT 0 NOT NULL,
	created_at timestamp without time zone DEFAULT current_timestamp NOT NULL,
	updated_at timestamp without time zone DEFAULT current_timestamp NOT NULL,
	CONSTRAINT bll_balance_pk PRIMARY KEY (balance_id)
);

ALTER TABLE public.bll_balance ADD CONSTRAINT fk__bll_balance__user_id__usr_users FOREIGN KEY (user_id) REFERENCES public.usr_users(user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS public.bll_balance;
-- +goose StatementEnd
