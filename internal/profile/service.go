package profile

import (
	"github.com/kevin-labs-org/kevin-lol-service/internal/riot"
)

type Service struct {
	riot  *riot.Client
	store Store
}

func NewProfileService(riot *riot.Client, store Store) *Service {
	return &Service{
		riot:  riot,
		store: store,
	}
}
