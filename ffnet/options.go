package ffnet

import (
	"math/rand"

	"github.com/spy16/pkg/ffnet/mat"
)

// Layer adds a layer of given size and activations to the network.
func Layer(size int, activation Differentiable) Option {
	return func(net *FFNet) error {
		inputSz := net.inputSz
		if len(net.layers) > 0 {
			inputSz = net.layers[len(net.layers)-1].layerSz
		}

		weights := mat.New(size, inputSz).Apply(func(_, _ int, val float64) float64 {
			return rand.Float64()
		})

		net.layers = append(net.layers, layer{
			inputSz:     inputSz,
			layerSz:     size,
			actFn:       activation,
			weights:     weights,
			biases:      mat.New(size, 1),
			activations: mat.New(size, 1),
		})
		return nil
	}
}
