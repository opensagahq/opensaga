create table if not exists opensaga.saga (
    id   uuid primary key,
    name text not null unique
);
