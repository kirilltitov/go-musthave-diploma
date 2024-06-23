create type public.order_status as enum ('NEW', 'PROCESSING', 'INVALID', 'PROCESSED');

create table public."user"
(
    id         uuid        not null constraint user_pk primary key,
    login      varchar     not null constraint login_key unique,
    password   varchar     not null,
    created_at timestamptz not null
);

create table public."order"
(
    id           uuid           not null constraint order_pk primary key,
    order_number varchar        not null constraint order_number_key unique,
    user_id      uuid           not null,
    status       order_status   not null,
    amount       numeric(10, 2) not null,
    created_at   timestamptz    not null,
    updated_at   timestamptz
);

create index order_created_at_index on public."order" (created_at);
create index order_user_id_index on public."order" (user_id);

create table public.account
(
    user_id           uuid                     not null constraint account_pk primary key,
    current_balance   numeric(10, 2) default 0 not null,
    withdrawn_balance numeric(10, 2) default 0 not null
);

create table public.withdrawal
(
    id           uuid           not null constraint withdrawal_pk primary key,
    user_id      uuid           not null,
    order_number varchar        not null constraint withdrawal_order_number_ukey unique,
    amount       numeric(10, 2) not null,
    created_at   timestamptz    not null
);

---- create above / drop below ----

drop table public."user";
drop table public."order";
drop table public.account;
drop table public.withdrawal;
drop type public.order_status;
