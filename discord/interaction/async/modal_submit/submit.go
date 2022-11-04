package modal_submit

import (
	"context"

	"github.com/bwmarrin/discordgo"
)

type Submit interface {
	Handle(ctx context.Context, session *discordgo.Session, interaction *discordgo.Interaction, data discordgo.ModalSubmitInteractionData) error
}
