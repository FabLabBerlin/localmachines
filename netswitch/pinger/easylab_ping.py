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
 
Laser_Cutter_Epilog_Zing=Drucker(231,0,0,"Laser Cutter Epilog Zing")
CNC_Router=Drucker(180,0,0,"CNC Router")
threeD_Printer_Vincent_1=Drucker(141,0,0,"3D Printer Vincent 1")
threeD_Printers_2345678=Drucker(25,0,0,"3D Printers 2,3,4,5,6,7,8")
Trotec_Laser=Drucker(109,0,0,"Trotec Laser")
Electronic_Desk=Drucker(185,0,0,"Electronic Desk")
Circle_Saw=Drucker(186,0,0,"Circle Saw")
Heat_Press=Drucker(242,0,0,"Heat Press")
PCB_CNC_Machine_LPKF=Drucker(132,0,0,"PCB CNC Machine (LPKF)")

l=[Laser_Cutter_Epilog_Zing,CNC_Router,threeD_Printer_Vincent_1,threeD_Printers_2345678,Trotec_Laser,Electronic_Desk,Circle_Saw]
#Heat Press is down und pcb cnc machine..
#no internet.. in m network..
#Laser_Cutter_Epilog_Zing,CNC_Router,threeD_Printer_Vincent_1,threeD_Printers_2345678,Trotec_Laser,Electronic_Desk,
#erst ethernet, dann wlan

def pingit(x):
    global counter
    address = "192.168.1." + str(x.ip_ending)
    if os.system("ping -c 1 " + address) == 0:
    #    if x.emailsent==1:
        #    server = smtplib.SMTP('smtp.gmail.com', 587)
       #     server.ehlo()
      #      server.starttls()
     #       server.ehlo()
    #        server.login("machinepinger", "raspberrypi")
   #         msg = "\nHi! The machine * "+str(x.name) + " * with ip address " + str(address) + " is WORKING AGAIN! :-)"
          #  server.sendmail("machinepinger@gmail.com", "julius@fablab.berlin", msg)
           # server.sendmail("machinepinger@gmail.com", "krisjanis.rijnieks@gmail.com", msg)
            #server.sendmail("machinepinger@gmail.com", "philip@fablab.berlin", msg)
        x.emailsent=0
        x.erreichbar=1
        counter=0
    else:           
        counter+=1
        if counter<=0:  #counter digit in secs (+waittime=60secs)
            pingit(x)
        else:
            if x.emailsent==0:
            #    server = smtplib.SMTP('smtp.gmail.com', 587)
           #     server.ehlo()
          #      server.starttls()
        #        server.ehlo()
       #         server.login("machinepinger", "raspberrypi")
         #       msg = "\nHi! The machine * "+str(x.name) + " * with ip address " + str(address) + " is DOWN. (There is no network connection available.)"
               # server.sendmail("machinepinger@gmail.com", "julius@fablab.berlin", msg)
                #server.sendmail("machinepinger@gmail.com", "krisjanis.rijnieks@gmail.com", msg)
                #server.sendmail("machinepinger@gmail.com", "philip@fablab.berlin", msg)
                x.emailsent=1
                x.erreichbar=0
            counter=0
counter=0


def jo():
    r=[]
    z=0
    for ping in l:
        pingit(ping)
        for i in l:
            if i.erreichbar ==0:
                r.append("The machine * "+str(i.name) + " * with ip address " + "192.168.1." + str(i.ip_ending) + " is NOT working :-(")
            else:
                r.append("The machine * "+str(i.name) + " * with ip address " + "192.168.1." + str(i.ip_ending) + " is working :-)")
        
        return r    
    #print(r)
     #   for i in l:
      #      print(i.emailsent)
       # print(r)
        #r=[]
        #time.sleep(2)       #Time period for checking ping in secs
#        if l[z].erreichbar ==0:
 #           r.append("The machine * "+str(l[z].name) + " * with ip address " + "192.168.1." + str(l[z].ip_ending) + " is NOT working :-(")
  #      else:
   #         r.append("The machine * "+str(l[z].name) + " * with ip address " + "192.168.1." + str(l[z].ip_ending) + " is working :-)")
    #    z+=1
    #return r    

#x  print(str(er[0]))

app = Flask(__name__, static_url_path='')

@app.route("/get_list")
def get_list():
    er=jo()
    return json.dumps({'status': er});#str(er[0])});

@app.route('/public/<string:filename>')
def send_js(filename):
    app.logger.error('filename: ' + filename)
    print("filename: " + filename)
    return send_from_directory("public", filename)


if __name__ == "__main__":
    app.run(debug=True)

#r=[]
#while True:
 #   for ping in l:
  #      pingit(ping)
   #     for i in l:
    #        if i.erreichbar ==0:
     #           r.append("The machine * "+str(i.name) + " * with ip address " + "192.168.1." + str(i.ip_ending) + " is NOT working :-(")
      #      else:
       #         r.append("The machine * "+str(i.name) + " * with ip address " + "192.168.1." + str(i.ip_ending) + " is working :-)")
    # #   for i in l:
     # #      print(i.emailsent)
     #   print(r)
      #  r=[]
       # time.sleep(2)       #Time period for checking ping in secs






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

