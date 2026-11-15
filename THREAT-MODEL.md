# Threat model: Sicheong 

## What is this?

A threat model answers three questions:

1. What do we want to protect?
2. Who might try to attack it, and how?
3. What have we done about it, and what is still open?

Sicheong is a small website where authors log in and post short links with a title and a comment. Visitors can read posts. Anyone can see the code on GitHub.

## What we protect

| Thing | Why it matters |
|---|---|
| Author passwords | Nobody should be able to take over someone else's account. People often reuse passwords on other sites, so a leaked password is also a risk elsewhere. |
| Login sessions | Whoever holds a session can post as that author |
| Post content and authorship | Posts should only be changed by their author |
| The site being available | It should not crash or go down |
| Deploy access | Whoever can deploy can change the live site |

## Who might attack it

- **Strangers on the internet**, who can see the site and try to break in.
- **Someone who has a valid login**, who may try to change other people's posts.
- **Someone who gets a copy of the code or database**, for example from GitHub or a backup.
- **Anyone who can push code** to the repository or trigger a deploy.

## The threats

Status: **Fixed** means it is handled in the code. **Partly fixed** means some protection exists. **Open** means it still needs work.

| # | What could go wrong | Plain-language explanation | Status | Next step |
|---|---|---|---|---|
| 1 | Default passwords | The first start creates three accounts with the password `abc`, and the README says so. Anyone can log in as them. | **Open, critical** | Remove the default passwords. Check the live site and change any still using `abc`. |
| 2 | Passwords stored as plain text | If someone gets the database, they can read every password. | **Open, high** | Store scrambled (hashed) passwords. The code for this already exists but is not used yet. |
| 3 | Password guessing | An attacker can try many passwords until one works. | **Open, high** | Block repeated failed logins for a short time. |
| 4 | Sessions never expire | A copied login cookie keeps working until the password changes or the server restarts. | **Partly fixed** | Expire sessions after a set time and clear them on logout. |
| 5 | Editing other people's posts | Any logged-in author could change or delete someone else's post. | **Fixed** | Add tests to keep it that way. |
| 6 | Hidden commands in form inputs (SQL injection) | An attacker types database commands into a form to read or delete data. | **Fixed** | Inputs are passed to the database as separate values, never as commands. Only two sort options are built into the query, and only from a fixed list. |
| 7 | Script injection (XSS) | An attacker posts text that runs code in other visitors' browsers. | **Fixed** | The page template escapes the text it shows. |
| 8 | Fake requests from other sites | Another website tricks your browser into posting on Sicheong. | **Mostly fixed** | Login cookies are only sent with requests from Sicheong itself. |
| 9 | The site crashes or goes down | Some errors stop the whole server, and a few login mistakes cause crashes. Some database bugs break the author filter. | **Open, medium** | Stop the fatal errors and fix the author filter. Bad input should show an error, not crash. |
| 10 | Unsafe deploys | Fly and Render deploy on every merge to `master` without waiting for the tests to pass. The GKE workflow deploys on pushes to any branch. | **Open, high** | Only deploy after tests pass. Limit GKE deploys to `master`. |
| 11 | Outdated libraries | Libraries can have known security holes. GitHub already warns about known problems in the libraries used. | **Partly fixed** | Confirm automatic update pull requests are on (add `.github/dependabot.yml` if not). Add a vulnerability check to CI. |
| 12 | Unencrypted traffic | Someone on the same network reads passwords or cookies. | **Fixed on fly.io** | The fly.io site only works over HTTPS. Render and GKE still need checking. |

## Other things to tidy up

These are smaller hardening tasks, not separate threats:

- Add basic browser security headers (they stop the site being shown inside another site).
- Stop the app running as the root user inside its container.
- Set limits on how large a post can be.
- Log logins, edits and deletes, so there is a record of what happened.

## Security measures already in place

1. Database inputs are always passed as values, never as commands.
2. Only the post's author can edit or delete it.
3. Author names are taken from the login, not from the form, so posts cannot be given to someone else.
4. Login cookies are encrypted, and they are marked so scripts cannot read them and other sites cannot send them.
5. Templates escape user text before showing it.
6. Automated checks run on every change: formatting, tests, a linter, and a security scanner (Gosec).
7. GitHub watches the project: alerts for vulnerable libraries, code scanning for common coding flaws, secret scanning, and private vulnerability reporting.
8. The live site on fly.io only works over TLS.
9. The database file is kept out of Git.

## Accepted risks

- The code is public, so attackers can see how it works. Nothing in it is a secret, so this is accepted.
- Deploy credentials are long-lived. Accepted until they are rotated.

## How we will check it works

1. Try to log in with `abc`. It should fail once threat 1 is fixed.
2. Add tests for: editing someone else's post, and the author filter.
3. Check that CI blocks a deploy when tests fail.
4. Try too many wrong passwords. The site should slow down or block you.

## Order of work

1. Remove default passwords (threat 1).
2. Hash passwords and limit login attempts (threats 2 and 3).
3. Fix the crashes and the author filter (threat 9).
4. Block unsafe deploys (threat 10).
5. Expire sessions (threat 4).
6. Everything else.
