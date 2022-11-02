package functions

import (
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/rmb938/gw2groups/discord/functions/interaction_endpoint"
	"github.com/rmb938/gw2groups/discord/functions/interaction_processor"
)

func init() {
	functions.HTTP("discordInteractionEndpoint", interaction_endpoint.InteractionEndpoint)
	functions.CloudEvent("discordInteractionProcessor", interaction_processor.InteractionProcessor)
}
