create database notesapp ;

create table notes (
    id varchar(64) primary key ,
    title varchar(1024) not null  ,
    content text null
);