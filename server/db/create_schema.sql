
create database if not exists playhouse;
use playhouse;

create table if not exists session(
    id serial primary key,
    is_available bool not null default false,
    due_at timestamp not null,
    created_at timestamp not null default current_timestamp()
);

create table if not exists video(
    id serial primary key,
    name string not null check (length(name) > 0),
    type string not null check (length(name) > 0),
    size int not null check(size > 0),
    is_uploaded bool not null default false,
    is_deleted bool not null default false,
    created_at timestamp not null default current_timestamp(),
    uploaded_at timestamp,
    session_id int references session(id) on delete set null on update no action
);

create table if not exists chunk(
        time_code int not null check( time_code > 0),
        size int not null check(size > 0),
        is_uploaded bool not null default false,
        created_at timestamp not null default current_timestamp(),
        uploaded_at timestamp,
        video_id int not null references video(id) on delete cascade on update no action,
        session_id int references session(id) on delete set null on update no action,
        primary key (video_id, time_code)
);
