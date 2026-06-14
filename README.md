<h1 align="center">Gator</h1>
<h3 align="center">Blog Aggregator</h1>

## Installation

- You will need Postgres and Go installed to run this program.
- Use the command `go install github.com/jacksterdealeo/gator@latest`

## Setup

- Configure Postgres:
  - Install Postgres with your package manager.
  - Create a database for gator, you can use the command `CREATE DATABASE gator;` inside PSQL.
    - On linux, you will need to set up a password afterward: `ALTER USER postgres PASSWORD 'PASSWORD123';`
- Create a config file in your home directory. It'll need to be a JSON file called ".gatorconfig.json" in this format:
```json
{ "db_url": "postgres://PCUSERNAME:PASSWORD123@localhost/gator?sslmode=disable" }
```
- Add yourself as a user with `gator register <your-name>`
- Add feeds with `gator addfeed <feed-name> <feed-URL>`
- Follow feeds after they have been added with `gator follow <feed-URL>`
- Run `gator agg <time>` to start start aggregation. The time format is specified [here](https://pkg.go.dev/time#ParseDuration), but I would personally recommend `10m` for 10 minutes. Very short durations can strain the servers.

## Commands

-	login:    {Username}: This sets the user to {Username}.
-	register: {Username}: This registers the user {Username} and logs that user in.
-	reset:    Prompts you if you want to delete all users.
-	users:    This lists all users.

-	agg:     {Duration}: This starts the aggregation process. Leave this running in the background to keep your system up to date on the latest posts.
-	addfeed: {Name} {URL}: This adds the feed, and makes it avalible to follow.
-	feeds:   This lists all feeds for all users.
-	follow:  {URL}: This follows a feed. Feed must be added first.

-	following: Lists all the feeds the user is following.
-	unfollow:  {URL}: This removes the feed the user is following.
-	browse:    {PostCount (Default is 2)}: Shows posts in order of most to least recent.

-	help: {Command (Optional)}: This displays all the commands and their help message.

- All the commands used are explained with the `gator help` command.
