package config

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type TelegramConfig struct {
	Enbaled      bool    `yaml:"enabled"`
	BotToken     string  `yaml:"bot_token"`
	AllowedUsers []int64 `yaml:"allowed_users"`
	Model        string  `yaml:"model"`
}
type Config struct {
	Model        string         `yaml:"model"`
	BaseURL      string         `yaml:"base_url"`
	APIKey       string         `yaml:"api_key"`
	MaxHistory   int            `yaml:"max_history"`
	Timeout      int            `yaml:"timeout"`
	SafeMode     bool           `yaml:"safe_mode"`
	AllowedPaths []string       `yaml:"allowed_paths"`
	SystemPrompt string         `yaml:"system_prompt"`
	Telegram     TelegramConfig `yaml:"telegram"`
}

type Session struct {
	ID        string    `json:"id"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Messages  []any     `json:"messages"`
	ToolCalls []ToolLog `json:"tool_calls"`
}

type ToolLog struct {
	Time   time.Time `json:"time"`
	Tool   string    `json:"tool"`
	Args   string    `json:"args"`
	Result string    `json:"result"`
}

func Load(path string) (*Config, error) {
	cfg := Default()
	data, err := os.ReadFile(path)
	if err != nil {
		return cfg, nil // fall back to default if no config file
	}

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

func Default() *Config {
	return &Config{
		Model:        "openai/gpt-oss-120b:free",
		BaseURL:      "https://openrouter.ai/api/v1",
		APIKey:       "",
		MaxHistory:   20,
		Timeout:      300,
		SafeMode:     true,
		AllowedPaths: []string{"./output", "/tmp"},
		SystemPrompt: `You are a professional cybersecurity red team expert.
	When a tool returns "unsafe path", tell the user explicitly that the path
	was rejected for security reasons and ask them to use a safe path instead.
	Never claim a task succeeded if a tool returned an error.`,
	}
}
