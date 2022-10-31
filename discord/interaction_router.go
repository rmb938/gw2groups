package discord

import (
	"context"

	"github.com/bwmarrin/discordgo"
	"github.com/rmb938/gw2groups/discord/interaction"
)

var interactionHandlers = map[discordgo.InteractionType]interaction.Interaction{
	discordgo.InteractionPing:               &interaction.Ping{},
	discordgo.InteractionApplicationCommand: &interaction.ApplicationCommand{},
	discordgo.InteractionMessageComponent:   &interaction.MessageComponent{},
	discordgo.InteractionModalSubmit:        &interaction.ModalSubmit{},
}

func InteractionRouter(ctx context.Context, session *discordgo.Session, i *discordgo.Interaction) (*discordgo.InteractionResponse, error) {
	if handler, ok := interactionHandlers[i.Type]; ok {
		return handler.Handler(ctx, session, i)
	}

	return nil, nil
}
