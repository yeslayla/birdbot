package core

import "strings"

// Config is used to modify the behavior of birdbot externally
type Config struct {
	Discord      DiscordConfig  `yaml:"discord"`
	Mastodon     MastodonConfig `yaml:"mastodon"`
	Feedback     Feedback       `yaml:"feedback"`
	StatusPortal StatusPortal   `yaml:"status_portal"`
	Features     Features       `yaml:"features"`
}

// DiscordConfig contains discord specific configuration
type DiscordConfig struct {
	Token         string `yaml:"token" env:"DISCORD_TOKEN"`
	ApplicationID string `yaml:"application_id" env:"DISCORD_APPLICATION_ID"`
	GuildID       string `yaml:"guild_id" env:"DISCORD_GUILD_ID"`

	EventCategory       string `yaml:"event_category" env:"DISCORD_EVENT_CATEGORY"`
	ArchiveCategory     string `yaml:"archive_category" env:"DISCORD_ARCHIVE_CATEGORY"`
	NotificationChannel string `yaml:"notification_channel" env:"DISCORD_NOTIFICATION_CHANNEL"`

	RoleSelections []RoleSelectionConfig `yaml:"role_selection"`

	ChatLinks map[string][]string `yaml:"chat_links"`
}

type Feedback struct {
	WebhookURL  string `yaml:"url" env:"BIRD_FEEDBACK_URL"`
	PayloadType string `yaml:"type" env:"BIRD_FEEDBACK_TYPE"`

	SuccessMessage string `yaml:"success_message"`
	FailureMessage string `yaml:"failure_message"`
}

type StatusPortal struct {
	URL string `yaml:"url" env:"BIRD_STATUS_PORTAL_URL"`
}

type RoleSelectionConfig struct {
	Title              string  `yaml:"title"`
	Description        string  `yaml:"description"`
	GenerateColorEmoji Feature `yaml:"generate_color_emoji"`

	SelectionChannel string       `yaml:"discord_channel"`
	Roles            []RoleConfig `yaml:"roles"`
}

type RoleConfig struct {
	RoleName string `yaml:"name"`
	Color    string `yaml:"color"`
}

// MastodonConfig contains mastodon specific configuration
type MastodonConfig struct {
	Server       string `yaml:"server" env:"MASTODON_SERVER"`
	Username     string `yaml:"user" env:"MASTODON_USER"`
	Password     string `yaml:"password" env:"MASTODON_PASSWORD"`
	ClientID     string `yaml:"client_id" env:"MASTODON_CLIENT_ID"`
	ClientSecret string `yaml:"client_secret" env:"MASTODON_CLIENT_SECRET"`
}

// Features contains all features flags that can be used to modify functionality
type Features struct {
	ManageEventChannels Feature `yaml:"manage_event_channels" env:"BIRD_EVENT_CHANNELS"`
	AnnounceEvents      Feature `yaml:"announce_events" env:"BIRD_ANNOUNCE_EVENTS"`
	RecurringEvents     Feature `yaml:"recurring_events" env:"BIRD_RECURRING_EVENTS"`
	RoleSelection       Feature `yaml:"role_selection" env:"BIRD_ROLE_SELECTION"`
	Feedback            Feature `yaml:"feedback" env:"BIRD_FEEDBACK"`
	LoadGamePlugins     Feature `yaml:"load_game_plugins" env:"BIRD_LOAD_GAME_PLUGINS"`
	StatusPortal        Feature `yaml:"status_portal" env:"BIRD_STATUS_PORTAL"`
}

// Feature is a boolean string used to toggle functionality
type Feature string

// IsEnabled returns true when a feature is set to be true
func (value Feature) IsEnabled() bool {
	return strings.ToLower(string(value)) == "true"
}

// IsEnabled returns true when a feature is set to be true
// or if the feature flag is not set at all
func (value Feature) IsEnabledByDefault() bool {
	v := strings.ToLower(string(value))
	if v == "" {
		v = "true"
	}
	return Feature(v).IsEnabled()
}
