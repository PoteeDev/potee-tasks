package shell

import (
	"encoding/json"
	"log"
	"strings"

	"gorm.io/gorm"
)

type Shell struct {
	DB       *gorm.DB
	Pwd      string
	Username string
	Commands *Commands
}

func InitShell(db *gorm.DB, user string) *Shell {
	s := Shell{
		DB:       db,
		Pwd:      "~",
		Username: user,
	}
	c := Commands{AllCommands: make(map[string]func(...string) string)}
	c.AddCommand("ls", s.Categories)
	c.AddCommand("tasks", s.Tasks)
	c.AddCommand("help", s.Help)
	c.AddCommand("vpn", s.GetVpnConf)
	c.AddCommand("whoami", s.Profile)
	c.AddCommand("who", s.Profile)
	c.AddCommand("w", s.Profile)
	c.FromFile("commands.yml")

	s.Commands = &c
	return &s
}

type Request struct {
	Pwd string `json:"pwd"`
	Cmd string `json:"cmd"`
}

type Response struct {
	Pwd    string `json:"pwd"`
	Answer string `json:"answer"`
}

func (s *Shell) Execute(cmd []byte) []byte {
	var req Request
	err := json.Unmarshal(cmd, &req)
	if err != nil {
		log.Print("unmarshal:", err)
		return []byte("")
	}
	s.Pwd = req.Pwd
	cmdList := strings.Split(req.Cmd, " ")
	var answer string

	if _, ok := s.Commands.AllCommands[cmdList[0]]; ok {
		answer = s.Commands.AllCommands[cmdList[0]](cmdList[1:]...)
	} else {
		answer = cmdList[0] + ": command not found\n"
	}

	// switch cmdList[0] {
	// case "cd":
	// 	s.Pwd = cmdList[1]
	// 	answer = ""
	// default:
	// 	answer = cmdList[0] + ": command not found\n"
	// }

	resp := Response{
		Pwd:    s.Pwd,
		Answer: answer,
	}
	result, _ := json.Marshal(resp)
	if err != nil {
		log.Print("marshal:", err)
		return []byte("")
	}
	return result
}
