package profile

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/rank1zen/kevin/internal/riot"
	"github.com/rank1zen/kevin/internal/store"
)

type GetProfileRequest struct {
	Region string
	Name   string
	Tag    string
}

type GetProfileResponse struct {
	Region        string
	PUUID         string
	Name          string
	Tag           string
	SummonerLevel int
	ProfileIconID string
	Rank          *Rank
}

func (s *Service) GetProfile(ctx context.Context, req GetProfileRequest) (GetProfileResponse, error) {
	summoner, err := s.store.GetSummonerByRegionNameTag(ctx, req.Region, req.Name, req.Tag)
	if errors.Is(err, pgx.ErrNoRows) {
		account, err := s.riotClient.GetAccountByRiotID(ctx, req.Region, req.Name, req.Tag)
		if err != nil {
			return GetProfileResponse{}, err
		}

		return s.firstVisitProfile(ctx, req.Region, account.GameName, account.TagLine, account.PUUID)
	}

	storeRank, err := s.store.GetRankByPUUID(ctx, summoner.PUUID)
	if err != nil {
		return GetProfileResponse{}, err
	}

	return GetProfileResponse{
		Region:        summoner.Region,
		PUUID:         summoner.PUUID,
		Name:          summoner.Name,
		Tag:           summoner.Tagline,
		SummonerLevel: summoner.SummonerLevel,
		ProfileIconID: summoner.ProfileIconID,
		Rank: &Rank{
			Tier:         storeRank.Tier,
			Rank:         storeRank.Division,
			LeaguePoints: storeRank.LeaguePoints,
		},
	}, nil
}

func (s *Service) firstVisitProfile(ctx context.Context, region, name, tag, puuid string) (GetProfileResponse, error) {
	riotSummoner, err := s.riotClient.GetSummoner(ctx, "NA1", puuid)
	if err != nil {
		return GetProfileResponse{}, err
	}

	summonerUpdateTS := time.Now()

	riotLeague, err := s.riotClient.GetLeagueEntriesByPUUID(ctx, "NA1", puuid)
	if err != nil {
		return GetProfileResponse{}, err
	}

	rankUpdateTS := time.Now()

	soloQEntry, rankFound := findSoloQ(riotLeague)

	inTx, err := s.store.Begin(ctx)
	if err != nil {
		return GetProfileResponse{}, err
	}
	defer func() {
		_ = inTx.Rollback(ctx)
	}()

	dbSummoner, err := inTx.CreateSummoner(ctx, store.CreateSummoner{
		PUUID:         riotSummoner.PUUID,
		Region:        region,
		Name:          name,
		Tagline:       tag,
		SummonerLevel: int(riotSummoner.SummonerLevel),
		ProfileIconID: strconv.Itoa(riotSummoner.ProfileIconID),
		LastUpdated:   summonerUpdateTS,
	})
	if err != nil {
		return GetProfileResponse{}, err
	}

	var createRank *store.CreateRank
	if rankFound {
		createRank = &store.CreateRank{
			PUUID:        riotSummoner.PUUID,
			Wins:         soloQEntry.Wins,
			Losses:       soloQEntry.Losses,
			Tier:         soloQEntry.Tier,
			Division:     soloQEntry.Division,
			LeaguePoints: soloQEntry.LeaguePoints,
			LastUpdated:  rankUpdateTS,
		}
	} else {
		createRank = &store.CreateRank{
			PUUID:       riotSummoner.PUUID,
			LastUpdated: rankUpdateTS,
		}

	}

	dbRank, err := inTx.CreateRank(ctx, *createRank)
	if err != nil {
		return GetProfileResponse{}, err
	}

	err = inTx.Commit(ctx)
	if err != nil {
		return GetProfileResponse{}, err
	}

	return GetProfileResponse{
		Region:        dbSummoner.Region,
		PUUID:         dbSummoner.PUUID,
		Name:          dbSummoner.Name,
		Tag:           dbSummoner.Tagline,
		SummonerLevel: dbSummoner.SummonerLevel,
		ProfileIconID: dbSummoner.ProfileIconID,
		Rank: &Rank{
			Tier:         dbRank.Tier,
			Rank:         dbRank.Division,
			LeaguePoints: dbRank.LeaguePoints,
		},
	}, nil
}

func findSoloQ(leagueEntries []riot.LeagueEntry) (riot.LeagueEntry, bool) {
	for _, entry := range leagueEntries {
		if entry.QueueType == riot.QueueTypeRankedSolo5x5 {
			return entry, true
		}
	}

	return riot.LeagueEntry{}, false
}
