# bitbucket-cli

A [Bitbucket Enterprise](https://bitbucket.org/product/enterprise) CLI.

```
Usage: bitbucket-cli [--username USERNAME] [--password PASSWORD] --url URL <command> [<args>]

Options:
  --username USERNAME, -u USERNAME
  --password PASSWORD, -p PASSWORD
  --url URL, -u URL
  --help, -h             display this help and exit

Commands:
  project
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

Most subcommands need to know which repository to operate on. You can supply this explicitly:

- `-k KEY` — project key (e.g. `TOOL`)
- `-n NAME` — repository slug (e.g. `my-repo`)

Or, if you run the command from inside a cloned Bitbucket repository, both are detected automatically
from the `origin` remote URL and can be omitted.

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
