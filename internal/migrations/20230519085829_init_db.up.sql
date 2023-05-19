create table if not exists "user" (
    id serial primary key,
    username varchar(100) not null,
    password varchar(255) not null,
    note jsonb,
    created_at timestamptz not null default now(),
    updated_at timestamptz,
    deleted_at timestamptz
);