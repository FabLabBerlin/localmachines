# Gateway Server

This is a lightweight server that can run in a lab.

It communicates mostly via XMPP, so it can run behind a firewall.  So it's possible to connect switches which are designed for local networks to localmachines.com.

Besides communication it also persists the relay states.  When plugging in a new switch, the Gateway Server makes sure that it gets turned off.  Or when a switch fails to receive a command, the Gateway server can resend the command later.
