# Threat model: Sicheong

Sicheong is a small website where authors log in and post short links with a title and a comment. Visitors can read posts. The code is public on GitHub.

This document answers three questions: what we protect, who might attack it, and what is still open. The measures we have put in place are listed in [SECURITY-MEASURES.md](SECURITY-MEASURES.md).

## What we protect


| Asset                | Why it matters                                                                    |
| -------------------- | --------------------------------------------------------------------------------- |
| Author passwords     | Nobody should take over someone else's account. People reuse passwords elsewhere. |
| Login sessions       | Whoever holds a session can post as that author.                                  |
| Posts and authorship | Only the author should change their posts.                                        |
| Availability         | The site should not crash or go down.                                             |
| Deploy access        | Whoever can deploy can change the live site.                                      |


## Who might attack it

- **Strangers on the internet**, who can see the site and try to break in.
- **Logged-in authors**, who may try to change other people's posts.
- **Anyone with a copy of the code or database**, for example from GitHub or a backup.
- **Anyone who can push code** or trigger a deploy.

## The threats

**Status:** *Fixed* means handled in the code. *Partly fixed* means some protection exists. *Open* means work is still needed.


| #   | What could go wrong                                                                                | Status              | Next step                                                                    |
| --- | -------------------------------------------------------------------------------------------------- | ------------------- | ---------------------------------------------------------------------------- |
| 1   | The first start creates accounts with the password `abc`, and the README says so.                  | **Open, critical**  | Remove the default passwords. Check the live site.                           |
| 2   | Passwords are stored as plain text.                                                                | **Open, high**      | Store hashed passwords. Helper code exists but is not used yet.              |
| 3   | Attackers can guess passwords by trying many.                                                      | **Open, high**      | Block repeated failed logins for a short time.                               |
| 4   | A copied session cookie works until the server restarts.                                           | **Partly fixed**    | Expire sessions and clear them on logout.                                    |
| 5   | Logged-in authors can edit or delete other people's posts.                                         | **Fixed**           | Add tests to keep it that way.                                               |
| 6   | Attackers type database commands into forms (SQL injection).                                       | **Fixed**           | Values are passed separately; sort options come from a fixed list.           |
| 7   | Attackers post script that runs in other visitors' browsers (XSS).                                 | **Fixed**           | Templates escape user text.                                                  |
| 8   | Other websites trick your browser into making requests to Sicheong.                                | **Mostly fixed**    | Cookies are only sent by requests from Sicheong.                             |
| 9   | Some errors crash the server, and a database bug breaks the author filter.                         | **Open, medium**    | Stop fatal errors and fix the author filter. Bad input should show an error. |
| 10  | Fly and Render deploy on merge to `master` without waiting for tests. GKE deploys from any branch. | **Open, high**      | Deploy only after tests pass. Limit GKE to `master`.                         |
| 11  | Libraries may have known security holes.                                                           | **Partly fixed**    | Keep automatic updates on. Add a vulnerability check to CI.                  |
| 12  | Someone on the same network reads passwords or cookies.                                            | **Fixed on fly.io** | Check Render and GKE use HTTPS too.                                          |


## Smaller tasks

These are hardening steps, not separate threats:

- Add browser security headers, so other sites cannot embed Sicheong.
- Stop the app running as root inside its container.
- Limit how large a post can be.
- Log logins, edits and deletes.

## Accepted risks

- The code is public. Nothing in it is secret, so this is accepted.
- Deploy credentials are long-lived. Accepted until they are rotated.

## How we will check it works

1. Try to log in with `abc`. It should fail once threat 1 is fixed.
2. Test that editing someone else's post is refused, and that the author filter works.
3. Check that CI blocks a deploy when tests fail.
4. Try many wrong passwords. The site should slow down or block you.

## Order of work

1. Remove default passwords (threat 1).
2. Hash passwords and limit login attempts (threats 2 and 3).
3. Fix crashes and the author filter (threat 9).
4. Block unsafe deploys (threat 10).
5. Expire sessions (threat 4).
6. Everything else.
