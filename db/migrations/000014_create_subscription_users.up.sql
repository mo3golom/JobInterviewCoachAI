create table if not exists subscription_users (
     user_id UUID primary key not null,
     type text not null default 'free',
     created_at TIMESTAMPTZ not null default NOW(),
     updated_at TIMESTAMPTZ not null default NOW()
)