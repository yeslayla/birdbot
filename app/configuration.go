package app

type Config struct {
	Discord DiscordConfig `yaml:"discord"`
}

type DiscordConfig struct {
	Token   string `yaml:"token" env:"DISCORD_TOKEN"`
	GuildID string `yaml:"guild_id" env:"DISCORD_GUILD_ID"`

	EventCategory       string `yaml:"event_category" env:"DISCORD_EVENT_CATEGORY"`
	ArchiveCategory     string `yaml:"archive_category" env:"DISCORD_ARCHIVE_CATEGORY"`
	NotificationChannel string `yaml:"notification_channel" env:"DISCORD_NOTIFICATION_CHANNEL"`
}
