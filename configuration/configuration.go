package configuration

import "github.com/sirupsen/logrus"

type EngineConfig struct {
	ScreenWidth  int
	ScreenHeight int

	FullScreen bool

	WindowTitle    string
	PolygonLines   bool
	CollisionLines bool

	ShowFPS bool

	Dimensions int

	Profiling      bool
	SingleMaterial bool

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
		FullScreen:     true,
		WindowTitle:    "game",
		PolygonLines:   false,
		CollisionLines: false,
		ShowFPS:        false,
		Dimensions:     Dimensions,
		Profiling:      false,
		SingleMaterial: false,
		Logger:         logrus.New(),
	}
}
