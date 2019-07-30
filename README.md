<p align="center">
    <a href="https://github.com/mitinarseny/telego">
        <img src="assets/logo.png" alt="telego logo" width="30%" />
    </a>
    <h1 align="center">telego</h1>
    <p align="center">Go Template for Telegram Bot</p>
    <p align="center">
      <a href="https://travis-ci.org/mitinarseny/telego"><img alt="TravisCI" src="https://img.shields.io/travis/mitinarseny/telego/master.svg?style=flat-square&logo=travis-ci"></a>
      <a href="https://golangci.com/r/github.com/mitinarseny/telego"><img src="https://golangci.com/badges/github.com/mitinarseny/telego.svg"></a>
      <a href="/LICENSE.md"><img alt="Software License" src="https://img.shields.io/badge/license-MIT-brightgreen.svg?style=flat-square"></a>
      <a href="https://saythanks.io/to/mitinarseny"><img alt="SayThanks.io" src="https://img.shields.io/badge/say-thanks-9933ff.svg?style=flat-square"></a>
    </p>
</p>

## Table of Contents
* [Usage](#usage)

## Usage
* Create new bot with [@BotFather](https://t.me/BotFather) and copy the token (example: `12345689:ABCdEFgHi1JKLMNO23P45rSTU6vw78xyz-a`)
* Paste copied token in [docker-compose.yaml](./docker-compose.yaml):
  ```yaml
  services:
    bot:
      # ...
      environment:
        TELEGO_BOT_TOKEN: "<your_token_here>"
  ```
* Run services:
  ```bash
  docker-compose up -d
  ```
