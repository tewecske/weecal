package passwordhash

type TestPasswordHash struct {
	Match          bool
	HashedPassword string
	Error          error
}

func (p *TestPasswordHash) ComparePasswordAndHash(password string, encodedHash string) (match bool, err error) {
	if p.Error != nil {
		return false, p.Error
	} else {
		return p.Match, nil
	}
}
func (p *TestPasswordHash) GenerateFromPassword(password string) (encodedHash string, err error) {
	if p.Error != nil {
		return "", p.Error
	} else {
		return p.HashedPassword, nil
	}
}
