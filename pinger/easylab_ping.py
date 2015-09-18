#!/usr/bin/env python3
import os, time
import smtplib
from flask import Flask, json, send_from_directory

#time.sleep(40)

class Drucker(object):
    def __init__(self,ip_ending,emailsent,erreichbar,name): 
        self.ip_ending=ip_ending
        self.emailsent=emailsent
        self.erreichbar=erreichbar
        self.name=name
 
Laser_Cutter_Epilog_Zing=Drucker(231,0,0,"Laser Cutter Epilog Zing")    #change first nr in brackets (ending of ip address of machine)
CNC_Router=Drucker(180,0,0,"CNC Router")                                #also change names of machines accordingly
threeD_Printer_Vincent_1=Drucker(141,0,0,"3D Printer Vincent 1")        #use these lines to remove/add machines
threeD_Printers_2345678=Drucker(25,0,0,"3D Printers 2,3,4,5,6,7,8")
Trotec_Laser=Drucker(109,0,0,"Trotec Laser")
Electronic_Desk=Drucker(185,0,0,"Electronic Desk")
Circle_Saw=Drucker(186,0,0,"Circle Saw")
Heat_Press=Drucker(242,0,0,"Heat Press")
PCB_CNC_Machine_LPKF=Drucker(132,0,0,"PCB CNC Machine (LPKF)")

#add/remove printers to this list l:
l=[Laser_Cutter_Epilog_Zing,CNC_Router,threeD_Printer_Vincent_1,threeD_Printers_2345678,Trotec_Laser,Electronic_Desk,Circle_Saw,Heat_Press,PCB_CNC_Machine_LPKF]

def pingit(x):
    global counter
    address = "192.168.1." + str(x.ip_ending)   #change ip adress here
    if os.system("ping -c 1 -W 100 " + address) == 0:      #-W is waittime in ms
    #for enabling email notifications remove all hashtags from here to x.emailsent=0
    #    if x.emailsent==1:
        #    server = smtplib.SMTP('smtp.gmail.com', 587)
       #     server.ehlo()
      #      server.starttls()
     #       server.ehlo()
    #        server.login("machinepinger", "raspberrypi")
   #         msg = "\nHi! The machine * "+str(x.name) + " * with ip address " + str(address) + " is WORKING AGAIN! :-)"
          #  server.sendmail("machinepinger@gmail.com", "julius@fablab.berlin", msg)        #fill in your email address(es) for notification here
           # server.sendmail("machinepinger@gmail.com", "krisjanis.rijnieks@gmail.com", msg)
            #server.sendmail("machinepinger@gmail.com", "philip@fablab.berlin", msg)
        x.emailsent=0
        x.erreichbar=1
        counter=0
    else:           
        counter+=1
        if counter<=5:  #counter digit in secs (+waittime=60secs)
            pingit(x)
        else:
            if x.emailsent==0:
            #for enabling email notifications remove all hashtags from here to x.emailsent=1
            #    server = smtplib.SMTP('smtp.gmail.com', 587)
           #     server.ehlo()
          #      server.starttls()
        #        server.ehlo()
       #         server.login("machinepinger", "raspberrypi")
         #       msg = "\nHi! The machine * "+str(x.name) + " * with ip address " + str(address) + " is DOWN. (There is no network connection available.)"
               # server.sendmail("machinepinger@gmail.com", "julius@fablab.berlin", msg)          #fill in your email address(es) for notification here
                #server.sendmail("machinepinger@gmail.com", "krisjanis.rijnieks@gmail.com", msg)
                #server.sendmail("machinepinger@gmail.com", "philip@fablab.berlin", msg)
                x.emailsent=1
                x.erreichbar=0
            counter=0
counter=0


def jo():
    r=[]
    for ping in l:
        pingit(ping)
    for i in l:
        adr="192.168.1." + str(i.ip_ending)
        dic={"name":str(i.name),"ip":str(adr),"status":"0"}
        if i.erreichbar ==0:
            dic["status"] = "DOWN"
            r.append(dic)
        else:
            dic["status"]="UP"
            r.append(dic)
    return r

app = Flask(__name__, static_url_path='')

@app.route("/get_list")
def get_list():
    er=jo()
    return json.dumps({'status': er});

@app.route('/public/<path:filename>')
def send_js(filename):
    app.logger.error('filename: ' + filename)
    print("filename: " + filename)
    return send_from_directory("public", filename)


if __name__ == "__main__":
    app.run(host="0.0.0.0")
    
    #app.run(debug=True)


#changed print it, pingit 10 to 1,   6 times #server.sandmail, timeslep
#192.168.1.231  Laser Cutter Epilog Zing
#192.168.1.180  CNC Router
#192.168.1.141  3D 1 Vincent vega   Rep2
#192.168.1.25:8080/a3  3D 2 Jules  Rep5
#192.168.1.25:8080/a1  3D 7 Fabienne Rep2
#192.168.1.25:8080/b2  3D 6 Honey Bunny i3
#192.168.1.226  Vinyl Cutter
#192.168.1.25:8080/a2  3D 4 Mia Rep5
#192.168.1.25:8080/b1  3D 5 Pumplin i3
#192.168.1.25:8080/b3   3D 3 Mr. Wallace RepZ18
#..25: 8 BUTCH 3D
#192.168.1.109  Trotec Laser
#192.168.1.185  Electronic Desk
#192.168.1.186  Circle Saw
#192.168.1.242  Heat Press
#192.168.1.132  PCB CNC Machine (LPKF)

