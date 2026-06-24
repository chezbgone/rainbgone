package server

import "github.com/joho/godotenv"

type config struct {
	pirateWeatherKey string
}

var appConfig = loadConfig()

func loadConfig() config {
	env, err := godotenv.Read()
	if err != nil {
		panic(err)
	}

	return config{
		pirateWeatherKey: env["PIRATE_WEATHER_KEY"],
	}
}
