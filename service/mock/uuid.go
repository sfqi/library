package mock

import (
	"github.com/stretchr/testify/mock"
)

type Generator struct {
	mock.Mock
}

func (g *Generator) Do() (string, error) {
	args := g.Called()
	if args.Get(0) != "" {
		return args.String(0), nil
	}

	return "", args.Error(1)
}
