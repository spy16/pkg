package ffnet

import (
	"context"
	"errors"
	"math/rand"
	"time"

	"github.com/spy16/pkg/ffnet/mat"
)

// SGDTrainer implements Stochastic Gradient Descent trainer for FFNet
// using BackPropagation algorithm.
type SGDTrainer struct {
	// Embed the network that needs to be trained.
	*FFNet
	// Eta is the learning rate to be used for training.
	Eta float64
	// Loss is the loss function to use for computing loss during training.
	Loss LossFunc
	// log function to be used for logging progress. Nil means no logs.
	LogFunc func(msg string, args ...interface{})
}

// Example represents a single training sample.
type Example struct {
	Inputs  []float64
	Outputs []float64
}

// Train runs training iterations for given number of epochs. Training loop
// can be stopped by cancelling the context.
func (t SGDTrainer) Train(ctx context.Context, epochs int, samples []Example) error {
	if err := t.init(); err != nil {
		return err
	}

	trainingStart := time.Now()
	for i := 0; i < epochs; i++ {
		startedAt := time.Now()
		shuffle(samples)

		for _, sample := range samples {
			x := mat.From(t.inputSz, 1, sample.Inputs...)
			y := mat.From(t.outputSz, 1, sample.Outputs...)

			zs, as, err := t.forwardPass(sample.Inputs)
			if err != nil {
				return err
			}

			yHat := zs[len(zs)-1]
			costGrad := t.Loss.FPrime(y, yHat)

			deltaB, deltaW := t.backPropagate(zs, as, x, costGrad)
			for i := 0; i < len(t.layers); i++ {
				t.layers[i].weights = mat.Sub(t.layers[i].weights, deltaW[i].Scale(t.Eta))
				t.layers[i].biases = mat.Sub(t.layers[i].biases, deltaB[i].Scale(t.Eta))
			}
		}

		t.LogFunc("epoch %d finished in %s", i, time.Since(startedAt))
	}

	t.LogFunc("training run of %d epochs finished in %s", epochs, time.Since(trainingStart))
	return nil
}

func (t *SGDTrainer) init() error {
	if t.FFNet == nil {
		return errors.New("field Net is not set, nothing to train")
	}

	if t.Loss.F == nil {
		t.Loss = SquaredError()
	}
	if t.Eta == 0 {
		t.Eta = 0.5
	}
	if t.LogFunc == nil {
		t.LogFunc = func(_ string, _ ...interface{}) {}
	}
	return nil
}

func shuffle(samples []Example) {
	for i := range samples {
		j := rand.Intn(i + 1)
		samples[i], samples[j] = samples[j], samples[i]
	}
}
