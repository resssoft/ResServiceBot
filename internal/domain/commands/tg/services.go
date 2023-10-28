package tgModel

type Service interface {
	Commands() Commands
	Name() string // TODO: use in the tgBot after append commands
	Configure(ServiceConfig)
}

type ServiceConfig struct {
}
