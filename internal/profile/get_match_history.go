package profile

import (
	"context"
	"time"

	"github.com/rank1zen/kevin/internal/store"
)

type GetMatchHistoryRequest struct {
	PUUID               string
	TimestampRangeStart time.Time
	TimestampRangeEnd   time.Time
	PageSize            int
	PageToken           string
}

type GetMatchHistoryResponse struct {
	Matches       []MatchHistory
	NextPageToken string
}

func (s *Service) GetMatchHistory(ctx context.Context, req GetMatchHistoryRequest) (GetMatchHistoryResponse, error) {
	matches, _, err := s.store.GetMatchByPUUID(ctx, req.PUUID, store.GetMatchByPUUIDPageParams{Size: req.PageSize})
	if err != nil {
		return GetMatchHistoryResponse{}, err
	}

	matchIDs := make([]string, 0, len(matches))
	for _, match := range matches {
		matchIDs = append(matchIDs, match.MatchID)
	}

	participants, err := s.store.GetParticipantByPUUIDAndMatchIDs(ctx, req.PUUID, matchIDs)
	if err != nil {
		return GetMatchHistoryResponse{}, err
	}

	byMatch := make(map[string]*store.Participant, len(participants))
	for _, participant := range participants {
		byMatch[participant.MatchID] = participant
	}

	out := make([]MatchHistory, 0, len(matches))
	for _, match := range matches {
		item := MatchHistory{
			Region:   match.Region,
			MatchID:  match.MatchID,
			PUUID:    req.PUUID,
			Date:     match.Date,
			Duration: match.Duration,
			Version:  match.Version,
			WinnerID: mustAtoi(match.WinnerID),
		}
		if participant, ok := byMatch[match.MatchID]; ok {
			item.TeamID = mustAtoi(participant.TeamID)
			item.ChampionID = mustAtoi(participant.ChampionID)
			item.ChampionLevel = participant.ChampionLevel
			item.TeamPosition = participant.TeamPosition
			item.SummonerSpell = append(item.SummonerSpell, participant.SummonerIDs...)
			item.RunePage = append(item.RunePage, participant.RuneIDs...)
			item.Kills = participant.Kills
			item.Deaths = participant.Deaths
			item.Assists = participant.Assists
		}
		out = append(out, item)
	}

	return GetMatchHistoryResponse{Matches: out, NextPageToken: ""}, nil
}
