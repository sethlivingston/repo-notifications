package server

type Config struct {
	Server HTTPServerConfig
	GitHub *GitHubConfig
	Slack  *SlackConfig
}

type HTTPServerConfig struct {
	BindAddress string `json:"bindAddress"`
}

type GitHubConfig struct {
	SigningSecret string `json:"signingSecret"`
}

type SlackConfig struct {
	BotToken      string `json:"botToken"`
	SigningSecret string `json:"signingSecret"`
}
