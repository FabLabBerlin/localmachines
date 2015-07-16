#!/usr/bin/env python
import time
from datetime import datetime
import web
import random
import sys

# Setup web server
urls = (
    #  Translates to /a1/on or /a1/off
    '/(.+)/(.+)', 'switch'  
)

app = web.application(urls, globals())

class switch:
    def GET(self, switchId, switchState):
        if round(random.random()):
            output = 'succ'
        else:
            raise web.internalerror('fail')
        return output 

if __name__ == "__main__":
    app.run()