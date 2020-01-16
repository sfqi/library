package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerator_Do(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)
	t.Run(("test"), func(t *testing.T) {
		g := &Generator{}
		got, err := g.Do()

		require.NoError(err, "Do() error")
		assert.Equal(36, len(got))
	})
}
