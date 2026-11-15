# Use cases: Si-cheong

Status key: 

**Implemented** means it works end to end in the app. 

**Partial** means only some of it exists. 

**Not implemented** means no code yet.

| # | Use case | Status | Notes |
| --- | --- | --- | --- |
| 1 | View the list of posts | Implemented | Home page, `GET /` |
| 2 | Sort posts by date or title, ascending or descending | Implemented | Query parameters `sort` and `direction` |
| 3 | View a single post | Implemented | `GET /post/{id}` |
| 4 | Log in as an author | Implemented | `GET` and `POST /login` |
| 5 | Log out | Implemented | `GET /logout` |
| 6 | Create a post while logged in | Implemented | `POST /post` |
| 7 | Edit your own post | Implemented | `POST /post/{id}`. Only the post's author can edit it. |
| 8 | Delete your own post | Implemented | `POST /post/delete/{id}`. Only the post's author can delete it. |
| 9 | Filter posts by author | Partial | Parameter `author` is parsed, but the query fails when it is set. See known issues. |
| 10 | Use the JSON API to list, view, create, edit, and delete posts | Implemented | `/api/v0/posts`, `/api/v0/post`, `/api/v0/post/{id}`. Edit and delete require the author. |
| 11 | Log in and out through the JSON API | Implemented | `POST /api/v0/login`, `GET /api/v0/logout` |
| 12 | Vote for a post | Partial | Stored in the database (`votes` table). No button or route in the UI yet. |
| 13 | Count the votes on a post | Partial | Database function exists. Not shown in the UI. |
| 14 | Comment on a post | Not implemented | Planned in the specs |
| 15 | Assign posts to categories | Not implemented | Planned in the specs |
| 16 | Delete any post as an admin | Not implemented | Planned in the specs |
| 17 | Create a new user account | Not implemented | Users are seeded at first start |

## Known issues

- **Use case 9:** Filtering by author panics or fails because the query runs a second time without its argument. The unfiltered list works.
- **Use case 12 and 13:** Votes can only be written and counted from code, not from the web or API.
- Passwords are stored in plain text.
