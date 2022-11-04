package async

import (
	"context"
	"fmt"

	"cloud.google.com/go/pubsub"
	"github.com/bwmarrin/discordgo"
	"github.com/rmb938/gw2groups/discord/interaction/async/message_component"
)

type MessageComponent struct{}

var messageComponents = map[string]message_component.Component{
	"button_gw2_api_key":            nil,
	"button_reset_lfg_selection":    &message_component.ButtonResetLFGSelection{},
	"select_menu_lfg_category":      &message_component.SelectMenuLFGCategory{},
	"select_menu_lfg_dungeon":       &message_component.SelectMenuLFGDungeon{},
	"button_lfg_dungeon_mode_story": &message_component.ButtonLFGDungeonModeStory{},
	// "button_lfg_dungeon_mode_exploration": nil,
}

func (i *MessageComponent) Handler(ctx context.Context, session *discordgo.Session, pubsubTopicPlayfabMatchmakingTickets *pubsub.Topic, interaction *discordgo.Interaction) error {
	data := interaction.MessageComponentData()

	// we only want to do things based on the last message
	// so if we get a interaction that isn't the last
	// ignore it
	messages, err := session.ChannelMessages(interaction.ChannelID, 1, "", "", "")
	if err != nil {
		return fmt.Errorf("error getting last channel message: %w", err)
	}

	if len(messages) == 0 {
		return nil
	}

	if messages[0].ID != interaction.Message.ID {
		return nil
	}

	// Real Message Handler Here
	if component, ok := messageComponents[data.CustomID]; ok {
		if component == nil {
			return nil
		}

		return component.Handle(ctx, session, pubsubTopicPlayfabMatchmakingTickets, interaction, data)
	}

	return fmt.Errorf("message component %s not implemented", data.CustomID)
}
