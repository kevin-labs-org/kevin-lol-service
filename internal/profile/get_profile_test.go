package profile_test

import (
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/rank1zen/kevin/internal/profile"
	"github.com/rank1zen/kevin/internal/riot"
	"github.com/rank1zen/kevin/internal/store"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestService_GetProfile(t *testing.T) {

	t.Run("should create profile on first visit", func(t *testing.T) {
		mockRiotClient := profile.NewMockRiotClient(t)
		mockStore := profile.NewMockStore(t)

		service := profile.NewService(mockRiotClient, mockStore)

		mockStore.EXPECT().
			GetSummonerByRegionNameTag(mock.Anything, "NA1", "test-name", "").
			Return(nil, pgx.ErrNoRows)

		mockRiotClient.EXPECT().
			GetAccountByRiotID(mock.Anything, "NA1", "test-name", "").
			Return(&riot.Account{
				PUUID:    "test-puuid",
				GameName: "test-name",
			}, nil)

		mockRiotClient.EXPECT().
			GetSummoner(mock.Anything, "NA1", "test-puuid").
			Return(&riot.Summoner{
				PUUID: "test-puuid",
			}, nil)

		mockRiotClient.EXPECT().
			GetLeagueEntriesByPUUID(mock.Anything, "NA1", "test-puuid").
			Return(riot.LeagueList{
				{
					QueueType: riot.QueueTypeRankedSolo5x5,
				},
			}, nil)

		inTxMockStore := profile.NewMockStore(t)

		mockStore.EXPECT().
			Begin(mock.Anything).
			Return(inTxMockStore, nil)

		inTxMockStore.EXPECT().
			CreateSummoner(mock.Anything, mock.Anything).
			Return(&store.Summoner{
				PUUID: "test-puuid",
				Name:  "test-name",
			}, nil)

		inTxMockStore.EXPECT().
			CreateRank(mock.Anything, mock.Anything).
			Return(&store.Rank{
				PUUID: "test-puuid",
			}, nil)

		inTxMockStore.EXPECT().
			Commit(mock.Anything).
			Return(nil)

		inTxMockStore.EXPECT().
			Rollback(mock.Anything).
			Return(pgx.ErrTxClosed)

		got, err := service.GetProfile(t.Context(), profile.GetProfileRequest{
			Region: "NA1",
			Name:   "test-name",
		})

		if assert.NoError(t, err) {
			assert.EqualValues(t, "test-name", got.Name)
		}
	})
}
