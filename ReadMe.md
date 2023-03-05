# Bird Bot

Bird Bot is a discord bot for managing and organizing events for a small discord community.

### Features

- Creating text channels for events
- Notifying when events are created & cancelled
- Delete text channels after events
- Archive text channels after events
- Create recurring weekly events

## Usage

To get up and running, install go and you can run `make run`!

## Using Docker

The container is expecting the config file to be located at `/etc/birdbot/birdbot.yaml`. The easily solution here is to mount the conifg with a volume.

Example: 
```bash
docker run -it -v `pwd`:/etc/birdbot yeslayla/birdbot:latest
```

In this example, your config is in the current directory and call `birdbot.yaml`