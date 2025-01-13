-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS public.wdr_withdrawals (
	withdrawal_id uuid DEFAULT gen_random_uuid() NOT NULL,
    user_id uuid NOT NULL,
    "order" varchar NOT NULL,
	sum numeric DEFAULT 0 NOT NULL,
	processed_at timestamp without time zone DEFAULT current_timestamp NOT NULL,
	CONSTRAINT wdr_withdrawals_pk PRIMARY KEY (withdrawal_id)
);

ALTER TABLE public.wdr_withdrawals ADD CONSTRAINT fk__wdr_withdrawals__user_id__usr_users FOREIGN KEY (user_id) REFERENCES public.usr_users(user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS public.wdr_withdrawals;
-- +goose StatementEnd
