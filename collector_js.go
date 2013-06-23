package main

var COLLECTOR_JS string = `

var endpoint = "{{.ServedFrom}}/collector-endpoint";
var period = 10;

var request = null;

if (window.XMLHttpRequest) {
    request = new XMLHttpRequest();
} else {
    try {
        request = new ActiveXObject("Microsoft.XMLHTTP");
    } catch (e) {
        // Ignore. Nothing can be done.
    }
}

var tick = function() {
    request.open("POST", endpoint, true);
    request.send();

    setTimeout(tick, period * 1000);
};

if (request != null) {
    tick();
}

`
