package errors

// NewWrapper returns an error handling wrapper. If ptrToResult is a pointer
// to any value, it will be made nil on first error. If a nil-value is passed
// for ptrToResult, this behaviour will be ignored. If ptrToResult is not a
// pointer, this function panics.
func NewWrapper() *Wrapper {
	wr := &Wrapper{}
	return wr
}

// Wrapper wraps multiple function calls as a chain of statements while
// taking care of erros on each function call
type Wrapper struct {
	err error
}

// Return finalizes the result and returns final error. handler can be
// used to transform/wrap the actual error with additional context. If
// nil handler is passed, the actual error will be returned as is.
// Handler will be called only if there was an error
func (wr *Wrapper) Return(handler func(error) error) error {
	if wr.err != nil && handler != nil {
		wr.err = handler(wr.err)
	}

	return wr.err
}

// ReturnOnError stops processing the chain further and returns
// when an error occurs
func (wr *Wrapper) ReturnOnError(fx func() error) *Wrapper {
	if wr.err != nil {
		wr.err = fx()
	}
	return wr
}

// PanicOnError panics when the closure returns error on call
func (wr *Wrapper) PanicOnError(fx func() error) *Wrapper {
	if wr.err != nil {
		wr.err = fx()
		if wr.err != nil {
			panic(wr.err)
		}
	}
	return wr
}
