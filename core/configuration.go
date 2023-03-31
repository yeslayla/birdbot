package core

type Config struct {
	Discord  DiscordConfig  `yaml:"discord"`
	Mastodon MastodonConfig `yaml:"mastodon"`
	Features Features       `yaml:"features"`
}

type DiscordConfig struct {
	Token   string `yaml:"token" env:"DISCORD_TOKEN"`
	GuildID string `yaml:"guild_id" env:"DISCORD_GUILD_ID"`

	EventCategory       string `yaml:"event_category" env:"DISCORD_EVENT_CATEGORY"`
	ArchiveCategory     string `yaml:"archive_category" env:"DISCORD_ARCHIVE_CATEGORY"`
	NotificationChannel string `yaml:"notification_channel" env:"DISCORD_NOTIFICATION_CHANNEL"`
}

type MastodonConfig struct {
	Server       string `yaml:"server" env:"MASTODON_SERVER"`
	Username     string `yaml:"user" env:"MASTODON_USER"`
	Password     string `yaml:"password" env:"MASTODON_PASSWORD"`
	ClientID     string `yaml:"client_id" env:"MASTODON_CLIENT_ID"`
	ClientSecret string `yaml:"client_secret" env:"MASTODON_CLIENT_SECRET"`
}

type Features struct {
	ManageEventChannels bool `yaml:"manage_event_channels" env:"BIRD_EVENT_CHANNELS" env-default:"true"`
	AnnounceEvents      bool `yaml:"announce_events" env:"BIRD_ANNOUNCE_EVENTS" env-default:"true"`
	ReccurringEvents    bool `yaml:"recurring_events" env:"BIRD_RECURRING_EVENTS" env-default:"true"`
	LoadGamePlugins     bool `yaml:"load_game_plugins" env:"BIRD_LOAD_GAME_PLUGINS" env-default:"true"`
}
