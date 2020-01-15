package mock

import (
	"fmt"
	"github.com/stretchr/testify/mock"
)

type Generator struct {
	mock.Mock
}

func (g *Generator) Do() (string, error) {
	args := g.Called()
	if args.Get(0) != "" {
		return fmt.Sprintf("%v", args.Get(0)), nil
	}
	return "", args.Error(1)
}
