# dcpfp

Retrieve a user's Discord profile picture and open it in the browser. Specify
a user ID as a target with `-t`, or a tag with `-g` in order to search your
friends (and blocked) list. 

Flags:
- `-T`: Set an authentication token. 
- `-t`: Set a target user ID. This works globally.
- `-g`: Set a target Discord tag (e.g. `Kaieteuria#9522`). This will only
  work on users you're friends with and users you've blocked. It's
  case-sensitive, but it also always returns the closest match to the query,
  so it never comes up empty. We hope. Note that this works whether the
  target has friended you in turn or not.
- `-me`: Set yourself as the target user.
- `-p`:  "Print-only." This flag causes the program to print the user's profile picture URL, as opposed to attempting to open a browser.

The simplest valid command line is something like `dcpfp -T <token> -g
Friend#1234`, or just `dcpfp -g Friend#1234` if you've already exported
`DCPFP_TOKEN`.
