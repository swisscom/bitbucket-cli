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
