# Bird Bot

Bird Bot is a discord bot for managing and organizing events for a small discord community.

## Features

- Creating text channels for events
- Delete/archive text channels after events
- Notifying when events are created & cancelled
- Role selection
- Plugin support

## Usage

To get up and running, install go and you can run `make run`!

## Using Docker

The container is expecting the config file to be located at `/etc/birdbot/birdbot.yaml`. The easily solution here is to mount the config with a volume.

Example:

```bash
docker run -it -v `pwd`:/etc/birdbot yeslayla/birdbot:latest
```

In this example, your config is in the current directory and call `birdbot.yaml`

### Persistant Data

The default location for container data is `/var/lib/birdbot/` so you can mount it like:

Example:

```bash
docker run -it -v `pwd`:/var/lib/birdbot/ yeslayla/birdbot:latest
```
