create table if not exists counter
(
    name text not null
    constraint counter_pk
    primary key,
    value bigint not null
);

alter table counter owner to postgres;

create unique index if not exists counter_name_uindex
    on counter (name);

create table if not exists gauge
(
    name text not null
    constraint gauge_pk
    primary key,
    value double precision not null
);

alter table gauge owner to postgres;

create unique index if not exists gauge_name_uindex
    on gauge (name);

