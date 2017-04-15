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
