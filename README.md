# gator 🐊

A CLI RSS aggregator written in Go.

## Prerequisites

Before you begin, ensure you have the following installed:

*   [Go](https://golang.org/doc/install)
*   [PostgreSQL](https://www.postgresql.org/download/)

## Installation

To install gator, run the following command:

```bash
go install github.com/mattnickolaus/gator@latest
```

This will install the `gator` binary in your Go bin directory.

## Configuration

Before running gator, you need to setup a Postgres database and create a configuration file.

### Postgres Schema

Enter the `psql` shell:
- Mac: `psql postgres`
- Linux: `sudo -u postgres psql`

Create a new database called gator: 

``` sql
CREATE DATABASE gator;
```

Connect to the database:

``` sql
\c gator
```

Create a new user and grant it privileges to the gator database. You can do this with the following commands in `psql`:

```sql
CREATE USER gator_user WITH PASSWORD 'your_password';
GRANT ALL PRIVILEGES ON DATABASE gator TO gator_user;
```

This username and password will then be used in the PostgreSQL connection string.

```
"postgres://username:password@localhost:5432/gator"
```

### Configuration File

Next you need to create a configuration file at `~/.gatorconfig.json`. 

``` bash
touch ~/.gatorconfig.json 
```

This file should contain the following:

```json
{
  "db_url": "postgres://user:password@localhost:5432/gator?sslmode=disable",
  "current_user_name": ""
}
```

Replace the `db_url` with your PostgreSQL connection string except add the `?sslmode=disable` parameter. The `current_user_name` will be automatically set when you log in and can be left blank.

### Set Up Schema 

Once you have configured gator, you can run the following command to ensure the schema is set. 

```bash
gator setup
```

### Account Setup

The last step to is to register your user account by just defining your username with the following command: 

```bash
gator register <username>
```

Upon running the register command, your user will be logged and you are all setup to start aggregating some RSS feeds! 

## Commands

Run `gator help` to see a list of command or reference the follwoing list you can use with gator:

*   `gator register <username>`: Register a new user.
*   `gator login <username>`: Log in as a user.
*   `gator users`: List all registered users.

*   `gator addfeed <url>`: Add a new RSS feed.
*   `gator feeds`: List all RSS feeds.
*   `gator follow <feed_id>`: Follow a feed.
*   `gator following`: List all feeds you are following.
*   `gator unfollow <feed_id>`: Unfollow a feed.
*   `gator browse`: Browse the latest posts from your followed feeds.
*   `gator agg`: Manually trigger the aggregator to fetch new posts.

*   `gator reset`: Reset the database and run all migrations.
*   `gator setup`: Setup the database schema for gator.
*   `gator takedown`: Takes down the database schema for gator.
*   `gator help`: Provides a similar list describe the usage for each command 
