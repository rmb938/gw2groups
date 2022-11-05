package async

import (
	"context"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/rmb938/gw2groups/discord/interaction/async/modal_submit"
)

type ModalSubmit struct{}

var modalSubmits = map[string]modal_submit.Submit{
	"modals_gw2_api_key": &modal_submit.ModalsGw2ApiKey{},
}

func (i *ModalSubmit) Handler(ctx context.Context, session *discordgo.Session, interaction *discordgo.Interaction) error {
	data := interaction.ModalSubmitData()

	// Real Message Handler Here
	if component, ok := modalSubmits[data.CustomID]; ok {
		return component.Handle(ctx, session, interaction, data)
	}

	return fmt.Errorf("modal submit %s not implemented", data.CustomID)
}
