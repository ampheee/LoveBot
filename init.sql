create database tacyproject;
grant all privileges on database tacyproject to postgres;
\c tacyproject;
create table if not exists images (
    id serial primary key,
    image bytea,
    seen bool default false,
    fromuser bool default false
);
create table if not exists compliments(
    id serial primary key,
    compliment text,
    seen bool default false
);
create table if not exists thougths(
    id serial primary key,
    thought text
);
create table if not exists users(
    userid serial primary key,
    role int
);

