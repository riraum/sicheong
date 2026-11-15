# 시청 (Sicheong) / City hall
> Go project to learn ([CRUD](https://en.wikipedia.org/wiki/Create,_read,_update_and_delete)).

# Background
Before, focused on the basics in Go and wrote a Hackernews CLI.
Had also done basic, static sites, with HTML, CSS and JavaScript.
Now, a simple CRUD project, in Go.

Wrote some basic CSS rules and for the rest used Pico CSS framework.
Instead a lot of time was used to learn about authentication and deployment.

Authentication should has quite some room for improvement, but deployment feels alright.

Researched usage of an ORM or SQL query builder such as GORM, ent, sqlc, sqlx, sqirrel and the likes.
Experiemented with GORM a bit, but reverted to pure SQL queries as it seemed to both make it easier and harder.

Also considered using PostgreSQL, instead of SQLite, as that's what "real projects" do. Instead decided to double down on SQLite and learn how to make it work.
Deployed the database on 2 persistent volumes on fly.io. Not yet done, but could be periodically saved to S3like storage.
Some posts that helped make that decision:

https://fly.io/blog/all-in-on-sqlite-litestream/

https://fly.io/blog/litestream-revamped/

https://kentcdodds.com/blog/i-migrated-from-a-postgres-cluster-to-distributed-sqlite-with-litefs

Other sources that helped apart from the Go documentation, are linked in the respective PRs.

# External packages/assets used
- Mattn Go SQLite driver: github.com/mattn/go-sqlite3
- Pico CSS framework: https://picocss.com/


# Database setup

- The app uses a SQLite file at `litefs/sq.db`. The folder is created on start if it doesn't exist.
- On first start, the app creates the tables and seeds 3 authors (Alpha, Bravo, Charlie, password `abc`) and 3 posts.
- To reset the data, stop the app, delete `litefs/sq.db`, and start it again.
- Tables: `authors`, `posts`, `votes`. Indexes: `posts_author_idx`, `votes_post_idx`. See the [ER model](ER-MODEL.md).
- Use cases and their status: [Use cases](USE-CASES.md).

## Data

A fresh database holds 3 authors, 3 posts, and 0 votes. The fixed seed values are filled in `Fill` in `db/db.go`. The votes table starts empty, so the many-to-many relationship is shown by the schema, not by seed data.


# Security
- See [Threat Model](THREAT-MODEL.md)
- And [Security Measures](SECURITY-MEASURES.md)

# Development
Requires
  - Go >=1.25

Then
- Clone repo
- Run it with `go run main.go`, then open http://localhost:8081
- Run the tests: `go test ./...`

For deployment on fly.io
- use their CLI
- create 2 persistent volumes for a persistent sqlite database, mount them as `litefs` folder
- read their guides for LiteFS and deployment in general
- can use the repos `fly-deploy.yml`, `fly.toml`, `litefs.yml` and LiteFS specific `Dockerfile-litefs` to see what you have to adjust


# Deployed:
- On fly.io
https://sicheong.fly.dev/
- On render.com
https://si-cheong.onrender.com/
- On Google Cloud Computing (GCP)
Deployed in the past on GCP's Google Kubernetes Engine (GKE)
