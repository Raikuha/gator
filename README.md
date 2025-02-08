# Gator (RSS Aggregator)

This project called gator is a CLI RSS aggregator that requires postgresql and Go to run

## Install

You can install gator by using `go install` on the terminal

```
$ go install github.com/Raikuha/gator@latest
```

## Getting started

You'll need a JSON config file named `.gatorconfig.json` in your home directory to store the connection to your postgresql database

```
{
    "db_url": "database_connection"
    "current_user_name": "name",
}
```

Once everything is set, you can use it through the command line with a set of simple commands

## User Commands

* Register a new user
> ./gator register [username]

* Login as an existing user
> ./gator login [username]

* List all users
> ./gator users

* Delete all users (and any related database entries)
> ./gator reset

## Feeds

* Add a feed (associated to current active user)
> ./gator addfeed [title] [url]

* List all feeds in the database
> ./gator feeds

* Follow a feed added by other users
> ./gator follow [url]

* Unfollow a feed
> ./gator unfollow [url]

* List all your current follows
> ./gator following

## Posts

* Retrieve new posts over an interval such as 1m, 1h...
> ./gator agg [interval]

* Browse through followed posts (default 2)
> ./gator browse [num]

While using the command Agg to update the database in the background, you can still open a new terminal to perform other commands