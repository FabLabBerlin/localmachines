# Virtual NetSwitch API
This is a modified version of the Python NetSwitch API that can be found [here](https://github.com/kr15h/sentry-py-api).

Use `./restapi.py` to launch the server and access it by `http://localhost:8080/a1/on` and `http://localhost:8080/a1/off`. It will randomly decide whether to fail or not. Use `./restapi.py 8888` (exectuable with a port argument) to launch the API on a custom port.

Before launching, make sure that the `web.py` module is installed. Install it with python-setuptools: `sudo easy_install web.py`.

It might be that you need to install the python serial module as well. To do that use `sudo easy_install pyserial`.

If your system does not recognise the command `easy_install`, use `sudo apt-get install python-setuptools` to set it up.