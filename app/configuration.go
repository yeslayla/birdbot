package app

type Config struct {
	Discord DiscordConfig `yaml:"discord"`
}

type DiscordConfig struct {
	Token   string `yaml:"token" env:"DISCORD_TOKEN"`
	GuildID string `yaml:"guild_id" env:"DISCORD_GUILD_ID"`

	EventCategory       string `yaml:"event_category"`
	ArchiveCategory     string `yaml:"archive_category"`
	NotificationChannel string `yaml:"notification_channel"`
}
