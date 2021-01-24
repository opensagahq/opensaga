create table if not exists opensaga.saga_call (
    idempotency_key uuid primary key,
    saga_id         uuid not null references opensaga.saga(id)
);
