create database entry_task;

use entry_task;

create table user (
	id int unsigned not null primary key auto_increment,
	username varchar(255) unique not null,
	nickname varchar(255) not null,
	password varchar(255) not null,
	photo varchar(255)
);