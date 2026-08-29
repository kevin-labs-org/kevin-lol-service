package profile

import (
	"context"
	"time"

	"uuid"

	"github.com/jackc/pgx/v5"
	"github.com/rank1zen/kevin/internal/riot"
	"github.com/rank1zen/kevin/internal/store"
)

type Service struct {
	riot *riot.Client

	store Store
}

func NewService(
	riot *riot.Client,
	store Store,
) *Service {
	return &Service{
		riot:  riot,
		store: store,
	}
}

type ServiceStore struct {
	tx pgx.Tx
}

// NewServiceStore creates a service store. Provide *pgxpool.Pool for concurrent db access.
func NewServiceStore(tx pgx.Tx) *ServiceStore {
	return &ServiceStore{
		tx: tx,
	}
}

func (s *ServiceStore) Begin(ctx context.Context) (Store, error) {
	tx, err := s.tx.Begin(ctx)
	if err != nil {
		return nil, err
	}
	return NewServiceStore(tx), nil
}

func (s *ServiceStore) Commit(ctx context.Context) error {
	return s.tx.Commit(ctx)
}

func (s *ServiceStore) Rollback(ctx context.Context) error {
	return s.tx.Rollback(ctx)
}

func (s *ServiceStore) GetSummonerByPUUID(ctx context.Context, puuid string) (*store.Summoner, error) {
	return store.NewSummonerStore(s.tx).GetSummonerByPUUID(ctx, puuid)
}

func (s *ServiceStore) SearchSummoner(ctx context.Context, name string) ([]*store.Summoner, error) {
	return store.NewSummonerStore(s.tx).SearchSummoner(ctx, name)
}

func (s *ServiceStore) GetSummonerByRegionNameTag(ctx context.Context, region, name, tag string) (*store.Summoner, error) {
	return store.NewSummonerStore(s.tx).GetSummonerByRegionNameTag(ctx, region, name, tag)
}

func (s *ServiceStore) UpdateSummonerByPUUID(ctx context.Context, puuid string, req store.UpdateSummoner) (*store.Summoner, error) {
	return store.NewSummonerStore(s.tx).UpdateSummonerByPUUID(ctx, puuid, req)
}

func (s *ServiceStore) CreateSummoner(ctx context.Context, req store.CreateSummoner) (*store.Summoner, error) {
	return store.NewSummonerStore(s.tx).CreateSummoner(ctx, req)
}

func (s *ServiceStore) GetRankByPUUID(ctx context.Context, puuid string) (*store.Rank, error) {
	return store.NewRankStore(s.tx).GetRankByPUUID(ctx, puuid)
}

func (s *ServiceStore) GetRankByPUUIDs(ctx context.Context, puuids []string) ([]*store.Rank, error) {
	return store.NewRankStore(s.tx).GetRankByPUUIDs(ctx, puuids)
}

func (s *ServiceStore) CreateRank(ctx context.Context, req store.CreateRank) (*store.Rank, error) {
	return store.NewRankStore(s.tx).CreateRank(ctx, req)
}

func (s *ServiceStore) UpdateRankByPUUID(ctx context.Context, puuid string, req store.UpdateRank) (*store.Rank, error) {
	return store.NewRankStore(s.tx).UpdateRankByPUUID(ctx, puuid, req)
}

func (s *ServiceStore) GetMatchByPUUID(ctx context.Context, puuid string, req store.GetMatchByPUUIDPageParams) ([]*store.Match, store.GetMatchByPUUIDPageParams, error) {
	return store.NewMatchStore(s.tx).GetMatchByPUUID(ctx, puuid, req)
}

func (s *ServiceStore) CreateMatch(ctx context.Context, req store.CreateMatch) (*store.Match, error) {
	return store.NewMatchStore(s.tx).CreateMatch(ctx, req)
}

func (s *ServiceStore) GetMatchByMatchID(ctx context.Context, matchID string) (*store.Match, error) {
	return store.NewMatchStore(s.tx).GetMatchByMatchID(ctx, matchID)
}

func (s *ServiceStore) GetRankHistoryByPUUID(ctx context.Context, puuid string, start, end time.Time) ([]*store.RankHistory, error) {
	return store.NewRankHistoryStore(s.tx).GetRankHistoryByPUUID(ctx, puuid, start, end)
}

func (s *ServiceStore) CreateRankHistory(ctx context.Context, req store.CreateRankHistory) (*store.RankHistory, error) {
	return store.NewRankHistoryStore(s.tx).CreateRankHistory(ctx, req)
}

func (s *ServiceStore) GetRankHistoryByID(ctx context.Context, id uuid.UUID) (*store.RankHistory, error) {
	return store.NewRankHistoryStore(s.tx).GetRankHistoryByID(ctx, id)
}

func (s *ServiceStore) GetParticipantByMatchIDAndPUUID(ctx context.Context, matchID, puuid string) (*store.Participant, error) {
	return store.NewParticipantStore(s.tx).GetParticipantByMatchIDAndPUUID(ctx, matchID, puuid)
}

func (s *ServiceStore) CreateParticipant(ctx context.Context, req store.CreateParticipant) (*store.Participant, error) {
	return store.NewParticipantStore(s.tx).CreateParticipant(ctx, req)
}

func (s *ServiceStore) GetParticipantByMatchID(ctx context.Context, matchID string) ([]*store.Participant, error) {
	return store.NewParticipantStore(s.tx).GetParticipantByMatchID(ctx, matchID)
}

func (s *ServiceStore) GetParticipantByPUUIDAndMatchIDs(ctx context.Context, puuid string, matchIDs []string) ([]*store.Participant, error) {
	return store.NewParticipantStore(s.tx).GetParticipantByPUUIDAndMatchIDs(ctx, puuid, matchIDs)
}

func (s *ServiceStore) AggregateAverageParticipantForChampionIDByPUUID(ctx context.Context, puuid string) ([]*store.AggregateAverageParticipantForChampionID, error) {
	return store.NewParticipantStore(s.tx).AggregateAverageParticipantForChampionIDByPUUID(ctx, puuid)
}
