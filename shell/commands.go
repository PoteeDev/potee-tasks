package shell

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"gopkg.in/yaml.v3"
)

type Commands struct {
	AllCommands map[string]func(...string) string
}

func (c *Commands) AddCommand(name string, function func(...string) string) {
	c.AllCommands[name] = function
}

func (c *Commands) FromFile(filename string) {
	yfile, err := ioutil.ReadFile(filename)

	if err != nil {
		log.Fatal(err)
	}

	data := make(map[string]string)
	if err = yaml.Unmarshal(yfile, &data); err != nil {

		log.Fatal(err)
	}
	for k, v := range data {
		function := func(v string) func(s ...string) string {
			return func(s ...string) string {
				v = strings.Replace(v, "\n", "\n\r", -1)
				return fmt.Sprintf("%s\n", v)
			}
		}(v)
		c.AddCommand(k, function)
	}
	log.Println(c.AllCommands)
}

func (c *Commands) Run(name string, args ...string) string {
	return c.AllCommands[name](args...)
}
