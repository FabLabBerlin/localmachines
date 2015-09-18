#Pinger
Pinger is used to ping all machines with an ip address in a network. When the program is run the status of the machines can be seen on a website with the status being updated continuously every 20 seconds. Optionally email notifications can be sent, when a machine is down and is running again.

## Table of contents
- [Set-up](#set-up)
	- [Installing Flask](#installing-flask)
	- [Machine Names and Addresses](#machine-names-and-addresses)
	- [Email Notifications](#email-notifications)
- [Running Pinger](#running-pinger)
	
## Set-up
### Installing Flask
To install Flask execute the following command in your terminal:
```
-pip install Flask
```
Flask is needed to send data from the back-end to the front-end and to set a url path for the data.

### Machine names and Addresses
You need to change the ip addresses and names of your machines according to the comments in `easylab_ping.py`.

### Email Notifications
If you want to receive email notifications (if a machine is down and back up), you have to enable it according to the comments in `easylab_ping.py`. Email notifications are set off by default.

## Running Pinger
To run Pinger you need to open your terminal and enter:
```
easylab_ping.py
```
To see the status of the machines go to `http://localhost:5000/public/ping.html` on your browser:
