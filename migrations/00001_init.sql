-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
create table if not exists counter
(
    name text not null
        constraint counter_pk
            primary key,
    value bigint not null
);

create table if not exists gauge
(
    name text not null
        constraint gauge_pk
            primary key,
    value double precision not null
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
drop table if exists counter;
drop table if exists gauge;
-- +goose StatementEnd
