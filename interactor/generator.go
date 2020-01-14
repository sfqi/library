package interactor

import "github.com/google/uuid"

type generator struct {
}

func (g *generator) GenerateUUID() (string, error) {
	uuid, err := uuid.NewUUID()
	if err != nil {
		return "", err
	}
	return uuid.String(), nil
}
