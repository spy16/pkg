package lua_test

import (
	"testing"

	"github.com/spy16/pkg/lua"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	called := false

	l, err := lua.New(
		lua.Globals(map[string]interface{}{
			"foo": func() {
				called = true
			},
		}),
	)
	assert.NoError(t, err)
	require.NotNil(t, l)
	defer l.Destroy()

	assert.False(t, called)
	err = l.Execute("foo()")
	assert.NoError(t, err)
	assert.True(t, called)
}
