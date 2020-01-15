package service

import "github.com/google/uuid"

type Generator struct{}

func (g *Generator) Do() (string, error) {
	uuid, err := uuid.NewUUID()
	if err != nil {
		return "", err
	}

	return uuid.String(), nil
}
