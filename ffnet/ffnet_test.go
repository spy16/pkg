package ffnet_test

import (
	"testing"

	"github.com/spy16/pkg/ffnet"
)

func BenchmarkFFNet_Predict_2_10(b *testing.B) {
	net, err := ffnet.New(2, ffnet.Layer(10, ffnet.Sigmoid()))
	if err != nil {
		panic(err)
	}
	for i := 0; i < b.N; i++ {
		_, _ = net.Predict(1, 1)
	}
}

func BenchmarkFFNet_Predict_2_100_500_100(b *testing.B) {
	net, err := ffnet.New(2,
		ffnet.Layer(100, ffnet.Sigmoid()),
		ffnet.Layer(500, ffnet.Sigmoid()),
		ffnet.Layer(100, ffnet.Sigmoid()))
	if err != nil {
		panic(err)
	}
	for i := 0; i < b.N; i++ {
		_, _ = net.Predict(1, 1)
	}
}
