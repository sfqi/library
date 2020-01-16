package service

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenerator_Do(t *testing.T) {
	require := require.New(t)
	t.Run(("test"), func(t *testing.T) {
		g := &Generator{}
		got, err := g.Do()
		if err != nil {
			t.Errorf("Do() error = %v", err)
			return
		}
		require.Equal(36, len(got))
	})
}
