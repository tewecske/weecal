package team

type Team struct {
	ID        int    `db:"id" json:"id"`
	Name      string `db:"name" json:"name"`
	ShortName string `db:"short_name" json:"shortName"`
}

type TeamForm struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	ShortName string `json:"shortName"`
}

type TeamStore interface {
	CreateTeam(team Team) error
	ListTeams() ([]Team, error)
	ReadTeam(id string) (Team, error)
	DeleteTeam(id string) error
}
