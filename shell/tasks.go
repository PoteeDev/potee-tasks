package shell

import (
	"console/challanges"
	"console/models"
	"fmt"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
)

func (s *Shell) Tasks(args ...string) string {
	if len(args) == 0 {
		t := table.NewWriter()
		t.SetTitle("Challenges")
		t.AppendRow(BoldTableRow(table.Row{"#", "Name", "Description", "Points"}))
		var chs []models.Challenge
		s.DB.Where("category = ?", s.Pwd).Find(&chs)
		if len(chs) == 0 {
			return fmt.Sprintf("invalid category name: %s\n", s.Pwd)
		}
		for i, ch := range chs {
			t.AppendRow(table.Row{i + 1, ch.Name, ch.Description, ch.Points})
		}
		t.SetStyle(table.StyleRounded)
		return fmt.Sprintf(t.Render() + "\n\r")
	}

	ch := challanges.InitChallenge(s.DB)

	switch args[0] {
	case "get":
		if len(args) >= 2 {
			for _, chName := range args[1:] {
				if err := ch.TakeChallenge(s.Username, chName); err != nil {
					return err.Error() + "\n"
				}
			}
			return fmt.Sprintf("%s added\n\r", strings.Join(args[1:], ","))
		}
		return "challenge name not set"
	case "submit":
		if len(args) == 3 {
			chName, flag := args[1], args[2]
			var exists bool
			s.DB.Model(&models.UsersChallenge{}).
				Select("count(*) > 0").
				Where("challenge_id = (?)", s.DB.Select("id").
					Where("name = ?", chName).
					Table("challenges")).
				Find(&exists)
			if !exists {
				return "wrong challenge name\n\r"
			}
			if err := ch.SubmitFlag(s.Username, chName, flag); err != nil {
				return err.Error() + "\n"
			}
			return fmt.Sprintf("challenge %s solved!\n\r", chName)
		}
		return "bad options\n"
	}
	return ""
}
