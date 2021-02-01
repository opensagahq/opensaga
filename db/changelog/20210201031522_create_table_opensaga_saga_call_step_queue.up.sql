create table if not exists opensaga.saga_call_step_queue (
    id           bigserial   not null primary key,
    saga_step_id uuid        not null references opensaga.saga_step(id),
    saga_call_id uuid        not null references opensaga.saga_call(id),
    enqueued_at  timestamptz not null default now(),
    locked_at    timestamptz,
    locked_by    text,
    payload      jsonb       not null,

    unique (saga_step_id, saga_call_id)
);
