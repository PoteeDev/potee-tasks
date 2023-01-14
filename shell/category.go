package shell

import (
	"console/models"
	"strings"
)

func (s *Shell) Categories(args ...string) string {
	var categories []string
	s.DB.Model(&models.Challenge{}).Select("DISTINCT ON (category) category").Find(&categories)
	return strings.Join(categories, " ") + "\n"
}
