create table if not exists opensaga.saga_step (
    id              uuid primary key,
    saga_id         uuid
        references opensaga.saga(id),
    next_on_success uuid
        references opensaga.saga_step(id)
            initially deferred,
    next_on_failure uuid
        references opensaga.saga_step(id)
            initially deferred,
    is_initial      bool not null default false,
    name            text not null,
    endpoint        text not null,

    unique (saga_id, name)
);
