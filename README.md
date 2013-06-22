Geigor
======

*A Geiger counter written in Go. (Hence the 'go').*

Copyright Joe Wass 2013
joe@afandian.com

Put a tracking script in your HTML. Listen to the geigier counter-style monitor to hear activity on your site.

To build
--------

go fetch "code.google.com/p/go.net/websocket"
go build

To run
------

./geigor <port number> <option>*

Options are:

 - index : show the helpful index page
 - demo : show the demo tracker page

Usage
-----

Include the script (see demo page) in your tracking pages. Visit the monitor page (./monitor) to hear ticks.

All paths (index, monitor, demo, js file) are relative so this should work nicely through a proxy or on a port of your choice. The tracker sents POST requests, so the tracker must be running on the same domain as the pages you're serving.