package specifications

import (
	"encoding/json"
	"github.com/p10r/pedro/pedro/domain/expect"
	"github.com/p10r/pedro/serve/discord"
	"sort"
	"testing"
)

func newDiscordMessage(t *testing.T, input []byte) discord.Message {
	var msg discord.Message
	err := json.Unmarshal(input, &msg)
	expect.NoErr(t, err)

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
	expect.NoErr(t, err)
	return marshal
}
