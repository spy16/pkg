package ffnet

import (
	"errors"
	"fmt"
	"sync"

	"github.com/spy16/pkg/ffnet/mat"
)

// New initializes a new artificial neural network with given configuration.
// Use Layer() option to add layers to the network.
func New(inputSz int, opts ...Option) (*FFNet, error) {
	net := &FFNet{inputSz: inputSz}
	for _, opt := range opts {
		if err := opt(net); err != nil {
			return nil, err
		}
	}
	if len(net.layers) == 0 {
		return nil, errors.New("need at-least an output layer")
	}
	net.outputSz = net.layers[len(net.layers)-1].layerSz
	return net, nil
}

// Option values can be used with New to configure network initialization.
type Option func(net *FFNet) error

// FFNet represents a fully-connected feed-forward artificial neural network.
type FFNet struct {
	mu       sync.RWMutex
	layers   []layer
	inputSz  int
	outputSz int
}

// Predict performs a forward pass of inputs and returns the predictions generated.
func (net *FFNet) Predict(inputs ...float64) ([]float64, error) {
	_, as, err := net.forwardPass(inputs)
	if err != nil {
		return nil, err
	}
	return as[len(as)-1].Values(), nil
}

func (net *FFNet) String() string {
	return fmt.Sprintf("FFNet{in=%d, out=%d}", net.inputSz, net.outputSz)
}

func (net *FFNet) forwardPass(inputs []float64) (zs, as []mat.Matrix, err error) {
	if len(inputs) != net.inputSz {
		return zs, as, fmt.Errorf("need exactly %d inputs, got %d", net.inputSz, len(inputs))
	}

	net.mu.RLock()
	defer net.mu.RUnlock()

	a := mat.From(net.inputSz, 1, inputs...)

	var z mat.Matrix
	for _, l := range net.layers {
		// find z = w.x + b and a = g(z)
		z = mat.Add(mat.Dot(l.weights, a), l.biases)
		a = l.actFn.F(z)

		zs = append(zs, z)
		as = append(as, a)
	}
	return zs, as, nil
}

func (net *FFNet) backPropagate(zs, as []mat.Matrix, x, costGrad mat.Matrix) (deltaB, deltaW []mat.Matrix) {
	deltas := make([]mat.Matrix, len(net.layers))
	deltaW = make([]mat.Matrix, len(net.layers))
	deltaB = make([]mat.Matrix, len(net.layers))

	L := len(net.layers) - 1 // index of last layer
	for l := len(net.layers) - 1; l >= 0; l-- {
		gPrime := net.layers[l].actFn.F(zs[l])
		deltas[l] = mat.Mul(costGrad, gPrime)

		// compute weight and bias updates for this layer
		deltaB[l] = deltas[l]
		if l == 0 {
			// for input layer (l=0), input is x (not available in 'as')
			deltaW[l] = mat.Dot(deltas[l], x.T())
		} else {
			// for non-input layers input is the activations of previous
			// layer (i.e., as[l-1])
			deltaW[l] = mat.Dot(deltas[l], as[l-1].T())
		}

		// update costGrad for next layer going backward except
		// for the output layer.
		if l < L {
			costGrad = mat.Dot(net.layers[l+1].weights.T(), deltas[l+1])
		}
	}

	return deltaB, deltaW
}

type layer struct {
	// layer configuration params
	inputSz int            // number of inputs for this layer
	layerSz int            // number of units in this layer
	actFn   Differentiable // activation function
	weights mat.Matrix     // current weights
	biases  mat.Matrix     // current biases

	// states of last forward pass
	weightedSum mat.Matrix // value of 'z'
	activations mat.Matrix // a = g(z)
}
