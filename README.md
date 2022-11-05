# gw2groups
Discord Bot to provide a better LFG experience for Guild Wars 2

### Feature TODO's

List of Features I would like to implement, these are not in priority order

If you see a feature not on this list feel free to make an Issue.

- [ ] Expose HTTP API
- [ ] Web UI
    * Requires HTTP API
    * I'm not a web developer so if someone wants to work on this let me know
- [ ] BlishHUD UI
    * Requires HTTP API
    * I'm not a C# developer so if someone wants to work on this let me know
- [ ] Large Community Private LFG
    * Currently, LFG create matches across all of Discord
    * Allow communities to make their own matching queues
- [ ] Dynamic Community Events
    * Allow communities to make events for everyone to join
    * i.e Event/HP Trains
    * Commander runs /lfg-train, enters the train details
    * Other players can join the train and it'll print out the /sqjoin command
- [ ] Automatic [Major Event Timers](https://wiki.guildwars2.com/wiki/Event_timers)
    * These probably aren't needed as there typically are always groups present in the native LFG
- [ ] Fractals
    * All Tiers
    * Agony Resistance Requirements
    * include role selection (dps, condi, quick, alac, ect..)
- [ ] Strike Missions
    * include role selection (dps, condi, quick, alac, ect..)
- [ ] Raids
    * include role selection (dps, condi, quick, alac, ect..)
- [ ] Holiday Event Groups
- [ ] [killproof.me](https://killproof.me/) Integration

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
    * `make run`

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
1. Run `make discord-ngrok` and copy the URL into the `Interactions Endpoint URL` and add `/discord/interactions` to the end
    * ex: `https://9899-2602-43-442-1800-886a-13ff-fe01-3516.ngrok.io/discord/interactions`
    * This URL will change every time, so you'll need to copy it every time

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
