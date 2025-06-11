# Welcome to the goGator - blog aggregator written in Golang!

It allows you to aggregate posts from websites' rss feeds into PostgreSQL database and then read them at your own convenience.

To run this program you need:
1. [go](https://webinstall.dev/golang/);
2. [postgreSQL database](https://webinstall.dev/postgres/);
The database should be migrated to the version 005_posts. I used [Goose command line tool](https://github.com/pressly/goose#install). To run a migration with it, run "goose postgres *connection_string* up" in folder "sql/schema "
3. JSON config file in your home directory named "~/.gatorconfig.json" with the following content:
```
{
  "db_url": "database_connection_string", 
  "current_user_name": "username"
}
```
Don't forget to substitute *database_connection_string* with your database connection string!

To install goGator do "go install" in the goGator folder. Now the goGator can be run!
GoGator accepts a few commandas and mandatory/optional arguments to run (goGator **command** arg1 arg2). Here is the list of commands you can use:
1. **register** username - to register user with the username;
2. **login** name - to switch the current user to the user with the name (if it exist);
3. **reset** - resets the databases;
4. **users** - shows the list of registered users and the current user;
5. **addfeed** "name" "url" - registers the rss-feed at url with the name;
6. **feeds** - shows the list of registered feeds;
7. **follow** "url" - to follow the registered feed by url with the current user;
8. **unfollow** "url" - to unfollow the registered feed by url with the current user;
9. **following** - shows the list of the feeds the current user is following;
10. **agg** time_delay - fetches the rss feeds, one at a time, from the registered url (starting from oldest) and saves the posts into the database. The fetch happens every time delay. Time format: 1m (for 1 mintute), 1h (for 1 hour). Please do not run it shorter than 1 minute.
This command runs an infinite loop. For this reason, run the second instance of the goGator for the posts aggregation or aggregate and then use Ctrl+C to stop it and continue using other commands.
11. **browse** number (optional) - browses number of posts (or 2 if not chosen) for the current user, showing newest first, according to the feeds the user is following. 

Have a pleasant goGator use! If you have any issues with running this program or have suggestions - don't hesitate to let me know!


