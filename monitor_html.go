package main

var MONITOR_HTML string = `<html>
    <head>
        <link href='http://fonts.googleapis.com/css?family=Cagliostro' rel='stylesheet' type='text/css' />
        <title>Geigor Monitor</title>
        <style type="text/css">
            body {
                font-family: 'Cagliostro', sans-serif;
                margin: 0px;
                padding: 0px;
                color: 2b3a42;
                background-color: #fff9f9;
            }

            a {
                text-decoration: underline;
                color: inherit;
            }

            #container {
                border: 1px solid #efe9e9;
                padding: 10px;
                margin: 20px;
                width: 415px;
            }

            h2 {
                margin: 0px;
            }

            .footnote {
               font-size: 0.8em;
            }
        </style>
    </head>

    <body>
        <div id="container">
            <h2>Geigor tracking</h2>
            <p>Current count: <span id="count">unknown</span>.</p>
            <canvas width="400" height="4" id="canvas-buffer"></canvas>
            <input type="range" id="ssh" style="width: 50px" max="40"><label for="ssh">vol</label>
            <p class="footnote"><a href="http://afandian.com/geigor">Geigor</a>. This uses HTML5 tech not supported on all browsers. Yet.</p>
        </div>

        <script type="text/javascript">
        var context;
        window.addEventListener('load', init, false);
        function init() {
          try {
            window.AudioContext = window.AudioContext||window.webkitAudioContext;
            context = new AudioContext();

            // Spacing of impulses in the audio buffer.
            var spacing = 1000;
            var rate = 44000;

            // Duration of ringbuffer.
            var duration = 1;
            var samplesInBuffer = rate * duration;

            // The buffer will have single-sample impulses written into it.
            var buffer = context.createBuffer(1, samplesInBuffer, rate);
            var bufferData = buffer.getChannelData(0);

            var source = context.createBufferSource();
            source.buffer = buffer;
            source.start(0);

            // Filter to make the impulses into 'knock' sounds.
            var filter = context.createBiquadFilter();
            filter.type = "bandpass";
            filter.frequency.value = 450;
            filter.Q.value = 10;
            filter.gain.value = 40;

            var gain = context.createGain();
            gain.gain.value = 20;

            // FFT analyer, but we're only going to be using its window sampling.
            // This will detect impulses as they're played and plot them.
            var sampleSize = 2048;
            var analyser = context.createAnalyser();
            analyser.fftSize = sampleSize;
            analyser.smoothingTimeConstant = 1;
            var fftBuffer = new Uint8Array(sampleSize);

            // Wire it all up.
            source.connect(analyser);
            analyser.connect(filter);
            filter.connect(gain);
            gain.connect(context.destination);

            // HTML element to report the current count.
            var count = document.getElementById("count");

            // Ssh.
            var sshButton = document.getElementById("ssh");
            sshButton.onchange = function(event) {
                gain.gain.value = event.target.valueAsNumber;
            }

            // Canvas.
            var canvasBuffer = document.getElementById("canvas-buffer");
            var contextBuffer = canvasBuffer.getContext("2d");
            var canvasWidth = 400;
            var canvasHeight = 4;
            var shiftWidth = canvasWidth - 1;
            var shiftAmount = 2;
            var shiftDeleteX = canvasWidth - shiftAmount;

            // Set pixels on the right hand side of the buffer when there's an impulse.
            setInterval(function(){
                analyser.getByteTimeDomainData(fftBuffer);
                for (var i = 0; i < sampleSize; i++) {
                    if (fftBuffer[i] != 128)
                    {
                        contextBuffer.fillRect(395, 0, 4, 4);
                    }
                }
            }, 50);

            // Keep shifting the data in the canvas left.
            function shift() {
                var imageData = contextBuffer.getImageData(
                    1, 0,
                    shiftWidth, canvasHeight
                    );

                contextBuffer.putImageData(imageData, 0, 0);

                contextBuffer.clearRect(shiftDeleteX, 0, shiftAmount, canvasHeight);
            }

            setInterval(shift, 10);

            var socket = new WebSocket("ws://" + window.location.hostname +  ":{{.Port}}/monitor-endpoint");
            socket.onmessage = function(message, x) {
                lastSecond = parseInt(message.data, 10);

                count.innerHTML = lastSecond.toString();

                for (x = 0; x < samplesInBuffer; x++) {
                    bufferData[x] = 0;
                }

                var pos;
                for (x = 0; x < lastSecond; x++) {
                    pos = Math.floor(Math.random() * duration * (rate / spacing));
                    bufferData[pos * spacing] = 1;
                };
            }

            source.start(0);
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
