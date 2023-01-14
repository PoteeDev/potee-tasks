package shell

import (
	"fmt"
	"strings"
)

func (s *Shell) GetVpnConf(args ...string) string {
	if len(args) > 0 {
		switch args[0] {
		case "get":
			return "your link -> wg0.conf\n"
		default:
			return fmt.Sprintf("bad command: vpn %s\n", strings.Join(args, ""))
		}
	}
	return fmt.Sprintf("bad command: vpn %s\n", strings.Join(args, ""))
}
