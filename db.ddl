CREATE TABLE public.balance
(
    user_id integer,
    amount double precision
)

CREATE TABLE public.transactions
(
    user_id integer,
    comment character varying COLLATE pg_catalog."default",
    amount double precision,
    date timestamp with time zone
)