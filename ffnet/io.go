package ffnet

import (
	"io"
)

// Load reads FFNet from a reader and returns. Stream must have been serialized
// using Save().
func Load(r io.Reader) (*FFNet, error) {
	panic("implement me")
}

// Save writes the given network to a writer using a custom serialization which
// can be read back using Load().
func Save(w io.Writer, net *FFNet) error {
	panic("implement me")
}
