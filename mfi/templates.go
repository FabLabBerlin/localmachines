package main

const SYSTEM_CFG = `
aaa.status=disabled
wpasupplicant.status=enabled
wireless.status=enabled
wireless.2.wmm=enabled
wireless.2.wds=disabled
wireless.2.status=enabled
wireless.2.ssid=mFi 0wned 
wireless.2.security=none
wireless.2.scan_list.status=disabled
wireless.2.mode=master
wireless.2.macclone=disabled
wireless.2.mac_acl.status=disabled
wireless.2.mac_acl.policy=deny
wireless.2.l2_isolation=disabled
wireless.2.is_guest=false
wireless.2.hide_ssid=disabled
wireless.2.devname=ath1
wireless.2.debug=0x80e81440
wireless.2.authmode=1
wireless.2.ap=
wireless.2.addmtikie=disabled
wireless.1.wmm=enabled
wireless.1.wds=disabled
wireless.1.status=enabled
wireless.1.ssid={{.WifiSSID}}
wireless.1.security=none
wireless.1.scan_list.status=disabled
wireless.1.mode=managed
wireless.1.macclone=disabled
wireless.1.mac_acl.status=disabled
wireless.1.mac_acl.policy=deny
wireless.1.l2_isolation=disabled
wireless.1.is_guest=false
wireless.1.hide_ssid=disabled
wireless.1.devname=ath0
wireless.1.debug=0x80e81440
wireless.1.authmode=1
wireless.1.ap=
wireless.1.addmtikie=disabled
vlan.status=disabled
users.status=enabled
users.1.status=enabled
users.1.password={{.SshPasswordHash}}
users.1.name=ubnt
upnp.status=enabled
upnp.1.status=enabled
upnp.1.devname=ath0
tshaper.status=disabled
telnetd.status=disabled
system.cfg.version=65537
syslog.status=enabled
snmp.status=disabled
route.status=enabled
route.1.status=disabled
resolv.status=enabled
resolv.nameserver.2.status=disabled
resolv.nameserver.1.status=enabled
resolv.nameserver.1.ip=0.0.0.0
resolv.host.1.name=mFi
radio.status=enabled
radio.countrycode=840
radio.1.virtual.1.status=enabled
radio.1.virtual.1.mode=master
radio.1.virtual.1.devname=ath1
radio.1.txpower=auto
radio.1.status=enabled
radio.1.rate.mcs=auto
radio.1.rate.auto=enabled
radio.1.puren=0
radio.1.mode=managed
radio.1.mcastrate=auto
radio.1.ieee_mode=
radio.1.forbiasauto=0
radio.1.devname=ath0
radio.1.cwm.mode=0
radio.1.cwm.enable=0
radio.1.countrycode=840
radio.1.clksel=1
radio.1.chanshift=
radio.1.channel=1
radio.1.ampdu.status=enabled
radio.1.acktimeout=23
radio.1.ackdistance=300
radio.1.ack.auto=disabled
qos.status=disabled
pwdog.status=disabled
ppp.status=disabled
ntpclient.status=disabled
ntpclient.1.status=disabled
ntpclient.1.server=pool.ntp.org
netmode=bridge
netconf.status=enabled
netconf.2.up=enabled
netconf.2.status=enabled
netconf.2.promisc=enabled
netconf.2.netmask=255.255.255.0
netconf.2.ip=192.168.2.20
netconf.2.devname=ath1
netconf.2.autoip.status=disabled
netconf.1.up=enabled
netconf.1.status=enabled
netconf.1.promisc=enabled
netconf.1.netmask=255.255.255.0
netconf.1.ip=192.168.1.20
netconf.1.devname=ath0
netconf.1.autoip.status=disabled
iptables.status=enabled
iptables.1.status=enabled
iptables.1.cmd=-t nat -A PREROUTING --in-interface ath1 -p tcp --dport 80 -j DNAT --to-destination 192.168.2.20
igmpproxy.status=disabled
httpd.status=enabled
ebtables.status=enabled
ebtables.3.status=enabled
ebtables.3.cmd=-t broute -A BROUTING --protocol 0x888e --in-interface ath0 -j DROP
ebtables.2.status=disabled
ebtables.2.cmd=-t nat -A POSTROUTING --out-interface ath0 -j arpnat --arpnat-target ACCEPT
ebtables.1.status=disabled
ebtables.1.cmd=-t nat -A PREROUTING --in-interface ath0 -j arpnat --arpnat-target ACCEPT
dnsmasq.status=enabled
dnsmasq.1.status=enabled
dnsmasq.1.devname=ath1
dhcpd.status=enabled
dhcpd.1.status=enabled
dhcpd.1.start=192.168.2.100
dhcpd.1.redirect=192.168.2.20
dhcpd.1.netmask=255.255.255.0
dhcpd.1.lease_time=3600
dhcpd.1.end=192.168.2.200
dhcpd.1.dnsproxy=enabled
dhcpd.1.dns.2.status=disabled
dhcpd.1.dns.2.server=
dhcpd.1.dns.1.status=disabled
dhcpd.1.dns.1.server=
dhcpd.1.devname=ath1
dhcpc.status=enabled
dhcpc.1.status=enabled
dhcpc.1.fallback=0.0.0.0
dhcpc.1.devname=ath0
bridge.status=disabled
aaa.1.status=disabled
wpasupplicant.device.1.status=enabled
wpasupplicant.device.1.devname=ath0
wpasupplicant.device.1.driver=madwifi
wpasupplicant.profile.1.network.1.proto.1.name=RSN
wpasupplicant.profile.1.network.1.pairwise.1.name=CCMP
wpasupplicant.profile.1.network.1.ssid={{.WifiSSID}}
wpasupplicant.device.1.profile=WPA-PSK
wpasupplicant.profile.1.status=enabled
wpasupplicant.profile.1.name=WPA-PSK
wpasupplicant.profile.1.network.1.key_mgmt.1.name=WPA-PSK
wpasupplicant.profile.1.network.1.psk={{.WifiPassword}}
wpasupplicant.profile.1.network.1.eap.1.status=disabled
`

const WLAN_OVERWRITE = `
essid={{.WifiSSID}}
security=wpapsk
wpakey={{.WifiPassword}}
`
