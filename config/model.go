package config

type TgBotConfig struct {
	Active         bool     `json:"active,omitempty" yaml:"active"`
	WebMode        bool     `json:"web,omitempty" yaml:"web"`
	Token          string   `json:"token,omitempty" yaml:"token"`
	Login          string   `json:"login,omitempty" yaml:"login"`
	AdminId        int64    `json:"admin,omitempty" yaml:"admin"`
	AdminLogin     string   `json:"adminlogin,omitempty" yaml:"adminlogin"`
	DevChat        int64    `json:"chat,omitempty" yaml:"chat"`
	Uri            string   `json:"uri,omitempty" yaml:"uri"`
	DefaultCommand string   `json:"command,omitempty" yaml:"command"`
	Services       []string `json:"services,omitempty" yaml:"services"`
	Description    string   `json:"description,omitempty" yaml:"description"`
}
