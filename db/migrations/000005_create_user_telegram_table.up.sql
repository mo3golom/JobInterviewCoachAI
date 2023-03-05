create table if not exists user_telegram (
    id bigint generated by default as identity primary key not null,
    user_id UUID not null,
    telegram_id integer not null,
    created_at TIMESTAMPTZ not null default NOW(),
    updated_at TIMESTAMPTZ not null default NOW(),

    foreign key (user_id) references "user" (id),
    constraint unique_user_id unique(user_id)
)