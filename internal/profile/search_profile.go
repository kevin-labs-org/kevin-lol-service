package profile

import (
	"context"

	"github.com/rank1zen/kevin/internal/store"
)

type SearchProfileRequest struct {
	Query string
}

type SearchProfileResponse struct {
	Results []*SearchResult
}

func (s *Service) SearchProfile(ctx context.Context, req SearchProfileRequest) (SearchProfileResponse, error) {
	results, err := s.store.SearchSummoner(ctx, req.Query)
	if err != nil {
		return SearchProfileResponse{}, err
	}

	if len(results) == 0 {
		return SearchProfileResponse{}, nil
	}

	puuidList := make([]string, len(results))
	for _, result := range results {
		puuidList = append(puuidList, result.PUUID)
	}

	puuidToRank := make(map[string]*store.Rank)
	ranks, err := s.store.GetRankByPUUIDs(ctx, puuidList)
	if err != nil {
		return SearchProfileResponse{}, err
	}

	for _, rank := range ranks {
		puuidToRank[rank.PUUID] = rank
	}

	var searchResults []*SearchResult
	for _, result := range results {
		storeRank, ok := puuidToRank[result.PUUID]
		if !ok {
			panic("rank with puuid not found in store")
		}

		searchResults = append(searchResults, &SearchResult{
			Region:        result.Region,
			PUUID:         result.PUUID,
			Name:          result.Name,
			Tag:           result.Tagline,
			SummonerLevel: result.SummonerLevel,
			ProfileIconID: result.ProfileIconID,
			Rank: &Rank{
				Tier:         storeRank.Tier,
				Rank:         storeRank.Division,
				LeaguePoints: storeRank.LeaguePoints,
			},
		})
	}

	return SearchProfileResponse{
		Results: searchResults,
	}, nil
}
