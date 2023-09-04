
create database if not exists playhouse;
use playhouse;

create table if not exists session(
    id serial primary key,
    is_available bool not null default false,
    due_at timestamp not null,
    created_at timestamp not null default current_timestamp(),
    user_id int references "user"(id) on delete cascade on update no action
);

create table if not exists video(
    id serial primary key,
    name string not null check (length(name) > 0),
    type string not null check (length(type) > 0),
    size int not null check(size > 0),
    url_to_stream string,
    pending_chunks int4 not null default 0,
    is_deleted bool not null default false,
    is_transcoded bool not null default false,
    created_at timestamp not null default current_timestamp(),
    uploaded_at timestamp,
    session_id int references session(id) on delete set null on update no action
);

create table if not exists chunk(
        code int not null check( code >= 0),
        size int not null default 0 check ( size >= 0),
        is_uploaded bool not null default false,
        created_at timestamp not null default current_timestamp(),
        uploaded_at timestamp,
        video_id int not null references video(id) on delete cascade on update no action,
        session_id int references session(id) on delete set null on update no action,
        primary key (video_id, code)
);

create table if not exists "user" (
       id           serial primary key,
       created_at   timestamp not null default current_timestamp(),
       email string not null check (length(email) > 0)
);