create database tech_dojo;
use tech_dojo;

create table users (
    id int not null auto_increment,
    name varchar(20) not null,
    primary key (id)
);

create table characters (
    id int not null auto_increment,
    name varchar(20),
    probability int,
    primary key (id)
);

insert into characters values (1,'yoshiharu habu', 10); 
insert into characters values (2,'sota hujii', 5);
insert into characters values (3,'amahiko sato', 30); 
insert into characters values (4,'akira watanabe', 30); 
insert into characters values (5,'toshiyuki moriuchi', 30); 
insert into characters values (6,'takuya nagase', 30); 
insert into characters values (7,'masayuki tyoshima', 30); 
insert into characters values (8,'shintaro saito', 30); 
insert into characters values (9,'yuki sasaki', 50); 

create table characters_possession (
    usercharacterid int not null auto_increment,
    userid int,
    characterid int,
    primary key (usercharacterid)
);

