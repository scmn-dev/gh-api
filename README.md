<p align="center">
  <img alt="octocat" src="https://github.githubassets.com/images/icons/emoji/octocat.png" />
  <h1 align="center">GitHub API</h1>
  <h3 align="center">Github API For Authentication and Manage Repos.</h3>
</p>

> Commands that need `gh-api`

* [**Auth**](https://docs.secman.dev/commands/auth)
  - [**Login**](https://docs.secman.dev/commands/auth#login)
  - [**Logout**](https://docs.secman.dev/commands/auth#logout)
  - [**Refresh**](https://docs.secman.dev/commands/auth#refresh)
  - [**Status**](https://docs.secman.dev/commands/auth#status)
* [**Repo**](https://docs.secman.dev/commands/repo)
  - [**Clone**](https://docs.secman.dev/commands/repo#clone)
  - [**Create**](https://docs.secman.dev/commands/repo#create)
  - [**Fork**](https://docs.secman.dev/commands/repo#fork)
  - [**List**](https://docs.secman.dev/commands/repo#list)
* [**Config**](https://docs.secman.dev/commands/config)
  - [**Get**](https://docs.secman.dev/commands/config#get)
  - [**Set**](https://docs.secman.dev/commands/config#set)

## Auth

```
secman auth -h
secman auth login
secman auth logout
secman auth refresh
secman auth status
```

## Repo

```
secman repo -h
secman repo clone
secman repo create
secman repo fork
secman repo list
```

## Config

```
secman gh-config -h
secman gh-config get ...
secman gh-config set ...
```
