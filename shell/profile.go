package shell

import (
	"console/models"
	"fmt"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

func (s *Shell) Profile(args ...string) string {
	var profile models.User
	result := s.DB.Model(&models.User{}).
		Preload("Group").
		Preload("Pool").
		Preload("UsersChallenges.Challenge").
		First(&profile, "login = ?", s.Username)

	if result.Error != nil {
		return result.Error.Error() + "\n"
	}

	t := table.NewWriter()
	// t.SetTitle("Profile")
	t.AppendRow(BoldTableRow(table.Row{"Profile", "Info"}), table.RowConfig{AutoMerge: true})
	t.AppendSeparator()
	t.AppendRow(table.Row{fmt.Sprintf("Login:\t%s", profile.Login), fmt.Sprintf("CIDR: %s", profile.Pool.Cidr)})
	t.AppendRow(table.Row{fmt.Sprintf("Name: \t%s %s", profile.SecondName, profile.FirstName)})
	t.AppendRow(table.Row{fmt.Sprintf("Group:\t%s", profile.Group.GroupCode)})
	t.SetStyle(table.StyleRounded)

	var score = 0
	t2 := table.NewWriter()
	t2.SetTitle("Challenges")
	if len(profile.UsersChallenges) > 0 {
		t2.AppendRow(BoldTableRow(table.Row{"Name", "Desctiption", "Address", "Solved"}))
		for _, ch := range profile.UsersChallenges {
			var ip string
			for i, addr := range profile.Pool.IPPool[1:] {
				if i+1 == ch.Challenge.Container {
					ip = addr
					break
				}
			}
			t2.AppendRow(table.Row{ch.Challenge.Name, ch.Challenge.Description, ip, ch.Solved})
			if ch.Solved {
				score += ch.Challenge.Points
			}
		}
	} else {
		t2.AppendRow(table.Row{"no challenge taken"})
	}

	t2.SetStyle(table.StyleRounded)
	t2.SetColumnConfigs([]table.ColumnConfig{{Number: 4, Transformer: solvedColor}})
	return fmt.Sprintf("%s\nScore: %d\n%s\n", t.Render(), score, t2.Render())
}

func BoldTableRow(rows table.Row) table.Row {
	for i, row := range rows {
		rows[i] = text.Bold.Sprint(row)
	}
	return rows
}

var solvedColor = text.Transformer(func(val interface{}) string {
	switch v := val.(type) {
	case bool:
		if v {
			return text.Colors{text.FgHiGreen}.Sprint("Yes")
		}
		return text.Colors{text.FgHiRed}.Sprint("No")

	default:
		return val.(string)
	}
})
