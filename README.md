# HttpRoller
An http server that dynamically serves static content depending upon the elapsed 
time within a sliding window.

The HttpRoller is a web server that serves static content based upon the 
time at which a request is made.  When started this server will examine the -path 
option to determine the location of one or more directories that are named according
to the elapsed seconds within a rolling window.  When the location of the content
to be served is determined based upon the time in seconds within the rolling
window then the request will be served from that location/

When starting the web server you would use a command such as:

<pre>
/home/pi/bin/HttpRoller -listen=127.0.0.1:12345 -path=~/pi-gateway/simulator/scenarios/default -window=60s
</pre>

If the directory scenarios/default contained two directories with a file named json:

scenarios/default/0/module/status/json
scenarios/default/30/module/status/json

Then the result would be that for the first 30 seconds of every minute after 
the server was started the first json would be served, and the second json served after
the 30 second period, and the rolling back to the after after that.

The finest granularity is 1 second for the selection of the directories.  Directories 
with names that represent numbers larger than the window option allows will be
ignored.

## Headless Installation

The HttpRoller is generally used as an interactive test tool however it can be 
installed as a systemd service.  Doing this requires that the HttpRoller.service
file is first modified to point at the test scenario that is needed.  The default
setting points to a test case within the pi-gateway suite.

The HttpRoller can be installed using systemd.  The HttpRoller.service file should be
modified to ensure that the appropriate listening socket is set and that the path
argument points to the appropriate location for the rolling test.  The service file
is then copied into the /lib/systemd/system/ directory.

These instructions assume that the HttpRoller has been git cloned into your
/home/pi/HttpRoller directory and that the binaries and other files are
provided to the systemd daemon from this location.  It also assumes for
the default test case that the pi-gateway has been git cloned also into the
/home/pi/pi-gateway directory and that the test scenario files are there.

If you have used a binary manually copied to the system or other arrangement 
you will need to modify the HttpRoller.service file before using the following
instructions.

<pre>
sudo cp HttpRoller.service /lib/systemd/system/HttpRoller.service
sudo chmod 644 /lib/systemd/system/HttpRoller.service

sudo systemctl daemon-reload
sudo systemctl enable HttpRoller.service
</pre>

After headless installation the Pi should be rebooted. Logging output
from the pi-gateway unit can be seen using the following command:

<pre>
sudo journalctl -u HttpRoller
</pre>

To follow the output use the '-f' option.

