create table user (
	`id` int not null auto_increment,
	`name` char(120) not null,
	`email` char(120) not null unique,
	`gender` char(6) default "",
	`profile` char(120) default "",
	`picture` char(120) default "",
	`isAdmin` bool default 0,
	`firstLogin` datetime not null,
	`lastLogin` datetime,
	primary key (`id`)
)
default charset = utf8;