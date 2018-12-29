package configuration

import "github.com/sirupsen/logrus"

type EngineConfig struct {
	ScreenWidth  int
	ScreenHeight int
	FullScreen   bool

	VSync           bool
	GammaCorrection bool
	AntiAliasing    bool

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
		WindowTitle: "game",

		// Screen dimensions
		ScreenWidth:  ScreenWidth,
		ScreenHeight: ScreenHeight,
		FullScreen:   true,
		Dimensions:   Dimensions,

		// Rendering
		VSync:           true,
		PolygonLines:    false,
		GammaCorrection: true,
		AntiAliasing:    true,

		// Misc
		CollisionLines: false,
		ShowFPS:        false,
		MaxFPS:         70,
		Profiling:      false,
		SingleMaterial: false,
		Logger:         logrus.New(),
	}
}
