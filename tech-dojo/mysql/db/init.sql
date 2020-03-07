create database tech_dojo;
use tech_dojo;

create table users (
    id int not null auto_increment,
    name varchar(20),
    primary key (id)
);

create table characters (
    id int not null auto_increment,
    name varchar(20),
    primary key (id)
);

create table characters_possession (
    userid int,
    characterid int
);