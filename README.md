# bitbucket-cli

A [Bitbucket Enterprise](https://bitbucket.org/product/enterprise) CLI.

```
Usage: bitbucket-cli [--debug] [--username USERNAME] [--password PASSWORD] [--access-token ACCESS-TOKEN] [--url URL] [--config CONFIG] <command> [<args>]

Options:
  --debug, -D
  --username USERNAME, -u USERNAME
  --password PASSWORD, -p PASSWORD
  --access-token ACCESS-TOKEN, -t ACCESS-TOKEN
                         A Personal Access Token
  --url URL, -u URL      URL to the REST API of Bitbucket, e.g: https://git.example.com/rest
  --config CONFIG, -c CONFIG
  --help, -h             display this help and exit

Commands:
  project
  repo
  pr
```

# Docker container

A docker container for this project can be obtained [here](https://github.com/swisscom/bitbucket-cli/pkgs/container/bitbucket-cli).

## Project

### List

Lists the repositories in a project

```
$ export BITBUCKET_USERNAME="my-bitbucket-username"
$ read -s BITBUCKET_PASSWORD # Type your password and then press ENTER
$ export BITBUCKET_PASSWORD
$ bitbucket-cli --url https://your-bitbucket-hostname/rest project list -k PRJKEY

project-1       https://your-bitbucket-hostname/scm/prjkey/project-1.git
project-2       https://your-bitbucket-hostname/scm/prjkey/project-2.git
project-3       https://your-bitbucket-hostname/scm/prjkey/project-3.git

```

### Clone

Clones all the repositories in a project:

```
$ export BITBUCKET_USERNAME="my-bitbucket-username"
$ read -r -s BITBUCKET_PASSWORD # Type your password and then press ENTER
$ export BITBUCKET_PASSWORD
$ bitbucket-cli --url https://your-bitbucket-hostname/rest project clone -k PRJKEY -o /tmp/test/
    
    head: 987a5d8c25d8adb5ba013cf1cb88cd56a189241e5048b9702f319fb6e641cf81 refs/heads/master
    head: df2b794192904e6a9265975f33510eebe680177013e86fd7002850f45389ad34 refs/heads/master
    head: 2cf20bee2c59c3b8cae6ec0820a1353ff0ca2adeecdb84ba773845cff91ab121 refs/heads/master

$ ls -la /tmp/test 
total 0
drwxr-xr-x 13 dvitali dvitali  260 Jul 21 18:09 .
drwxrwxrwt 29 root    root    1400 Jul 21 18:11 ..
drwxr-xr-x  3 dvitali dvitali  120 Jul 21 18:09 project-1
drwxr-xr-x  4 dvitali dvitali  140 Jul 21 18:09 project-2
drwxr-xr-x  3 dvitali dvitali  100 Jul 21 18:09 project-3
```


## Repo

Most subcommands need to know which repository to operate on.  You can supply
this explicitly:

- `-k KEY` / `--key KEY`: the project key (e.g. `TOOL`).  Can also be set
  through `BITBUCKET_PROJECT`.
- `-n NAME` / `--name NAME`: the slug of an existing repository (e.g.
  `my-repo`).  Can also be set through `BITBUCKET_REPO`.

Or, if you run the command from inside a cloned Bitbucket repository, both are
detected automatically from the `origin` remote URL and can be omitted.

The project key is always needed.  The slug is needed by every subcommand
except `create`, which names the new repository with its own `--display-name`
option instead and ignores `--name`.  Both flags belong to `repo` itself and
are accepted before or after the subcommand name, but a subcommand's own
`--help` does not list them.  A subcommand that needs the slug and cannot get
one stops with an error before doing anything.

### Get

Shows a single repository.  Requires `REPO_READ` permission.

```
$ bitbucket-cli repo --key KEY --name bitbucket-playground get
Name:          bitbucket-playground
Slug:          bitbucket-playground
ID:            42
Project:       KEY
State:         AVAILABLE
Public:        false
Forkable:      true
Description:   A playground repository
Clone (https): https://your-bitbucket-hostname/scm/key/bitbucket-playground.git
Clone (ssh):   ssh://git@your-bitbucket-hostname:7999/key/bitbucket-playground.git
Web:           https://your-bitbucket-hostname/projects/KEY/repos/bitbucket-playground/browse
```

`--output json` prints the repository exactly as the server returned it, which
is handy together with `jq`:

```
$ bitbucket-cli repo --key KEY --name bitbucket-playground get --output json \
    | jq --raw-output '.links.clone[] | select(.name == "http") | .href'
https://your-bitbucket-hostname/scm/key/bitbucket-playground.git
```

##### Usage

```plain
Usage: bitbucket-cli repo get [--output OUTPUT]

Options:
  --output OUTPUT, -o OUTPUT
                         Output format: text (default) or json
  --help, -h             display this help and exit
```

### Create

Creates a repository in the project.  Requires `PROJECT_ADMIN` permission.

The new repository is named with `--display-name`; Bitbucket derives the
slug from it (`"My Repo"` becomes `my-repo`).  Neither the parent's `--name`
nor a slug detected from the git remote is used by this command; only the
project key is.

```
$ bitbucket-cli repo --key KEY create --display-name "My Repo" --description "Something new" --forkable=false
Name:          My Repo
Slug:          my-repo
ID:            43
Project:       KEY
State:         AVAILABLE
Public:        false
Forkable:      false
Description:   Something new
Clone (https): https://your-bitbucket-hostname/scm/key/my-repo.git
Clone (ssh):   ssh://git@your-bitbucket-hostname:7999/key/my-repo.git
Web:           https://your-bitbucket-hostname/projects/KEY/repos/my-repo/browse
```

Only the options you pass are sent to the server, so its defaults apply to
everything else.  Boolean options take their value with an equals sign:
`--forkable=false` works, `--forkable false` does not.  Older Bitbucket versions
only honour the name and silently ignore the other options.  `--output json` is
available as for `get`.

##### Usage

```plain
Usage: bitbucket-cli repo create --display-name DISPLAY-NAME [--description DESCRIPTION] [--forkable] [--public] [--default-branch DEFAULT-BRANCH] [--output OUTPUT]

Options:
  --display-name DISPLAY-NAME
                         Display name of the new repository; Bitbucket derives the slug from it
  --description DESCRIPTION, -d DESCRIPTION
                         Description of the repository
  --forkable             Whether the repository can be forked (use --forkable=false to disable)
  --public               Whether the repository is publicly accessible (use --public=false to disable)
  --default-branch DEFAULT-BRANCH
                         Default branch of the new repository, e.g: main
  --output OUTPUT, -o OUTPUT
                         Output format: text (default) or json
  --help, -h             display this help and exit
```

### Update

Changes one or more settings of a repository.  Requires `REPO_ADMIN` permission.
At least one option must be given, and only the given ones are changed.

```
$ bitbucket-cli repo --key KEY --name my-repo update --new-name "Renamed Repo" --description ""
Name:          Renamed Repo
Slug:          renamed-repo
...
```

Renaming a repository may change its slug; the updated repository, including the
new slug, is printed afterwards.  Passing `--description ""` clears the
description.  `--to-project OTHERKEY` moves the repository to another project,
which needs admin rights on both projects.  As with `create`, booleans take
their value with an equals sign (`--public=false`).

##### Usage

```plain
Usage: bitbucket-cli repo update [--new-name NEW-NAME] [--description DESCRIPTION] [--forkable] [--public] [--default-branch DEFAULT-BRANCH] [--to-project TO-PROJECT] [--output OUTPUT]

Options:
  --new-name NEW-NAME    New name of the repository (renaming may change the slug)
  --description DESCRIPTION, -d DESCRIPTION
                         New description of the repository (pass "" to clear it)
  --forkable             Whether the repository can be forked (use --forkable=false to disable)
  --public               Whether the repository is publicly accessible (use --public=false to disable)
  --default-branch DEFAULT-BRANCH
                         New default branch, e.g: main
  --to-project TO-PROJECT
                         Key of the project to move the repository to
  --output OUTPUT, -o OUTPUT
                         Output format: text (default) or json
  --help, -h             display this help and exit
```

### Delete

Schedules a repository for deletion.  Requires `REPO_ADMIN` permission.

Without `--yes` the command looks the repository up and asks you to type its
slug on standard input; any other answer aborts.  Scripts must pass `--yes`,
since there is nobody to answer the prompt.

```
$ bitbucket-cli repo --key KEY --name my-repo delete
Delete repository KEY/my-repo ("My Repo")? Type the slug to confirm: my-repo
repository KEY/my-repo scheduled for deletion
```

```
$ bitbucket-cli repo --key KEY --name my-repo delete --yes
repository KEY/my-repo scheduled for deletion
```

##### Usage

```plain
Usage: bitbucket-cli repo delete [--yes]

Options:
  --yes, -y              Skip the confirmation prompt
  --help, -h             display this help and exit
```

### PR

This subcommand deals with PRs, please check its subcommands.

#### Create

Creates a Pull Request.  When run from inside a cloned Bitbucket repository, the project key, slug,
and source branch are all detected automatically.  If the source branch has exactly one commit ahead
of the target, the PR title and description are pre-populated from that commit's subject and body.

Minimal invocation (from inside the repo, on a single-commit feature branch):

```
$ bitbucket-cli repo pr create -T "refs/heads/master"
```

Explicit invocation:

```
$ bitbucket-cli repo -k "KEY" \
    -n "bitbucket-playground" \
    pr create \
    -t "Some Title" \
    -d "Some Description :thumbsup:" \
    -F "refs/heads/feature/2" -T "refs/heads/master"
```

##### Usage

```
Usage: bitbucket-cli repo [-k KEY] [-n NAME] pr create [-t TITLE] [-d DESCRIPTION] [-F FROM-REF] --to-ref TO-REF [--from-key FROM-KEY] [--from-slug FROM-SLUG] [--reviewers REVIEWERS]

Options:
  --title TITLE, -t TITLE
                         Title of this PR; defaults to the commit subject when the branch has
                         exactly one commit ahead of the target
  --description DESCRIPTION, -d DESCRIPTION
                         Description of the PR; defaults to the commit body when the branch has
                         exactly one commit ahead of the target
  --from-ref FROM-REF, -F FROM-REF
                         Source branch, e.g: refs/heads/feature-ABC-123; defaults to the current branch
  --to-ref TO-REF, -T TO-REF
                         Target branch, e.g: refs/heads/master
  --from-key FROM-KEY, -K FROM-KEY
                         Project key of the source repository (if different from target)
  --from-slug FROM-SLUG, -S FROM-SLUG
                         Repository slug of the source repository (if different from target)
  --reviewers REVIEWERS, -r REVIEWERS
                         Comma-separated list of reviewers
  --help, -h             display this help and exit
```

#### List

Lists all the PRs for the chosen repository

```
$ bitbucket-cli repo -k KEY -n bitbucket-playground pr list
Some Title (ID: 2)
feature 1 (ID: 1)
```

```
$ bitbucket-cli repo -k KEY -n bitbucket-playground pr list -s DECLINED
feature 1 (ID: 1)
```

##### Usage

```plain
Usage: bitbucket-cli repo pr list [--state STATE]

Options:
  --state STATE, -s STATE
                         PR State, any of: ALL, OPEN, DECLINED, MERGED
  --help, -h             display this help and exit
```

### Security

#### Scan

##### Usage

```plain
Usage: bitbucket-cli repo security scan

Options:
  --help, -h             display this help and exit
```

##### Example

```plain
bitbucket-cli repo -k ABC -n some-repo security scan
```


## PR

Dashboard-level operations on Pull Requests. Unlike `repo pr`, these commands operate across all repositories visible to the authenticated user.

### List

Lists all Pull Requests visible on your dashboard.

```
$ bitbucket-cli pr list
Some Title - Jane Doe - https://your-bitbucket-hostname/projects/KEY/repos/some-repo/pull-requests/2
feature 1 - John Doe - https://your-bitbucket-hostname/projects/KEY/repos/some-repo/pull-requests/1
```

```
$ bitbucket-cli pr list -s OPEN
```

##### Usage

```plain
Usage: bitbucket-cli pr list [--state STATE] [--output OUTPUT] [--filter-title FILTER-TITLE] [--filter-desc FILTER-DESC]

Options:
  --state STATE, -s STATE
  --output OUTPUT, -o OUTPUT
  --filter-title FILTER-TITLE, -t FILTER-TITLE
  --filter-desc FILTER-DESC, -d FILTER-DESC
  --help, -h             display this help and exit
```

### Create

Creates a Pull Request in the specified repository.

```
$ bitbucket-cli pr create \
  -k "KEY" \
  -n "bitbucket-playground" \
  -t "Some Title" \
  -d "Some Description :thumbsup:" \
  -F "refs/heads/feature/2" -T "refs/heads/master"
```

##### Usage

```plain
Usage: bitbucket-cli pr create --key KEY --name NAME --title TITLE [--description DESCRIPTION] --from-ref FROM-REF --to-ref TO-REF [--from-key FROM-KEY] [--from-slug FROM-SLUG] [--reviewers REVIEWERS]

Options:
  --key KEY, -k KEY      Project key
  --name NAME, -n NAME   Repository slug
  --title TITLE, -t TITLE
                         Title of this PR
  --description DESCRIPTION, -d DESCRIPTION
                         Description of the PR
  --from-ref FROM-REF, -F FROM-REF
                         Reference of the incoming PR, e.g: refs/heads/feature-ABC-123
  --to-ref TO-REF, -T TO-REF
                         Target reference, e.g: refs/heads/master
  --from-key FROM-KEY, -K FROM-KEY
                         Project Key of the "from" repository
  --from-slug FROM-SLUG, -S FROM-SLUG
                         Repository slug of the "from" repository
  --reviewers REVIEWERS, -r REVIEWERS
                         Comma separated list of reviewers
  --help, -h             display this help and exit
```
