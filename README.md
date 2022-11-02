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
1. Run the following in different terminals
    * `make pubsub-emulator`
    * `make run-discordInteractionEndpoint`
    * `make run-discordInteractionProcessor`

TODO: `make run-playfabCloudScriptHTTP`

#### Discord Setup

Discord is used for the Bot

1. Create a Discord Server
1. Create a Discord Application [here](https://discord.com/developers/applications)
    * Copy the Application ID into the `.env` file ex: `DISCORD_APP_ID=123456789`
    * Copy the Public Key into the `.env` file ex: `DISCORD_APP_PUBLIC_KEY=4e5...847d1`
1. Enable the Bot functionality
    * Copy the Token into the `.env` file ex: `DISCORD_APP_BOT_TOKEN=MTA...5m10y8P6kw`
1. Under OAuth2/URL Generator select the following
    * `bot`, `applications.commands`, `messages.read`
    * Copy the URL into your browser and invite the Bot to your server
1. Run `make discord-ngrok` and copy the URL into the `Interactions Endpoint URL`
    * This URL will change every time so you'll need to copy it every time

TODO: command to register slash commands

#### Playfab Setup

Playfab is used for automatic matchmaking to form groups.
Discord User IDs and GW2 API tokens are stored in player objects.

1. Register for a free [Playfab](https://playfab.com/) account
1. Create a new title
    * Copy the Title ID into the `.env` file ex: `PLAYFAB_TITLE_ID=ABCD0`
1. In Title API Features disable client API access
    * If this is not disabled anyone that knows your Title ID and Discord User IDs can login and use the API
    * This is a big security risk so 100% make sure this is disabled
1. Create a secret key in the title and copy it into the `.env` file ex `PLAYFAB_TITLE_SECRET_KEY=VEW4...46S`

TODO: command to setup matchmaking queues
TODO: command to simulate multiple players for matchmaking testing
