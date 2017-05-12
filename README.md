# HttpRoller
An http server that dynamically serves static content depending upon the elapsed 
time within a sliding window.

The HttpRoller is a web server that serves static content based upon the 
time at which a request is made.  When started this server will examine the -path 
option to determine the location of one or more directories that are named according
to the elapsed seconds within a rolling window.  When the location of the content
to be served is determined based upon the time in seconds within the rolling
window then the request will be served from that location.

The time slot at which the test should be restarted and rewound is defined by
having a slot directory that contains a file "finish".

When starting the web server you would use a command such as:

<pre>
/home/pi/bin/HttpRoller -listen=127.0.0.1:12345 -path=/home/pi/pi-gateway/simulator/scenarios/portal_builds
</pre>

If the directory scenarios/default contained two directories with a file named json:

scenarios/default/0/module/status/json
scenarios/default/30/module/status/json
scenarios/default/60/finish

Then the result would be that for the first 30 seconds of every minute after 
the server was started the first json would be served, and the second json served after
the 30 second period, and the rolling back to the after after that.

The finest granularity is 1 second for the selection of the directories.  Directories 
with names that represent numbers larger than the window option allows will be
ignored.

## Scenario Push

This feature allows for completely remotely controlled testing to occur.

The HttpRoller is equiped with a special HTTP Endpoint, "/configure".  This
endpoint is enabled only when the -remote option is enabled.  When enabled
an Http GET operation that starts with "/configure" will result in the URI path 
without the "/configure" portion being interpreted as file system path, and 
that in turn is used to populate a new test scenario to be loaded into the server.

As an example using the following linux command would reset the test
server to initate the XM_drain_all test from the pi-gateway test
suite and serve it up to the pi-gateway.

<pre>
wget -O- http://127.0.0.1:12345/configure/home/pi/pi-gateway/simulator/scenarios/XM_drain_all
</pre>

If you were running the HttpRoller in a shell session and ran the above command
at approximately the 9 second mark you could well see something
like the following:

<pre>
bin/HttpRoller -path scenario/default -loglevel=debug --listen=127.0.0.1:12345 -remote
21:36:27.178541 DBG HttpRoller loaded scenario scenario/default
21:36:27.679789 DBG HttpRoller using scenario/default/0
21:36:28.180561 DBG HttpRoller using scenario/default/0
21:36:28.679786 DBG HttpRoller using scenario/default/0
21:36:29.179780 DBG HttpRoller using scenario/default/2
21:36:29.679787 DBG HttpRoller using scenario/default/2
21:36:30.179997 DBG HttpRoller using scenario/default/2
21:36:30.679968 DBG HttpRoller using scenario/default/2
21:36:31.179981 DBG HttpRoller using scenario/default/5
21:36:31.180872 DBG HttpRoller loaded scenario scenario/default
21:36:31.679987 DBG HttpRoller using scenario/default/0
21:36:32.180013 DBG HttpRoller using scenario/default/0
21:36:32.680036 DBG HttpRoller using scenario/default/0
21:36:33.180026 DBG HttpRoller using scenario/default/2
21:36:33.679865 DBG HttpRoller using scenario/default/2
21:36:34.179799 DBG HttpRoller using scenario/default/2
21:36:34.623002 DBG HttpRoller forced load of /home/pi/pi-gateway/simulator/scenarios/XM_drain_all occurring
21:36:34.623955 DBG HttpRoller loaded scenario /home/pi/pi-gateway/simulator/scenarios/XM_drain_all
21:36:34.679859 DBG HttpRoller using /home/pi/pi-gateway/simulator/scenarios/XM_drain_all/0
21:36:35.179881 DBG HttpRoller using /home/pi/pi-gateway/simulator/scenarios/XM_drain_all/0
21:36:35.679719 DBG HttpRoller using /home/pi/pi-gateway/simulator/scenarios/XM_drain_all/0
21:36:36.179736 DBG HttpRoller using /home/pi/pi-gateway/simulator/scenarios/XM_drain_all/0
21:36:36.679737 DBG HttpRoller using /home/pi/pi-gateway/simulator/scenarios/XM_drain_all/0
21:36:37.179753 DBG HttpRoller using /home/pi/pi-gateway/simulator/scenarios/XM_drain_all/5
</pre>

The server will attempt to ignore and return errors for any relative path 
references used.

The new scenario will be applied immediately.

This of course has a major potential for absue should a third party get 
access to the server and so should only be used in private networks, or
under strictly controlled circumstances.

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

