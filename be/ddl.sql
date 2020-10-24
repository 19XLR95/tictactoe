create table game_difficulties(
	id int auto_increment primary key,
	name varchar(64) not null,
	created_at timestamp not null default current_timestamp,
	updated_at timestamp not null default current_timestamp
);
insert into game_difficulties(name) values('Easy');
insert into game_difficulties(name) values('Impossible');

create table games(
	id int auto_increment primary key,
	game_key varchar(12) unique not null,
	username varchar(64) not null,
	difficulty int references game_difficulties(id),
	game_state varchar(25) not null default '[[0,0,0],[0,0,0],[0,0,0]]',
	is_cpu_turn bool not null default false,
	cpu_turn_started_at timestamp null default null,
	finished_at timestamp null default null,
	who_win tinyint(1) default null,
	created_at timestamp not null default current_timestamp,
	updated_at timestamp not null default current_timestamp,
	foreign key (difficulty) references game_difficulties(id)
);
