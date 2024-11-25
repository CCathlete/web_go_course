-- To automatically copy into the psql container and run use:
-- docker cp ~/Repos/web_go_course/psql_commands.pgsql web_go_course_db_1:/; docker-compose exec db psql -U ccat -d webgo -f psql_commands.pgsql

SELECT id, email FROM users;