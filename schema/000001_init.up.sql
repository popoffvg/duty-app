CREATE TABLE users
(
    id serial not null unique CONSTRAINT users_pk PRIMARY KEY,
    name varchar(64) not null,
    username varchar(64) not null unique,
    password varchar(512) not null
);

CREATE TABLE teams
(
    id serial not null unique CONSTRAINT teams_pk PRIMARY KEY,
    user_id int references users (id) on delete cascade not null,
    title varchar(255) not null
);

CREATE TABLE teammates
(
    id serial not null unique CONSTRAINT teammates_pk PRIMARY KEY,
    team_id int references teams (id) on delete cascade not null,
    name varchar(255) not null,
    duty_readiness bool not null,
    duties int not null
);

CREATE TABLE duties
(
    id serial not null unique CONSTRAINT duties_pk PRIMARY KEY,
    team_id int references teams (id) on delete cascade not null,
    teammate_id int references teammates (id) on delete cascade not null,
    is_daily bool not null,
    date date not null,
    CONSTRAINT duties_unique_constraint UNIQUE (team_id, teammate_id, is_daily, date)
);


