<p align="center">
    <a href="https://github.com/mitinarseny/telego">
        <img src="_assets/logo.png" alt="telego logo" width="20%" />
    </a>
    <h1 align="center">telego</h1>
    <p align="center">Docker Go Template for creating new <a href="https://core.telegram.org/bots">Telegram Bots</a></p>
    <p align="center">
      <a href="https://travis-ci.org/mitinarseny/telego"><img alt="TravisCI" src="https://img.shields.io/travis/mitinarseny/telego/master.svg?style=flat-square&logo=travis-ci"></a>
      <a href="https://golangci.com/r/github.com/mitinarseny/telego"><img src="https://golangci.com/badges/github.com/mitinarseny/telego.svg"></a>
      <a href="/LICENSE.md"><img alt="Software License" src="https://img.shields.io/badge/license-MIT-brightgreen.svg?style=flat-square"></a>
      <a href="https://saythanks.io/to/mitinarseny"><img alt="SayThanks.io" src="https://img.shields.io/badge/say-thanks-9933ff.svg?style=flat-square"></a>
    </p>
</p>

## Table of Contents
* [Usage](#usage)
* [Debug](#debug)

## Usage
### Create Bot
Create new bot with [@BotFather](https://t.me/BotFather).
### Copy token
Create file `./docker-compose.secret.yaml` with the following structure and paste the token from [@BotFather](https://t.me/BotFather):
```yaml
# ./docker-compose.secret.yaml

version: '3.7'

services:
  bot:
    environment:
      TELEGO_BOT_TOKEN: "12345689:ABCdEFgHi1JKLMNO23P45rSTU6vw78xyz-a"
```
#### Notifier
You can enable Telegram notifications on your bot's status (`UP` or `DOWN`) by creating another bot and a group chat with this bot. Then edit `./docker-compose.secret.yaml`:
```yaml
# ...
environment:
  TELEGO_NOTIFIER_BOT_TOKEN: "<token>"
  TELEGO_NOTIFIER_CHAT_ID: "<chat_id>"
```

### Run
```bash
docker-compose \
  -f docker-compose.yaml \
  -f docker-compose.secret.yaml \
  up --build -d
```

## Debug
You can debug your code with [Delve](https://github.com/go-delve/delve) debugger. 
### Build & Run
To enable [dlv](https://github.com/go-delve/delve) debugger inside the container run:
```bash
docker-compose \
  -f docker-compose.yaml \ 
  -f docker-compose.dev.yaml \
  -f docker-compose.secret.yaml \
  up --build -d
``` 
### Attach
```bash
${GOPATH}/bin/dlv connect localhost:40000
```
