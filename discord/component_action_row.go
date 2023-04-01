package discord

import "github.com/bwmarrin/discordgo"

type ActionRow struct {
	components []Component
}

// NewActionRow creates an empty action row component
func (discord *Discord) NewActionRow() *ActionRow {
	return &ActionRow{
		components: []Component{},
	}
}

// NewActionRowWith creates an action row with a set of components
func (discord *Discord) NewActionRowWith(comp []Component) *ActionRow {
	return &ActionRow{
		components: comp,
	}
}

// AddComponent adds a component to the action row
func (row *ActionRow) AddComponent(comp Component) {
	row.components = append(row.components, comp)
}

func (row *ActionRow) toMessageComponent() discordgo.MessageComponent {

	comps := make([]discordgo.MessageComponent, len(row.components))
	for i, v := range row.components {
		comps[i] = v.toMessageComponent()
	}

	return discordgo.ActionsRow{
		Components: comps,
	}
}
