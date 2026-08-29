package profile

import (
	"context"
	"time"
)

type GetMatchRequest struct {
	Region  string
	MatchID string
}

type GetMatchResponse struct {
	Region       string
	MatchID      string
	Date         time.Time
	Duration     time.Duration
	Version      string
	WinnerID     string
	Participants []Participant
}

func (s *Service) GetMatch(ctx context.Context, req GetMatchRequest) (GetMatchResponse, error) {
	storeMatch, err := s.store.GetMatchByMatchID(ctx, req.MatchID)
	if err != nil {
		return GetMatchResponse{}, err
	}

	storeParticipants, err := s.store.GetParticipantByMatchID(ctx, req.MatchID)
	if err != nil {
		return GetMatchResponse{}, err
	}

	// Convert storeMatch to GetMatchResponse
	response := GetMatchResponse{
		Region:       storeMatch.Region,
		MatchID:      storeMatch.MatchID,
		Date:         storeMatch.Date,
		Duration:     storeMatch.Duration,
		Version:      storeMatch.Version,
		WinnerID:     storeMatch.WinnerID,
		Participants: make([]Participant, len(storeParticipants)),
	}

	for i, participant := range storeParticipants {
		response.Participants[i] = newParticipant(*participant)
	}

	return response, nil
}
