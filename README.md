# Gator

📝 A commandline rss aggre-gator, implemented in Go! 🐭

## Setup & Usage

This project requires the following dependencies:

- Go (only required for installation)
- Postgres, either installed locally or accessible over the network.
- Goose (only required once, to configure the Postgres database at installation)

### Setup and Installation

With Postgres installed, connect (either with `psql` or a tool like `PGAdmin`) and create a new database named `gator`
```SQL
CREATE DATABASE gator
```

To configure the postgres database, execute the below commands from the application source directory, providing a **valid** connection string
```bash
cd sql/schema && goose postgres <connection string> up
```

Gator uses a config file stored at `~/.gatorconfig.json`, and requires a **valid** connection string, as shown below. (Your config should match this format).
```json
{
  "db_url": "postgres://username@localhost:5432/gator?sslmode=disable"
}
```

To run the program *without* installing it, run:
```bash
go run . <command>
```
or, to permanently install gator into $GOBIN
```bash
go install
```

If gator is permanently installed, to run a command: `gator <command>`

### Available Commands
- `register <name>` - Register a user
- `login <name>` - Log in as a registered user
- `users` - Display a list of registered users
- `reset` - Remove all users, feeds & posts from the database
- `addfeed <feed name> <feed url>` - Register a new feed
- `feeds` - Display the list of registered feeds
- `follow <feed url>` - Follow the provided (registered) feed url
- `unfollow <feed url>` - Unfollow the provided (registered) feed url
- `following` - Display the list of followed feeds
- `agg <duration>` - Poll the registered feeds, at the provided interval. Duration format: 1m | 10m | 1h etc.
- `browse [number]` - display the selected number of posts, ordered most recent first. (defaults to 2).

---

This project was built as part of the Build a Blog Aggregator course on [Boot.Dev](https://www.boot.dev/courses/build-blog-aggregator-golang)
