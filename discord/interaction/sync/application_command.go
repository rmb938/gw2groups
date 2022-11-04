package sync

import (
	"context"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/rmb938/gw2groups/discord/interaction/sync/application_command"
)

type ApplicationCommand struct{}

var applicationCommands = map[string]application_command.Command{
	"lfg": &application_command.LFG{},
}

func (i *ApplicationCommand) Handler(ctx context.Context, session *discordgo.Session, interaction *discordgo.Interaction) (*discordgo.InteractionResponse, error) {
	data := interaction.ApplicationCommandData()

	if command, ok := applicationCommands[data.Name]; ok {
		if interaction.Member != nil {
			if !command.CanGuild() {
				return &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Slash Command is not usable in Servers",
					},
				}, nil
			}
		}

		if interaction.User != nil {
			if !command.CanDM() {
				return &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Slash Command is not usable in DMs",
					},
				}, nil
			}
		}

		return command.Handle(ctx, session, interaction, data)
	}

	return nil, fmt.Errorf("command %s not implemented", data.Name)
}
