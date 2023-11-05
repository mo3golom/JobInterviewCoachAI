create table if not exists subscription (
     user_id UUID primary key not null,
     start_at TIMESTAMPTZ not null default NOW(),
     end_at TIMESTAMPTZ not null default NOW(),
     created_at TIMESTAMPTZ not null default NOW(),
     updated_at TIMESTAMPTZ not null default NOW()
)