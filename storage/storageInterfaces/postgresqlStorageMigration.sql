-- Project is kinda small so just going to dump an SQL file in to initialize a database in a usable state.
-- In practice a migration system would be better to apply changes according to code version.
-- TODO: Are there any packages in go that could satisfy this need?

CREATE DATABASE go_crawler;
CREATE SCHEMA public;

CREATE TABLE frontier(
    url           VARCHAR PRIMARY KEY     NOT NULL,
    status        CHAR(50)
);

CREATE INDEX frontier_status_index
    ON public.frontier (status)
;