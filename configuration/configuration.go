package configuration

import "github.com/sirupsen/logrus"

type EngineConfig struct {
	ScreenWidth  int
	ScreenHeight int

	WindowTitle    string
	PolygonLines   bool
	CollisionLines bool

	Dimensions int

	Logger *logrus.Logger
}

func NewEngineConfig(
	ScreenWidth,
	ScreenHeight,
	Dimensions int,
) EngineConfig {
	return EngineConfig{
		ScreenWidth:    ScreenWidth,
		ScreenHeight:   ScreenHeight,
		WindowTitle:    "game",
		PolygonLines:   false,
		CollisionLines: false,
		Dimensions:     Dimensions,
		Logger:         logrus.New(),
	}
}
