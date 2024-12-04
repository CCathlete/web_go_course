create table sessions (
  id serial primary key,
  user_id bigint unique,
  token_hash text unique not null,
  foreign key (user_id) references users(id)
);