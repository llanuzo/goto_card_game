package config

type App struct {
	AppName      string
	Env          string
	HttpPort     int
	MetricsPort  int
	PProfEnabled bool
}

func NewApp() App {
	return App{
		AppName:      "card-game",
		Env:          "local",
		HttpPort:     8080,
		MetricsPort:  10001,
		PProfEnabled: true,
	}
}
