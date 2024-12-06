-- create table sessions (
--   id serial primary key,
--   user_id bigint unique,
--   token_hash text unique not null,
--   foreign key (user_id) references users(id)
--   on delete cascade
-- );
-- select
--   email,
--   sessions.token_hash as token_hash
-- from
--   users
--   inner join sessions on users.id = sessions.user_id
-- ;
-- select
--   email,
--   sessions.token_hash as token_hash
-- from
--   users
--   left join sessions on users.id = sessions.user_id
-- ;
-- select
--   email,
--   sessions.token_hash as token_hash
-- from
--   users
--   right join sessions on users.id = sessions.user_id
-- ;
-- select
--   email,
--   sessions.token_hash as token_hash
-- from
--   users
--   full outer join sessions on users.id = sessions.user_id
-- ;
-- select
--   users.id,
--   users.email,
--   users.password_hash
-- from
--   sessions
--   join users on users.id = sessions.user_id
-- where
--   sessions.token_hash = $1
-- ;
insert into
  sessions (user_id, token_hash)
values
  ($1, $2) on conflict (user_id) do
update
set
  token_hash = $2 returning id
;