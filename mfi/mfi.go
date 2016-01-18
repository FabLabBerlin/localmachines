// CLI tool for automatic configuration of mfi swithces
//
// TODO: use golang.org/x/crypto/ssh

package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"text/template"
)

const (
	SSH_INITIAL_PASSWORD = "ubnt"
	SSH_PASSWORD_HASH    = "KQiBBQ7dx8sx2" // = hash("ubnt")
	TEMPFILE_PREFIX      = "mfi_deploy"
)

func usage() {
	fmt.Printf("Please hold the power switch's reset button until it\n")
	fmt.Printf("blinks blue-orange.  Afterwards wait until a Wifi\n")
	fmt.Printf("named 'mfi...' comes up and connect to it.\n")
	fmt.Printf("(The device is now on 192.168.2.20)\n")
}

type Config struct {
	WlanOverwriteFilename string
	SystemCfgFilename     string
	WifiPassword          string
	SshPasswordHash       string
	HwAddr                string
}

func (c *Config) Run() (err error) {
	if err = c.generate(); err != nil {
		return fmt.Errorf("generate: %v", err)
	}
	fmt.Printf("\nWhen asked for an SSH password, just enter '%v'\n\nPress [enter] to continue...\n", SSH_INITIAL_PASSWORD)
	var string tmp
	fmt.Scanln(&tmp)
	if err = c.scp(); err != nil {
		return fmt.Errorf("scp: %v", err)
	}
	if c.HwAddr, err = c.getHwAddr(); err != nil {
		return fmt.Errorf("get hw addr: %v", err)
	}
	if err = c.reboot(); err != nil {
		return fmt.Errorf("reboot: %v", err)
	}
	os.Remove(c.SystemCfgFilename)
	os.Remove(c.WlanOverwriteFilename)
	return
}

func (c *Config) generate() (err error) {
	fmt.Printf("What is the Wifi password?\n")
	if _, err = fmt.Scanln(&c.WifiPassword); err != nil {
		return fmt.Errorf("scan ln: %v", err)
	}
	c.SystemCfgFilename, err = c.generateFromFile("templates/system.cfg")
	if err != nil {
		return fmt.Errorf("generate from system.cfg: %v", err)
	}
	c.WlanOverwriteFilename, err = c.generateFromFile("templates/wlan_overwrite")
	if err != nil {
		return fmt.Errorf("generate from wlan_overwrite: %v", err)
	}
	return
}

func (c *Config) generateFromFile(templateFilename string) (resultFilename string, err error) {
	t, err := template.ParseFiles(templateFilename)
	if err != nil {
		return "", fmt.Errorf("parse template: %v", err)
	}
	f, err := ioutil.TempFile(os.TempDir(), TEMPFILE_PREFIX)
	if err != nil {
		return "", fmt.Errorf("temp file: %v", err)
	}
	resultFilename = f.Name()
	defer f.Close()
	if err = t.Execute(f, *c); err != nil {
		return "", fmt.Errorf("execute template: %v", err)
	}
	return
}

func (c *Config) scp() (err error) {
	cmd := exec.Command("scp", c.SystemCfgFilename, "ubnt@192.168.2.20:/tmp/system.cfg")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err = cmd.Run(); err != nil {
		return fmt.Errorf("scp system.cfg: %v", err)
	}
	cmd = exec.Command("scp", c.WlanOverwriteFilename, "ubnt@192.168.2.20:/tmp/wlan_overwrite")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err = cmd.Run(); err != nil {
		return fmt.Errorf("scp system.cfg: %v", err)
	}
	return
}

func (c *Config) getHwAddr() (hwAddr string, err error) {
	cmd := exec.Command("ssh", "ubnt@192.168.2.20", "ifconfig | grep wifi0")
	buf := bytes.NewBufferString("")
	cmd.Stdout = buf
	cmd.Stderr = os.Stderr
	if err = cmd.Run(); err != nil {
		return "", fmt.Errorf("ifconfig: %v", err)
	}
	s := strings.TrimSpace(buf.String())
	tmp := strings.Split(s, " ")
	if len(tmp) == 0 {
		return "", fmt.Errorf("unexpected ifconfig output: '%v'", s)
	}
	log.Printf("tmp = '%v'", tmp)
	log.Printf("s = '%v'", s)
	return tmp[len(tmp)-1], nil
}

func (c *Config) reboot() (err error) {
	cmd := exec.Command("ssh", "ubnt@192.168.2.20", "cfgmtd -w")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err = cmd.Run(); err != nil {
		return fmt.Errorf("cfgmtd: %v", err)
	}
	cmd = exec.Command("ssh", "ubnt@192.168.2.20", "reboot")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err = cmd.Run(); err != nil {
		return fmt.Errorf("reboot: %v", err)
	}
	return
}

func main() {
	c := &Config{}
	if err := c.Run(); err == nil {
		fmt.Printf("Your switch is properly configured and it's hardware")
		fmt.Printf(" address is: '%v'\n", c.HwAddr)
		fmt.Printf("Wait until the LED starts blinking blue...\n")
	} else {
		log.Printf("config: %v", err)
		fmt.Printf("Make sure the switch is properly connected:\n\n")
		usage()
	}
}
