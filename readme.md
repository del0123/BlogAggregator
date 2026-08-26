# Gator (Blog Aggregator CLI)

Gator is a multi-user, terminal-based RSS feed aggregator built in Go and powered by PostgreSQL. It allows users to manage feed subscriptions, scrape post updates continuously in the background, and read articles directly from the command line.

---

## Prerequisites

Before installing and running Gator, ensure you have the following software installed on your machine:

* **[Go](https://go.dev/doc/install)** (v1.20 or later)
* **[PostgreSQL](https://www.postgresql.org/download/)** running locally

---

## Installation

Gator is compiled into a single static binary. You can build and install it directly to your `$GOPATH/bin` using `go install`:

```bash
go install github.com/YOUR_GITHUB_USERNAME/BlogAggregator@latest

```

> **Note:** Replace `YOUR_GITHUB_USERNAME` with your actual GitHub username. Ensure that your `$GOPATH/bin` (or `~/go/bin`) directory is included in your system's `PATH` environment variable so you can run the binary from anywhere.

Alternatively, for local development, you can build the executable directly inside the project root:

```bash
go build -o gator .

```

---

## Configuration & Database Setup

1. **Create the Database:**
Ensure your PostgreSQL service is running, and create a local database for the project (e.g., named `gator`).
2. **Run Migrations:**
Apply your database schema migrations to set up the necessary tables (`users`, `feeds`, `feed_follows`, `posts`).
3. **Set Up `.gatorconfig.json`:**
Gator expects a JSON configuration file located in your user's home directory (`~/.gatorconfig.json`). Create this file manually with your local database connection details:
```json
{
  "db_url": "postgres://postgres:postgres@localhost:5432/gator?sslmode=disable",
  "current_user_name": ""
}

```


*(Modify the database URL credentials, host, port, and database name to match your PostgreSQL setup.)*

---

## Usage & Commands

Once configured, run `gator` directly from your terminal. Below is a list of key commands:

### User Management

* **Register a new user:**
```bash
gator register <username>

```


* **Log in as an existing user:**
```bash
gator login <username>

```


* **List all users:**
```bash
gator users

```



### Feeds & Subscriptions

* **Add a new feed:**
```bash
gator addfeed "Feed Name" "https://example.com/feed.xml"

```


* **List all available feeds:**
```bash
gator feeds

```


* **Follow an existing feed:**
```bash
gator follow "https://example.com/feed.xml"

```


* **Unfollow a feed:**
```bash
gator unfollow "https://example.com/feed.xml"

```


* **List feeds followed by the logged-in user:**
```bash
gator following

```



### Aggregation & Reading

* **Start the continuous background feed scraper:**
```bash
gator agg 1m

```


*(Fetches posts every 1 minute. Accepts duration strings like `1s`, `1m`, `1h`. Press `Ctrl+C` to stop.)*
* **Browse saved posts:**
```bash
gator browse

```


*(Defaults to displaying 2 posts. Optionally pass a numerical limit, e.g., `gator browse 10`)*