package runtime_test

import (
	"context"
	"log"
	"testing"

	"github.com/kevin-labs-org/kevin-lol-service/internal/pgtestcontainer"
)

var DefaultPGInstance *pgtestcontainer.PGInstance

func TestMain(m *testing.M) {
	DefaultPGInstance = pgtestcontainer.NewPGInstance(context.Background())

	defer func(DefaultPGInstance *pgtestcontainer.PGInstance, ctx context.Context) {
		err := DefaultPGInstance.Terminate(ctx)
		if err != nil {
			log.Fatal(err)
		}
	}(DefaultPGInstance, context.Background())

	m.Run()
}
