package configuration

import "github.com/sirupsen/logrus"

type EngineConfig struct {
	ScreenWidth  int
	ScreenHeight int

	FullScreen bool
	VSync      bool

	WindowTitle    string
	PolygonLines   bool
	CollisionLines bool

	ShowFPS bool

	MaxFPS int

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
		VSync:          true,
		WindowTitle:    "game",
		PolygonLines:   false,
		CollisionLines: false,
		ShowFPS:        false,
		MaxFPS:         60,
		Dimensions:     Dimensions,
		Profiling:      false,
		SingleMaterial: false,
		Logger:         logrus.New(),
	}
}
