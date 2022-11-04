package async

import (
	"context"
	"fmt"

	"cloud.google.com/go/pubsub"
	"github.com/bwmarrin/discordgo"
	"github.com/rmb938/gw2groups/discord/interaction/async/application_command"
)

type ApplicationCommand struct{}

var applicationCommands = map[string]application_command.Command{
	"lfg": &application_command.LFG{},
}

func (i *ApplicationCommand) Handler(ctx context.Context, session *discordgo.Session, pubsubTopicPlayfabMatchmakingTickets *pubsub.Topic, interaction *discordgo.Interaction) error {
	data := interaction.ApplicationCommandData()

	if command, ok := applicationCommands[data.Name]; ok {
		return command.Handle(ctx, session, interaction, data)
	}

	return fmt.Errorf("command %s not implemented", data.Name)
}
