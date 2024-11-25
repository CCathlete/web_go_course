-- To automatically copy into the psql container and run use:
-- docker cp ~/Repos/web_go_course/psql_commands.pgsql web_go_course_db_1:/; docker-compose exec db psql -U ccat -d webgo -f psql_commands.pgsql

-- insert into users (first_name, last_name, age, email) values 
-- ('Ronny', 'Dio', 70, 'foolfool@example.com'),
-- ('Tim', 'Minchin', 55, 'onlyaginger@example.com');
update users set first_name = 'Ronnie' where first_name = 'Ronny' and last_name = 'Dio';
select * from users;