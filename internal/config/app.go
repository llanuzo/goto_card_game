package config

type App struct {
	HttpPort int
}

func NewApp() App {
	return App{
		HttpPort: 8080,
	}
}
