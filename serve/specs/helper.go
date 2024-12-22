package specifications

import (
	"encoding/json"
	"github.com/p10r/pedro/serve/discord"
	"github.com/quii/go-graceful-shutdown/assert"
	"sort"
	"testing"
)

func newDiscordMessage(t *testing.T, input []byte) discord.Message {
	var msg discord.Message
	err := json.Unmarshal(input, &msg)
	assert.NoError(t, err)

	return orderLeagues(msg)
}

// we order the leagues to make sure the output json has always the same structure
func orderLeagues(msg discord.Message) discord.Message {
	sort.Slice(msg.Embeds[0].Fields, func(i, j int) bool {
		leagueName1 := msg.Embeds[0].Fields[i].Name
		leagueName2 := msg.Embeds[0].Fields[j].Name

		return len(leagueName1) < len(leagueName2)
	})

	return msg
}

func prettyPrinted(t *testing.T, msg discord.Message) []byte {
	marshal, err := json.MarshalIndent(msg, "", " ")
	assert.NoError(t, err)
	return marshal
}
