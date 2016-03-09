/*
Commandline tool for automatic configuration of mfi switches.

TODO: use github.com/ziutek/telnet
      (golang.org/x/crypto/ssh doesn't support cbc ciphers because they're
      unsafe)
      use parser for hwaddr to check for equality

*/
package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"text/template"
	"time"
)

const (
	SSH_INITIAL_PASSWORD = "ubnt"
	// http://blog.vucica.net/2014/08/mfi-mpower-basic-use-without-cloud-and-controller.html
	SSH_PASSWORD_HASH = "KQiBBQ7dx8sx2" // = hash("my_password")
	TEMPFILE_PREFIX   = "mfi_deploy"
)

type Config struct {
	EthernetConfig        bool
	WlanOverwriteFilename string
	SystemCfgFilename     string
	WifiSSID              string
	WifiPassword          string
	SshPasswordHash       string
	HwAddr                string
}

func (c *Config) usage() {
	fmt.Printf("Please hold the power switch's reset button until it\n")
	fmt.Printf("blinks blue-orange.  ")
	if c.EthernetConfig {
		fmt.Printf("Be sure to be connected with\n")
		fmt.Printf("the device with an Ethernet cable.\n")
		fmt.Printf("Be on 192.168.1.x (The device is now on 192.168.1.20)\n")
	} else {
		fmt.Printf("Afterwards wait until a Wifi\n")
		fmt.Printf("named 'mfi...' comes up and connect to it.\n")
		fmt.Printf("(The device is now on 192.168.2.20)\n")
	}
}

func (c *Config) Host() string {
	if c.EthernetConfig {
		return "192.168.1.20"
	} else {
		return "192.168.2.20"
	}
}

func (c *Config) Run() (err error) {
	if c.WifiPassword, err = c.getWifiPw(); err == nil {
		log.Printf("Wifi password automatically obtained")
	} else {
		log.Printf("Cannot obtain Wifi password automatically")
		fmt.Printf("What is the Wifi password?\n")
		if _, err = fmt.Scanln(&c.WifiPassword); err != nil {
			return fmt.Errorf("scan ln: %v", err)
		}
	}
	if err = c.generate(); err != nil {
		return fmt.Errorf("generate: %v", err)
	}
	if err = c.scp(); err != nil {
		return fmt.Errorf("scp: %v", err)
	}
	if c.HwAddr, err = c.getHwAddr(); err != nil {
		return fmt.Errorf("get hw addr: %v", err)
	}
	if err = c.finalize(); err != nil {
		return fmt.Errorf("reboot: %v", err)
	}
	os.Remove(c.SystemCfgFilename)
	os.Remove(c.WlanOverwriteFilename)
	return
}

func (c *Config) getWifiPw() (hwAddr string, err error) {
	cmd := exec.Command("sshpass", "-p", "ubnt", "ssh", "-o", "UserKnownHostsFile=/dev/null", "-o", "StrictHostKeyChecking=no", "ubnt@"+c.Host(), "cat /etc/wpasupplicant_WPA-PSK.conf | grep psk")
	buf := bytes.NewBufferString("")
	cmd.Stdout = buf
	cmd.Stderr = os.Stderr
	if err = cmd.Run(); err != nil {
		return "", fmt.Errorf("ifconfig: %v", err)
	}
	s := strings.TrimSpace(buf.String())
	if len(s) < 8 {
		return "", fmt.Errorf("unexpected ifconfig output: '%v'", s)
	}
	s = s[len(`psk="`):]
	s = s[:len(s)-len(`"`)]
	return s, nil
}

func (c *Config) generate() (err error) {
	c.SystemCfgFilename, err = c.generateFromTemplate(SYSTEM_CFG)
	if err != nil {
		return fmt.Errorf("generate from system.cfg: %v", err)
	}
	c.WlanOverwriteFilename, err = c.generateFromTemplate(WLAN_OVERWRITE)
	if err != nil {
		return fmt.Errorf("generate from wlan_overwrite: %v", err)
	}
	return
}

func (c *Config) generateFromTemplate(templateText string) (resultFilename string, err error) {
	t := template.New("cfg")
	t, err = t.Parse(templateText)
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
	cmd := exec.Command("sshpass", "-p", "ubnt", "scp", "-o", "UserKnownHostsFile=/dev/null", "-o", "StrictHostKeyChecking=no", c.SystemCfgFilename, "ubnt@"+c.Host()+":/tmp/system.cfg")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err = cmd.Run(); err != nil {
		return fmt.Errorf("scp system.cfg: %v", err)
	}
	cmd = exec.Command("sshpass", "-p", "ubnt", "scp", "-o", "UserKnownHostsFile=/dev/null", "-o", "StrictHostKeyChecking=no", c.WlanOverwriteFilename, "ubnt@"+c.Host()+":/tmp/wlan_overwrite")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err = cmd.Run(); err != nil {
		return fmt.Errorf("scp system.cfg: %v", err)
	}
	return
}

func (c *Config) getHwAddr() (hwAddr string, err error) {
	cmd := exec.Command("sshpass", "-p", "ubnt", "ssh", "-o", "UserKnownHostsFile=/dev/null", "-o", "StrictHostKeyChecking=no", "ubnt@"+c.Host(), "ifconfig | grep wifi0")
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
	return tmp[len(tmp)-1], nil
}

func (c *Config) finalize() (err error) {
	sshCmds := NewSshCommands()
	sshCmds.Add("cfgmtd -w")
	sshCmds.Add("sleep 3")
	rcPoststart := `#!/bin/sh

/usr/bin/echo 0 > /proc/power/output1
/usr/bin/echo 0 > /proc/power/output2
/usr/bin/echo 0 > /proc/power/output3
/usr/bin/echo 0 > /proc/power/output4
/usr/bin/echo 0 > /proc/power/output5
/usr/bin/echo 0 > /proc/power/output6
	`
	sshCmds.AddFile("/etc/persistent/rc.poststart", rcPoststart)
	sshCmds.Add("chmod a+x /etc/persistent/rc.poststart")
	sshCmds.Add("cfgmtd -w -p /etc")
	sshCmds.Add("sync")
	sshCmds.Add("reboot")
	if err := sshCmds.Exec(c.Host()); err != nil {
		return fmt.Errorf("ssh cmds exec: %v", err)
	}
	return
}

func (c *Config) FindDeviceOn(netmask string) (resultIp net.IP, err error) {
	ip, ipnet, err := net.ParseCIDR(netmask)
	if err != nil {
		return net.IP{}, fmt.Errorf("parse cidr: %v", err)
	}
	ch := make(chan net.IP, 1)
	for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); incIp(ip) {
		scanIp := make(net.IP, len(ip))
		copy(scanIp, ip)
		go func(ip net.IP) {
			if c.deviceOn(ip) {
				log.Printf("found it on %v!!!!!!", ip.String())
				select {
				case ch <- ip:
				default:
					log.Fatalf("fatal error: %v", ip)
				}
			}
		}(scanIp)
	}
	select {
	case resultIp = <-ch:
		break
	case <-time.After(30 * time.Second):
		err = fmt.Errorf("could not find device")
	}
	return
}

func incIp(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

type MfiStatus struct {
	Wlan MfiWlan `json:"wlan"`
}

type MfiWlan struct {
	HwAddr string `json:"hwaddr"`
}

func (c *Config) deviceOn(ip net.IP) bool {
	resp, err := http.Get("http://" + ip.String() + "/status.cgi")
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	dec := json.NewDecoder(resp.Body)
	status := MfiStatus{}
	if err := dec.Decode(&status); err != nil {
		return false
	}
	hwAddr := strings.Replace(c.HwAddr, "-", ":", -1)
	return len(hwAddr) > 10 && strings.HasPrefix(hwAddr, status.Wlan.HwAddr)
}

func main() {
	ethernetConfig := flag.Bool("ethConfig", false, "Connect Switch to computer via Ethernet Cable (be on 192.168.1.x)")
	netmask := flag.String("network", "172.26.0.128/24", "Wifi network's netmask on which we auto discover the switch")
	ssid := flag.String("ssid", "FabLab-M", "Wifi SSID")
	flag.Parse()

	c := &Config{
		EthernetConfig: *ethernetConfig,
		WifiSSID:       *ssid,
	}
	if err := c.Run(); err == nil {
		fmt.Printf("Your switch is properly configured and its hardware")
		fmt.Printf(" address is: '%v'\n", c.HwAddr)
		fmt.Printf("Wait until the LED starts blinking blue, connect to the")
		fmt.Printf(" normal Wifi and then press enter...\n")
		var tmp string
		fmt.Scanln(&tmp)
		if ip, err := c.FindDeviceOn(*netmask); err == nil {
			fmt.Printf("Successfully discovered device: %v\n", ip.String())
		} else {
			fmt.Printf("Unable to find device on %v\n.", *netmask)
		}
	} else {
		log.Printf("config: %v", err)
		fmt.Printf("Make sure the switch is properly connected:\n\n")
		c.usage()
	}
}
