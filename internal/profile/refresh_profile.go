package profile

import (
	"context"
	"strconv"
	"time"

	"github.com/rank1zen/kevin/internal/store"
)

type RefreshProfileRequest struct {
	PUUID string
}

type RefreshProfileResponse struct {
	Region        string
	PUUID         string
	Name          string
	Tag           string
	SummonerLevel int
	ProfileIconID string
	Rank          *Rank
}

func newRefreshProfileResponse(summoner *store.Summoner, rank *store.Rank) RefreshProfileResponse {
	return RefreshProfileResponse{
		Region:        summoner.Region,
		PUUID:         summoner.PUUID,
		Name:          summoner.Name,
		Tag:           summoner.Tagline,
		SummonerLevel: summoner.SummonerLevel,
		ProfileIconID: summoner.ProfileIconID,
		Rank: &Rank{
			Tier:         rank.Tier,
			Rank:         rank.Division,
			LeaguePoints: rank.LeaguePoints,
		},
	}
}

func (s *Service) RefreshProfile(ctx context.Context, req RefreshProfileRequest) (RefreshProfileResponse, error) {
	account, err := s.riotClient.GetAccountByPUUID(ctx, "a", req.PUUID)
	if err != nil {
		return RefreshProfileResponse{}, err
	}

	leagueEntries, err := s.riotClient.GetLeagueEntriesByPUUID(ctx, "a", req.PUUID)
	if err != nil {
		return RefreshProfileResponse{}, err
	}

	rankUpdateTS := time.Now()

	soloq, found := findSoloQ(leagueEntries)

	summoner, err := s.riotClient.GetSummoner(ctx, "a", req.PUUID)
	if err != nil {
		return RefreshProfileResponse{}, err
	}

	summonerUpdateTS := time.Now()

	inTx, err := s.store.Begin(ctx)
	if err != nil {
		return RefreshProfileResponse{}, err
	}
	defer func() {
		_ = inTx.Rollback(ctx)
	}()

	storeSummoner, err := inTx.UpdateSummonerByPUUID(ctx, req.PUUID, store.UpdateSummoner{
		Name:          account.GameName,
		Tagline:       account.TagLine,
		SummonerLevel: int(summoner.SummonerLevel),
		ProfileIconID: strconv.Itoa(summoner.ProfileIconID),
		LastUpdated:   summonerUpdateTS,
	})
	if err != nil {
		return RefreshProfileResponse{}, err
	}

	var updateRank *store.UpdateRank
	if !found {
		updateRank = &store.UpdateRank{}
	} else {
		updateRank = &store.UpdateRank{
			Wins:         soloq.Wins,
			Losses:       soloq.Losses,
			Tier:         soloq.Tier,
			Division:     soloq.Division,
			LeaguePoints: soloq.LeaguePoints,
			LastUpdated:  rankUpdateTS,
		}
	}

	storeRank, err := inTx.UpdateRankByPUUID(ctx, req.PUUID, *updateRank)
	if err != nil {
		return RefreshProfileResponse{}, err
	}

	err = inTx.Commit(ctx)
	if err != nil {
		return RefreshProfileResponse{}, err
	}

	return newRefreshProfileResponse(storeSummoner, storeRank), nil
}
