create table if not exists opensaga.saga_call (
    id              uuid primary key,
    -- todo add request-id
    idempotency_key uuid unique,
    saga_id         uuid  not null references opensaga.saga(id),
    content         jsonb not null
    -- todo add status
);
