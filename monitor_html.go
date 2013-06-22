package main

var MONITOR_HTML string = `<html>
    <body>
        <p>Geigor tracking. Current count: <span id="count">unknown</span>.</p>
        <script type="text/javascript">
        var context;
        window.addEventListener('load', init, false);
        function init() {
          try {
            window.AudioContext = window.AudioContext||window.webkitAudioContext;
            context = new AudioContext();

            var spacing = 1000;
            var rate = 44000;
            var duration = 1;

            // The buffer will have single-sample impulses written into it.
            var buffer = context.createBuffer(1, rate * duration, rate);
            var bufferData = buffer.getChannelData(0);

            var source = context.createBufferSource();
            source.buffer = buffer;

            var filter = context.createBiquadFilter();
            filter.type = "bandpass";
            filter.frequency.value = 450;
            filter.Q.value = 10;
            filter.gain.value = 40;

            var gain = context.createGain();
            gain.gain.value = 20;

            source.connect(filter);
            filter.connect(gain);
            gain.connect(context.destination);

            count = document.getElementById("count");

            var socket = new WebSocket("ws://" + window.location.hostname +  ":{{.Port}}/monitor-endpoint");
            socket.onmessage = function(message, x) {
                lastSecond = parseInt(message.data, 10);

                count.innerHTML = lastSecond.toString();

                for (x = 0; x < rate * duration; x++) {
                    bufferData[x] = 0;
                }

                for (x = 0; x < lastSecond; x++) {
                    pos = Math.floor(Math.random() * duration * (rate / spacing) * spacing);
                    bufferData[pos] = 1;
                };
            }

            source.onended = function(){console.log("END")}
            source.loop = true;
            source.start(context.currentTime);

          }
          catch(e) {
            alert('Web Audio API is not supported in this browser' + e);
          }
        }
        </script>
    </body>
</html>`
