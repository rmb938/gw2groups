package sync

import (
	"context"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/rmb938/gw2groups/discord/interaction/sync/message_component"
)

type MessageComponent struct{}

var messageComponents = map[string]message_component.Component{
	"button_gw2_api_key":                  &message_component.ButtonGw2ApiKey{},
	"button_reset_lfg_selection":          nil,
	"select_menu_lfg_category":            nil,
	"select_menu_lfg_dungeon":             nil,
	"button_lfg_dungeon_mode_story":       nil,
	"button_lfg_dungeon_mode_exploration": nil,
}

func (i *MessageComponent) Handler(ctx context.Context, session *discordgo.Session, interaction *discordgo.Interaction) (*discordgo.InteractionResponse, error) {
	data := interaction.MessageComponentData()

	if component, ok := messageComponents[data.CustomID]; ok {
		if component == nil {
			return nil, nil
		}

		return component.Handle(ctx, session, interaction, data)
	}

	return nil, fmt.Errorf("message component %s not implemented", data.CustomID)
}
