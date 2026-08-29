package profile_test

import (
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/rank1zen/kevin/internal/profile"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestService_GetProfile(t *testing.T) {

	t.Run("should create profile on first visit", func(t *testing.T) {
		mockStore := profile.NewMockStore(t)

		service := profile.NewService(nil, mockStore)

		mockStore.EXPECT().
			GetSummonerByRegionNameTag(mock.Anything, "NA1", "test-name", "").
			Return(nil, pgx.ErrNoRows)

		got, err := service.GetProfile(t.Context(), profile.GetProfileRequest{
			Region: "NA1",
			Name:   "test-name",
		})

		if assert.NoError(t, err) {
			assert.EqualValues(t, "test-name", got.Name)
		}
	})
}
