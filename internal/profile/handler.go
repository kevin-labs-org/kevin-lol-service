package profile

import (
	"context"
	"strconv"
	"time"

	"buf.build/gen/go/kevin-labs/lol-service/connectrpc/go/kevin/lolservice/v1/lolservicev1connect"
	lolservicev1 "buf.build/gen/go/kevin-labs/lol-service/protocolbuffers/go/kevin/lolservice/v1"
	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Handler struct {
	lolservicev1connect.UnimplementedProfileServiceHandler

	service *Service
}

func (s Handler) GetProfile(ctx context.Context, c *connect.Request[lolservicev1.GetProfileRequest]) (*connect.Response[lolservicev1.GetProfileResponse], error) {
	resp, err := s.service.GetProfile(ctx, GetProfileRequest{
		Region: c.Msg.Region,
		Name:   c.Msg.Name,
		Tag:    c.Msg.Tag,
	})
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&lolservicev1.GetProfileResponse{
		Region:        resp.Region,
		Puuid:         resp.PUUID,
		Name:          resp.Name,
		Tag:           resp.Tag,
		SummonerLevel: int32(resp.SummonerLevel),
		ProfileIconId: int32(mustAtoi(resp.ProfileIconID)),
		Rank:          toProtoRank(resp.Rank),
	}), nil
}

func (s Handler) RefreshProfile(ctx context.Context, c *connect.Request[lolservicev1.RefreshProfileRequest]) (*connect.Response[lolservicev1.RefreshProfileResponse], error) {
	resp, err := s.service.RefreshProfile(ctx, RefreshProfileRequest{PUUID: c.Msg.Puuid})
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&lolservicev1.RefreshProfileResponse{
		Region:        resp.Region,
		Puuid:         resp.PUUID,
		Name:          resp.Name,
		Tag:           resp.Tag,
		SummonerLevel: int32(resp.SummonerLevel),
		ProfileIconId: int32(mustAtoi(resp.ProfileIconID)),
		Rank:          toProtoRank(resp.Rank),
	}), nil
}

func (s Handler) SearchProfile(ctx context.Context, c *connect.Request[lolservicev1.SearchProfileRequest]) (*connect.Response[lolservicev1.SearchProfileResponse], error) {
	resp, err := s.service.SearchProfile(ctx, SearchProfileRequest{Query: c.Msg.Query})
	if err != nil {
		return nil, err
	}
	results := make([]*lolservicev1.SearchResult, 0, len(resp.Results))
	for _, r := range resp.Results {
		results = append(results, &lolservicev1.SearchResult{
			Region:        r.Region,
			Puuid:         r.PUUID,
			Name:          r.Name,
			Tag:           r.Tag,
			SummonerLevel: int32(r.SummonerLevel),
			ProfileIconId: int32(mustAtoi(r.ProfileIconID)),
			Rank:          toProtoRank(r.Rank),
		})
	}
	return connect.NewResponse(&lolservicev1.SearchProfileResponse{Results: results}), nil
}

func (s Handler) GetMatch(ctx context.Context, c *connect.Request[lolservicev1.GetMatchRequest]) (*connect.Response[lolservicev1.GetMatchResponse], error) {
	resp, err := s.service.GetMatch(ctx, GetMatchRequest{Region: c.Msg.Region, MatchID: c.Msg.MatchId})
	if err != nil {
		return nil, err
	}
	participants := make([]*lolservicev1.Participant, 0, len(resp.Participants))
	for _, p := range resp.Participants {
		participants = append(participants, &lolservicev1.Participant{
			Puuid:         p.PUUID,
			MatchId:       p.MatchID,
			TeamId:        int32(mustAtoi(p.TeamID)),
			ParticipantId: int32(mustAtoi(p.ParticipantID)),
			ChampionId:    int32(mustAtoi(p.ChampionID)),
			ChampionLevel: int32(p.ChampionLevel),
			TeamPosition:  p.TeamPosition,
		})
	}
	return connect.NewResponse(&lolservicev1.GetMatchResponse{
		Region:       resp.Region,
		MatchId:      resp.MatchID,
		Date:         timestamppb.New(resp.Date),
		Duration:     durationpb.New(resp.Duration),
		Version:      resp.Version,
		WinnerId:     int32(mustAtoi(resp.WinnerID)),
		Participants: participants,
	}), nil
}

func (s Handler) GetMatchHistory(ctx context.Context, c *connect.Request[lolservicev1.GetMatchHistoryRequest]) (*connect.Response[lolservicev1.GetMatchHistoryResponse], error) {
	start := time.Time{}
	if c.Msg.TimestampRangeStart != nil {
		start = c.Msg.TimestampRangeStart.AsTime()
	}
	end := time.Time{}
	if c.Msg.TimestampRangeEnd != nil {
		end = c.Msg.TimestampRangeEnd.AsTime()
	}
	resp, err := s.service.GetMatchHistory(ctx, GetMatchHistoryRequest{
		PUUID:               c.Msg.Puuid,
		TimestampRangeStart: start,
		TimestampRangeEnd:   end,
		PageSize:            int(c.Msg.PageSize),
		PageToken:           c.Msg.PageToken,
	})
	if err != nil {
		return nil, err
	}
	matches := make([]*lolservicev1.MatchHistory, 0, len(resp.Matches))
	for _, m := range resp.Matches {
		matches = append(matches, &lolservicev1.MatchHistory{
			Region:        m.Region,
			MatchId:       m.MatchID,
			Puuid:         m.PUUID,
			Date:          timestamppb.New(m.Date),
			Duration:      durationpb.New(m.Duration),
			Version:       m.Version,
			WinnerId:      int32(m.WinnerID),
			ChampionId:    int32(m.ChampionID),
			ChampionLevel: int32(m.ChampionLevel),
			TeamPosition:  m.TeamPosition,
			Kills:         int32(m.Kills),
			Deaths:        int32(m.Deaths),
			Assists:       int32(m.Assists),
			Rank:          toProtoRank(m.Rank),
		})
	}
	return connect.NewResponse(&lolservicev1.GetMatchHistoryResponse{
		Matches:       matches,
		NextPageToken: resp.NextPageToken,
	}), nil
}

func (s Handler) GetRankHistory(ctx context.Context, c *connect.Request[lolservicev1.GetRankHistoryRequest]) (*connect.Response[lolservicev1.GetRankHistoryResponse], error) {
	start := time.Time{}
	if c.Msg.StartDate != nil {
		start = c.Msg.StartDate.AsTime()
	}
	end := time.Time{}
	if c.Msg.EndDate != nil {
		end = c.Msg.EndDate.AsTime()
	}
	resp, err := s.service.GetRankHistory(ctx, GetRankHistoryRequest{
		PUUID:     c.Msg.Puuid,
		StartDate: start,
		EndDate:   end,
	})
	if err != nil {
		return nil, err
	}
	ranks := make([]*lolservicev1.RankHistory, 0, len(resp.Ranks))
	for _, r := range resp.Ranks {
		ranks = append(ranks, &lolservicev1.RankHistory{
			Puuid:        r.PUUID,
			Date:         timestamppb.New(r.Date),
			Tier:         r.Tier,
			Rank:         r.Rank,
			LeaguePoints: int32(r.LeaguePoints),
		})
	}
	return connect.NewResponse(&lolservicev1.GetRankHistoryResponse{
		Ranks:        ranks,
		TotalResults: int32(resp.TotalResults),
	}), nil
}

func (s Handler) GetChampionAggregate(ctx context.Context, c *connect.Request[lolservicev1.GetChampionAggregateRequest]) (*connect.Response[lolservicev1.GetChampionAggregateResponse], error) {
	_, err := s.service.GetChampionAggregate(ctx, GetChampionAggregateRequest{
		Region:  "",
		MatchID: "",
	})
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&lolservicev1.GetChampionAggregateResponse{}), nil
}

func toProtoRank(r *Rank) *lolservicev1.Rank {
	if r == nil {
		return nil
	}
	return &lolservicev1.Rank{
		Tier:         r.Tier,
		Rank:         r.Rank,
		LeaguePoints: int32(r.LeaguePoints),
	}
}

func mustAtoi(s string) int {
	v, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return v
}
