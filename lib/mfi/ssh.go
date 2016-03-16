package mfi

import (
	"os"
	"os/exec"
	"strings"
)

type SshCommands struct {
	cmds []string
}

func NewSshCommands() (s *SshCommands) {
	s = &SshCommands{}
	s.cmds = make([]string, 0, 100)
	return
}

func (s *SshCommands) Add(cmd string) {
	s.cmds = append(s.cmds, cmd)
}

func (s *SshCommands) AddFile(filename, text string) {
	lines := strings.Split(text, "\n")
	for i, line := range lines {
		if strings.Contains(line, "'") {
			panic("line contains '!!!")
		}
		cmd := `echo '` + line + `' >`
		if i > 0 {
			cmd += `>`
		}
		cmd += ` ` + filename
		s.cmds = append(s.cmds, cmd)
	}
}

// Exec a command via ssh on the mfi switch.  That's always on the mfi's own Wifi,
// so the IP is usually 192.168.2.20.  Therefore its SSH Fingerprint always changes.
// We don't want SSH to complain about that: http://bit.ly/1Os4sx5
func (s *SshCommands) Exec(host string) (err error) {
	arg := []string{
		"-p", "ubnt",
		"ssh",
		"-o", "UserKnownHostsFile=/dev/null",
		"-o", "StrictHostKeyChecking=no",
		"ubnt@" + host,
		strings.Join(s.cmds, "\n"),
	}
	cmd := exec.Command("sshpass", arg...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
