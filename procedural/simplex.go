package procedural

import (
	"github.com/ojrac/opensimplex-go"
)

type SimplexGenerator struct {
	frequency float64

	amplitude float64

	octaves int

	persistence float64

	generator opensimplex.Noise
}

func NewSimplexGenerator(frequency, amplitude, persistence float64, octaves int, seed int64) SimplexGenerator {
	return SimplexGenerator{
		frequency:   frequency,
		amplitude:   amplitude,
		octaves:     octaves,
		persistence: persistence,
		generator:   opensimplex.NewNormalized(seed),
	}
}

func (sg *SimplexGenerator) Noise2D(x, y float64) float64 {
	total := float64(0)
	max := float64(0)

	amp := sg.amplitude
	freq := sg.frequency

	for i := 0; i < sg.octaves; i++ {
		total += sg.generator.Eval2(x*freq, y*freq) * amp

		max += amp
		amp *= sg.persistence

		freq *= 2
	}

	return total / max
}

func (sg *SimplexGenerator) Noise1D(x float64) float64 {
	total := float64(0)
	max := float64(0)

	amp := sg.amplitude
	freq := sg.frequency

	for i := 0; i < sg.octaves; i++ {
		total += sg.generator.Eval2(x*freq, 1) * amp

		max += amp
		amp *= sg.persistence

		freq *= 2
	}

	return total / max
}
