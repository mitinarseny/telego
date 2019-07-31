<p align="center">
    <a href="https://github.com/mitinarseny/telego">
        <img src="_assets/logo.png" alt="telego logo" width="20%" />
    </a>
    <h1 align="center">telego</h1>
    <p align="center">Docker Go template for creating <a href="https://core.telegram.org/bots">Telegram Bots</a></p>
    <p align="center">
      <a href="https://travis-ci.org/mitinarseny/telego"><img alt="TravisCI" src="https://img.shields.io/travis/mitinarseny/telego/master.svg?style=flat-square&logo=travis-ci"></a>
      <a href="https://golangci.com/r/github.com/mitinarseny/telego"><img src="https://golangci.com/badges/github.com/mitinarseny/telego.svg"></a>
      <a href="/LICENSE.md"><img alt="Software License" src="https://img.shields.io/badge/license-MIT-brightgreen.svg?style=flat-square"></a>
      <a href="https://saythanks.io/to/mitinarseny"><img alt="SayThanks.io" src="https://img.shields.io/badge/say-thanks-9933ff.svg?style=flat-square"></a>
    </p>
</p>

## Table of Contents
* [Usage](#usage)
  * [Create Bot](#create-bot)
  * [Copy Token](#copy-token)
    * [Notifier](#notifier)
  * [Code](#code)
    * [Logic](#logic)
    * [Handlers](#handlers)
  * [Run](#run)
* [Debug](#debug)
  * [Build & Run](#build--run)
  * [Attach](#attach)

## Usage
### Create Bot
Create new bot with [@BotFather](https://t.me/BotFather).
### Copy Token
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
### Code
#### Logic
Main logic of the bot should be implemented inside `HandleUpdates` method of `Bot` in [`bot/handlers/core.go`](bot/handlers/core.go):
```go
func (b *Bot) HandleUpdates(updates tgbotapi.UpdatesChannel, errCh chan <- error) error {
    for update := range updates {
        go func() {
            switch {
            case update.Message != nil:
                switch {
                case update.Message.Command() == "hello":
                    if err := b.HandleHello(update); err != nil {
                        errCh <- err
                    }
                }
            }
        }()
    }
    return nil
}
```
Take a look at more complex example in [`bot/handlers/core.go`](bot/handlers/core.go).
#### Handlers
All hanlders should be placed in [`bot/handlers/`](bot/handlers). Here is an example from [`hello.go`](bot/handlers/hello.go):
```go
func (b *Bot) HandleHello(update tgbotapi.Update) error {
    msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
    msg.ReplyToMessageID = update.Message.MessageID

    _, err := b.Send(msg)
    return err
}

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
