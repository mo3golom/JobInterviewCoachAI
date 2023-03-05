create table if not exists "user" (
    id UUID primary key not null,
    created_at TIMESTAMPTZ not null default NOW(),
    updated_at TIMESTAMPTZ not null default NOW()
)