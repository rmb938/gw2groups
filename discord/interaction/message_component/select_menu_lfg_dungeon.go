package message_component

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	playFabAPI "github.com/rmb938/gw2groups/pkg/playfab/api"
	"k8s.io/utils/pointer"
)

type SelectMenuLFGDungeon struct{}

func (c *SelectMenuLFGDungeon) Handle(ctx context.Context, session *discordgo.Session, interaction *discordgo.Interaction, data discordgo.MessageComponentInteractionData) (*discordgo.InteractionResponse, error) {
	playFabClient := playFabAPI.NewPlayFabClient()
	_, err := playFabClient.LoginWithCustomID(ctx, &playFabAPI.LoginWithCustomIDRequest{
		CustomId: pointer.String(interaction.User.ID),
	})
	if err != nil {
		return nil, fmt.Errorf("error logging in with playfab customid: %w", err)
	}

	var dungeonsList []string

	hasAny := false
	for _, value := range data.Values {
		if value == "any" {
			hasAny = true
			break
		}
	}

	selectOptions := interaction.Message.Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.SelectMenu).Options
	if hasAny {
		data.Values = []string{}
		for _, option := range selectOptions {
			if option.Value == "any" {
				continue
			}
			data.Values = append(data.Values, option.Value)
		}
	}

	for _, option := range selectOptions {
		for _, selectedOption := range data.Values {
			if option.Value == selectedOption {
				dungeonsList = append(dungeonsList, option.Label)
			}
		}
	}

	dataValues, err := json.Marshal(data.Values)
	if err != nil {
		return nil, fmt.Errorf("error converting dungeon values to string: %w", err)
	}

	_, err = playFabClient.UpdateUserData(ctx, &playFabAPI.UpdateUserDataRequest{
		Data: map[string]string{
			"lfg_dungeon": string(dataValues),
		},
	})
	if err != nil {
		return nil, fmt.Errorf("error updating playfab user data: %w", err)
	}

	// TODO: ask if story or exploration mode
	//  if exploration as for path (any, path 1, path 2, path 3, save data, start matchmaking
	//  if story, save data, start matchmaking

	return &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("> Category: **Dungeons**\n> Dungeon: **%s**\n \nSelect the dungeon mode", strings.Join(dungeonsList, ", ")),
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.Button{
							Label:    "Story Mode",
							CustomID: "button_lfg_dungeon_mode_story",
							Style:    discordgo.PrimaryButton,
						},
						discordgo.Button{
							Label:    "Exploration Mode",
							CustomID: "button_lfg_dungeon_mode_exploration",
							Style:    discordgo.PrimaryButton,
						},
					},
				},
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.Button{
							Label:    "Reset LFG Selections",
							CustomID: "button_reset_lfg_selection",
							Style:    discordgo.DangerButton,
						},
					},
				},
			},
		},
	}, nil

	// See https://learn.microsoft.com/en-us/gaming/playfab/features/multiplayer/matchmaking/config-examples#hostsearcher-or-role-based-requirements when doing Raids
	//  Specifically the "Games may have role requirements" part for selecting healer, support, tank, ect.. (quickness, alacrity, ect..)

	// TODO: move this all this bellow
	// TODO: do matchmaking here
	//  use https://wiki.guildwars2.com/wiki/API:2/worlds to filter to correct world, find first digit and set lfg_world=1 or =2 in ticket props

	return &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("> Category: **Dungeons**\n> Dungeon: **%s**\n \nNow matchmaking for Dungeons", strings.Join(dungeonsList, ", ")),
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.Button{
							Label:    "Reset LFG Selections",
							CustomID: "button_reset_lfg_selection",
							Style:    discordgo.DangerButton,
						},
					},
				},
			},
		},
	}, nil
}
