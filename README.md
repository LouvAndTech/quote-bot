<p align="center">
  <a href="icon.png" target="_blank">
    <img src="icon.png" alt="Logo" width="128" height="128">
  </a>
</p>

# Quote Bot
![GitHub Actions Workflow Status](https://img.shields.io/github/actions/workflow/status/louvandtech/quote-bot/docker-build-push.yml?style=for-the-badge&label=Build%20%26%20Push)
![GitHub Actions Workflow Status](https://img.shields.io/github/actions/workflow/status/louvandtech/quote-bot/docker-build-validation.yml?style=for-the-badge&label=Build%20on%20main)

![Docker Pulls](https://img.shields.io/docker/pulls/louvandtech/quote-bot?style=for-the-badge)
![Docker Image Size (latest by date)](https://img.shields.io/docker/image-size/louvandtech/quote-bot?style=for-the-badge)
![Docker Image Version (latest by date)](https://img.shields.io/docker/v/louvandtech/quote-bot?style=for-the-badge)

![GitHub License](https://img.shields.io/github/license/louvandtech/quote-bot?style=for-the-badge)



## What is Quote Bot?
It is a bot that will help you format you quotes quickly and easily for your discord server.
It is also able to give you a random quote from your own serer.

### It's particularly:
It does not have any serverside database, and rely on discord's own server to store the quotes, which means that the quotes are only available in the server where the bot is invited and cannot be query from outside.

## How to use it?

### Invite the bot to your server:
1. First of all, invite the bot to your server by clicking [here](https://discord.com/api/oauth2/authorize?client_id=1090565153495453716&permissions=11264&scope=bot). 
2. Then You need to create a channel that include `quotes` in it's name in your server.

After that you are ready to go: 
- You can use the bot by typing `/quotization` to format your quotes
- You can type `/quote` to get a random quote

### Host it yourself?
You can find the docker image [here](https://hub.docker.com/r/louvandtech/quote-bot)

To deploy the bot, you can use the following `docker-compose.yml` file:
```yaml
version: "3.8"
services:
  quote-bot:
    image: louvandtech/quote-bot:latest
    container_name: quote-bot
    restart: unless-stopped
    environment:
      TOKEN: <YOUR-DISCORD-TOKEN>
```
