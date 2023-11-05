create table if not exists subscription_free_attempt (
     user_id UUID primary key not null,
     attempts integer default 10,
     created_at TIMESTAMPTZ not null default NOW(),
     updated_at TIMESTAMPTZ not null default NOW()
)