create table if not exists payment (
    id UUID primary key not null,
    user_id UUID not null,
    external_id text default null,
    amount_penny integer not null,
    "type" text not null,
    description text not null,
     redirect_url text default null,
     status text not null,
     created_at TIMESTAMPTZ not null default NOW(),
     updated_at TIMESTAMPTZ not null default NOW()
)