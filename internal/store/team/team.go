package team

type Team struct {
	ID        int    `db:"id" json:"id"`
	Name      string `db:"name" json:"name"`
	ShortName string `db:"short_name" json:"shortName"`
	CreatedAt string `db:"created_at" json:"createdAt"`
	UpdatedAt string `db:"updated_at" json:"updatedAt"`
}

type TeamForm struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	ShortName string `json:"shortName"`
	UpdatedAt string `json:"updatedAt"`
}

func NewTeamForm(team Team) TeamForm {
	return TeamForm{
		ID:        team.ID,
		Name:      team.Name,
		ShortName: team.ShortName,
		UpdatedAt: team.UpdatedAt,
	}
}

type TeamStore interface {
	CreateTeam(team Team) error
	ListTeams() ([]Team, error)
	ReadTeam(id int) (Team, error)
	DeleteTeam(id int) error
	UpdateTeam(team Team) error
}
