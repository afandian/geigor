Geigor
======

*A Geiger counter written in Go. (Hence the 'go').*

Copyright Joe Wass 2013

MIT License, see LICENSE file.

joe@afandian.com

Put a tracking script in your HTML. Listen to the geigier counter-style monitor to hear activity on your site. Install Geigor on your webserver (it's written in Go), start it with Supervisor and off you go. The Geigor tracking should work on any browser, the monitor currently only works on latest versions of FF, Firefox and Safari, see http://caniuse.com/#feat=audio-api .

To build
--------

    go get "code.google.com/p/go.net/websocket"
    go build

To run
------

You need to supply the 'served from' address, port number, and diagnostic options. The served from address will be the base URL for the host name by default, but if you are using a proxy server, this might be different.

    ./geigor <served_from> <port number> <option>*


e.g. default:

    ./geigor http://mydomain.com:9319/ 9319 index demo


or, if you're proxying


    ./geigor http://mydomain.com/geigor 80 index demo


You can use supervisor to keep this running on your server. There is a sample config file at etc/supervisor/geigor.conf

Options are:

 - index : show the helpful index page
 - demo : show the demo tracker page

Usage
-----

Include the script (see demo page) in your tracking pages. Visit the monitor page (./monitor) to hear ticks.

All paths (index, monitor, demo, js file) are relative so this should work nicely through a proxy or on a port of your choice. The tracker sents POST requests, so the tracker must be running on the same domain as the pages you're serving.