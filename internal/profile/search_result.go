package profile

type SearchResult struct {
	Region        string
	PUUID         string
	Name          string
	Tag           string
	SummonerLevel int
	ProfileIconID string
	Rank          *Rank
}
