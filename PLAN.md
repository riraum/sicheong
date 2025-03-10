# PLAN

End date: 2025-05-23

> [!NOTE]
> Write issues and milestones for all those when you can.
> 
> E.g. https://github.com/riraum/project451/milestone/1
>
> (Those should not be really acted on, all the work should be done in the final repo)

- Milestone 00 - Due 2025-03-12 - Project setup
  - Find a name
  - Create a repository
  - Setup CI
  - Write a `script/run` script that runs the code
  - Write a dummy `main.go`
  - Write a dummy `main_test.go`
  - Write in issues at least milestone 01 and 02

- Milestone 01 - Due 2025-03-19 - Basic website
  - Write a webserver that respond to `GET /` with `200 "OK"`
  - Write a router to handle responses to
    - `GET /` => `200 "OK"`
    - `GET /api/v0/posts` => `200 "[]"`
    - `POST /api/v0/posts` => `201`
    - Other paths should result in => `404`
  - Write tests for each route
  - Resources:
    - https://pkg.go.dev/net/http#hdr-Servers
    - https://pkg.go.dev/net/http#Server
    - https://pkg.go.dev/net/http#ServeMux
    - https://http.cat/

- Milestone 02 - Due 2025-03-26 - Basic DB
  - Create an SQLite db
  - Define a simple schema with a single table
    - `Post`: `id, title, link`
  - Write go code to initialize the database with the table
  - Write go code to write random data into the db  
  - Write go code to read all posts from the db
  - Write tests against a test db
  - Resources:
    - https://go.dev/doc/database/
    - https://pkg.go.dev/database/sql
    - https://go.dev/wiki/SQLDrivers
    - https://github.com/mattn/go-sqlite3 (most popular lib)
    - https://github.com/mattn/go-sqlite3/blob/master/_example/simple/simple.go

- Milestone 03 - Due 2025-04-02 - Basic frontend rendering

- Milestone 04 - Due 2025-04-09 - 

- Milestone 05 - Due 2025-04-16 -

- Milestone 06 - Due 2025-04-23 -

- Milestone 07 - Due 2025-04-30 -

- Milestone 08 - Due 2025-05-07 -

- Milestone 09 - Due 2025-04-14 -

- Milestone 10 - Due 2025-04-21 -

- Milestone 11 - Due 2025-04-25 - Final review and cleanup
