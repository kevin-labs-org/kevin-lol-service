package profile_test

import (
	"testing"

	"github.com/rank1zen/kevin/internal/profile"
	"github.com/stretchr/testify/assert"
)

func TestService_GetProfile(t *testing.T) {

	service := profile.NewService(nil, nil)

	t.Run("should create profile on first visit", func(t *testing.T) {
		got, err := service.GetProfile(t.Context(), profile.GetProfileRequest{
			Region: "NA1",
			Name:   "test-name",
		})

		if assert.NoError(t, err) {
			assert.EqualValues(t, "test-name", got.Name)
		}
	})
}
