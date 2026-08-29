package profile

import (
	"context"
	"time"
)

type GetRankHistoryRequest struct {
	PUUID     string
	StartDate time.Time
	EndDate   time.Time
}

type GetRankHistoryResponse struct {
	Ranks        []RankHistory
	TotalResults int
}

func (s *Service) GetRankHistory(ctx context.Context, req GetRankHistoryRequest) (GetRankHistoryResponse, error) {
	history, err := s.store.GetRankHistoryByPUUID(ctx, req.PUUID, req.StartDate, req.EndDate)
	if err != nil {
		return GetRankHistoryResponse{}, err
	}

	out := make([]RankHistory, 0, len(history))
	for _, item := range history {
		out = append(out, RankHistory{
			Region:       "",
			RankID:       item.ID.String(),
			PUUID:        item.PUUID,
			Date:         item.ValidFrom,
			Tier:         item.Tier,
			Rank:         item.Division,
			LeaguePoints: int32(item.LeaguePoints),
		})
	}

	return GetRankHistoryResponse{Ranks: out, TotalResults: len(out)}, nil
}
