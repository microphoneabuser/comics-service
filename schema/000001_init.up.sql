CREATE TABLE users
(
    id serial not null unique,
    name varchar(255) not null,
    login varchar(255) not null unique,
    password_hash varchar(255) not null
);

CREATE TABLE comics
(
    id serial not null unique,
    title varchar(255) not null,
    user_id int references users (id) on delete cascade not null,
    date date not null,
    img varchar(255) not null,
    description varchar(255)
);