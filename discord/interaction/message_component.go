package interaction

import (
	"context"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/rmb938/gw2groups/discord/interaction/message_component"
)

type MessageComponent struct{}

var messageComponents = map[string]message_component.Component{
	"button_gw2_api_key":                        &message_component.ButtonGw2ApiKey{},
	"select_menu_lfg_category":                  &message_component.SelectMenuLFGCategory{},
	"button_reset_lfg_selection":                &message_component.ButtonResetLFGSelection{},
	"select_menu_lfg_dungeon":                   &message_component.SelectMenuLFGDungeon{},
	"button_reset_lfg_dungeon_mode_story":       nil,
	"button_reset_lfg_dungeon_mode_exploration": nil,
}

func (i *MessageComponent) Handler(ctx context.Context, session *discordgo.Session, interaction *discordgo.Interaction) (*discordgo.InteractionResponse, error) {
	data := interaction.MessageComponentData()

	// we only want to do things based on the last message
	// so if we get a interaction that isn't the last
	// ignore it
	messages, err := session.ChannelMessages(interaction.ChannelID, 1, "", "", "")
	if err != nil {
		return nil, fmt.Errorf("error getting last channel message: %w", err)
	}

	if len(messages) == 0 {
		return &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseDeferredMessageUpdate,
		}, nil
	}

	if messages[0].ID != interaction.Message.ID {
		return &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseDeferredMessageUpdate,
		}, nil
	}

	// Real Message Handler Here
	if component, ok := messageComponents[data.CustomID]; ok {
		return component.Handle(ctx, session, interaction, data)
	}

	return nil, fmt.Errorf("message component %s not implemented", data.CustomID)
}
