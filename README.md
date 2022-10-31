# gw2groups
Discord Bot to provide a better LFG experience for Guild Wars 2

### Development

#### Requirements

* [ngrok](https://ngrok.com/)
* Golang 1.19+

#### Steps

1. Clone this repo
1. Create a `.env` file in the same folder as this repo
1. Follow [Discord Setup](#discord-setup) steps
1. Follow [Playfab Setup](#playfab-setup) steps
1. Run `make run-discord`

TODO: `make run-playfab`

#### Discord Setup

1. Create a Discord Server
1. Create a Discord Application [here](https://discord.com/developers/applications)
    * Copy the Application ID into the `.env` file ex: `DISCORD_APP_ID=123456789`
    * Copy the Public Key into the `.env` file ex: `DISCORD_APP_PUBLIC_KEY=4e5...847d1`
1. Enable the Bot functionality
    * Copy the Token into the `.env` file ex: `DISCORD_APP_BOT_TOKEN=MTA...5m10y8P6kw`
1. Under OAuth2/URL Generator select the following
    * `bot`, `applications.commands`, `messages.read`
    * Copy the URL into your browser and invite the Bot to your server
1. Run `ngrok http 8080` and copy the URL into the `Interactions Endpoint URL`
    * This URL will change every time so you'll need to copy it every time

#### Playfab Setup

1. Create a new title
    * Copy the Title ID into the `.env` file ex: `PLAYFAB_TITLE_ID=ABCD0`

TODO: command to setup matchmaking queues

