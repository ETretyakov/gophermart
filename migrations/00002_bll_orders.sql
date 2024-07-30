-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS public.bll_orders (
	order_id uuid DEFAULT gen_random_uuid() NOT NULL,
    user_id uuid NOT NULL,
	number varchar NOT NULL,
    status varchar DEFAULT 'NEW' NOT NULL,
	accrual numeric DEFAULT 0 NOT NULL,
	uploaded_at timestamp without time zone DEFAULT current_timestamp NOT NULL,
	CONSTRAINT bll_orders_pk PRIMARY KEY (order_id)
);

ALTER TABLE public.bll_orders ADD CONSTRAINT fk__bll_orders__user_id__usr_users FOREIGN KEY (user_id) REFERENCES public.usr_users(user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS public.bll_orders;
-- +goose StatementEnd
