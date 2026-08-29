package profile

import (
	"context"
	"time"

	"uuid"

	"github.com/rank1zen/kevin/internal/store"
)

// Store provides only the needed methods from internal/store. It is used for mocking.
type Store interface {
	// Begin begins a transaction
	Begin(ctx context.Context) (Store, error)

	// Commit commits the transaction
	Commit(ctx context.Context) error

	// Rollback rolls back the transaction
	Rollback(ctx context.Context) error

	GetSummonerByPUUID(ctx context.Context, puuid string) (*store.Summoner, error)
	SearchSummoner(ctx context.Context, name string) ([]*store.Summoner, error)
	GetSummonerByRegionNameTag(ctx context.Context, region, name, tag string) (*store.Summoner, error)
	UpdateSummonerByPUUID(ctx context.Context, puuid string, req store.UpdateSummoner) (*store.Summoner, error)
	CreateSummoner(ctx context.Context, req store.CreateSummoner) (*store.Summoner, error)
	GetRankByPUUID(ctx context.Context, puuid string) (*store.Rank, error)
	CreateRank(ctx context.Context, req store.CreateRank) (*store.Rank, error)
	UpdateRankByPUUID(ctx context.Context, puuid string, req store.UpdateRank) (*store.Rank, error)
	GetRankByPUUIDs(ctx context.Context, puuids []string) ([]*store.Rank, error)
	GetMatchByPUUID(ctx context.Context, puuid string, req store.GetMatchByPUUIDPageParams) ([]*store.Match, store.GetMatchByPUUIDPageParams, error)
	CreateMatch(ctx context.Context, req store.CreateMatch) (*store.Match, error)
	GetMatchByMatchID(ctx context.Context, matchID string) (*store.Match, error)
	GetRankHistoryByPUUID(ctx context.Context, puuid string, start, end time.Time) ([]*store.RankHistory, error)
	CreateRankHistory(ctx context.Context, req store.CreateRankHistory) (*store.RankHistory, error)
	GetRankHistoryByID(ctx context.Context, id uuid.UUID) (*store.RankHistory, error)
	GetParticipantByMatchIDAndPUUID(ctx context.Context, matchID, puuid string) (*store.Participant, error)
	CreateParticipant(ctx context.Context, req store.CreateParticipant) (*store.Participant, error)
	GetParticipantByMatchID(ctx context.Context, matchID string) ([]*store.Participant, error)
	GetParticipantByPUUIDAndMatchIDs(ctx context.Context, puuid string, matchIDs []string) ([]*store.Participant, error)
	AggregateAverageParticipantForChampionIDByPUUID(ctx context.Context, puuid string) ([]*store.AggregateAverageParticipantForChampionID, error)
}
