create table if not exists users (
	id int primary key,
	username varchar(256) not null,	
	email varchar(256) not null unique,
	firstname varchar(256) not null,
	lastname varchar(256) not null,
	usertype varchar(256) not null,
	password varchar(256) not null,
);

create table if not exists submissions (
	id int primary key,
	filetype varchar(256) not null,
	location varchar(256) not null,
	email varchar(256) not null,
	firstname varchar(256) not null,
	lastname varchar(256) not null
);
