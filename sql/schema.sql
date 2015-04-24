create table user (
	`id` int not null auto_increment,
	`name` varchar(120) not null,
	`email` varchar(120) not null unique,
	`profile` varchar(120),
	`picture` varchar(120),
	`isAdmin` bool default 0,
	`firstLogin` datetime not null,
	`lastLogin` datetime,
	primary key (`id`)
)
default charset = utf8;