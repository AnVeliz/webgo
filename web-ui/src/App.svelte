<script lang="ts">
  // Forbid DevTools
  window.addEventListener("keydown", function (event) {
    if (event.keyCode == 116) {
      // block F5 (Refresh)
      event.preventDefault();
      event.stopPropagation();
      return false;
    } else if (event.keyCode == 122) {
      // block F11 (Fullscreen)
      event.preventDefault();
      event.stopPropagation();
      return false;
    } else if (event.keyCode == 123) {
      // block F12 (DevTools)
      event.preventDefault();
      event.stopPropagation();
      return false;
    } else if (event.ctrlKey && event.shiftKey && event.keyCode == 73) {
      // block Strg+Shift+I (DevTools)
      event.preventDefault();
      event.stopPropagation();
      return false;
    } else if (event.ctrlKey && event.shiftKey && event.keyCode == 74) {
      // block Strg+Shift+J (Console)
      event.preventDefault();
      event.stopPropagation();
      return false;
    }
  });

  window.addEventListener("contextmenu", (e) => e.preventDefault());
  window.addEventListener("selectstart", function (e) {
    e.preventDefault();
  });

  // Backend connection
  function ConnectToBackendWebSocket() {
    if ("WebSocket" in window) {
      //alert("WebSocket is supported by your Browser!");
      // Let us open a web socket
      var ws = new WebSocket("ws://localhost:8089/");
      window["App"] = {};
      window["App"].websocket = ws;

      ws.onopen = function () {
        // Web Socket is connected, send data using send()
        //ws.send("Message to send");
        //alert("Message is sent...");
        console.log("websocket connected");
      };

      ws.onmessage = function (evt) {
        console.log("message received");
        const element = document.getElementById("helloworldtxt");
        if (!element) {
          return;
        }
        var received_msg = evt.data;
        element.innerText = new Date(received_msg).toLocaleTimeString();

        //alert("Message is received..." + received_msg);
      };

      ws.onclose = function () {
        // websocket is closed.
        //alert("Connection is closed...");
        window.close();
        console.log("websocket disconnected");
      };

      ws.onerror = function () {
        window.close();
      };
    } else {
      // The browser doesn't support WebSocket
      alert("WebSocket is NOT supported by your Browser!");
    }
  }
  ConnectToBackendWebSocket();
</script>

<main>
  <div class="sign">
    <span class="fast-flicker" id="helloworldtxt">Hello!</span>
  </div>
</main>

<style>
  @font-face {
    font-family: Clip;
    src: url("https://acupoftee.github.io/fonts/Clip.ttf");
  }
  .sign {
    position: absolute;
    display: flex;
    justify-content: center;
    align-items: center;
    width: 50%;
    height: 50%;
    background-image: radial-gradient(
      ellipse 50% 35% at 50% 50%,
      #6b1839,
      transparent
    );
    transform: translate(-50%, -50%);
    letter-spacing: 2;
    left: 50%;
    top: 50%;
    font-family: "Clip";
    text-transform: uppercase;
    font-size: 6em;
    color: #ffe6ff;
    text-shadow: 0 0 0.6rem #ffe6ff, 0 0 1.5rem #ff65bd,
      -0.2rem 0.1rem 1rem #ff65bd, 0.2rem 0.1rem 1rem #ff65bd,
      0 -0.5rem 2rem #ff2483, 0 0.5rem 3rem #ff2483;
    animation: shine 2s forwards, flicker 3s infinite;
  }

  @keyframes blink {
    0%,
    22%,
    36%,
    75% {
      color: #ffe6ff;
      text-shadow: 0 0 0.6rem #ffe6ff, 0 0 1.5rem #ff65bd,
        -0.2rem 0.1rem 1rem #ff65bd, 0.2rem 0.1rem 1rem #ff65bd,
        0 -0.5rem 2rem #ff2483, 0 0.5rem 3rem #ff2483;
    }
    28%,
    33% {
      color: #ff65bd;
      text-shadow: none;
    }
    82%,
    97% {
      color: #ff2483;
      text-shadow: none;
    }
  }

  .flicker {
    animation: shine 2s forwards, blink 3s 2s infinite;
  }

  .fast-flicker {
    animation: shine 2s forwards, blink 10s 1s infinite;
  }

  @keyframes shine {
    0% {
      color: #6b1839;
      text-shadow: none;
    }
    100% {
      color: #ffe6ff;
      text-shadow: 0 0 0.6rem #ffe6ff, 0 0 1.5rem #ff65bd,
        -0.2rem 0.1rem 1rem #ff65bd, 0.2rem 0.1rem 1rem #ff65bd,
        0 -0.5rem 2rem #ff2483, 0 0.5rem 3rem #ff2483;
    }
  }

  @keyframes flicker {
    from {
      opacity: 1;
    }

    4% {
      opacity: 0.9;
    }

    6% {
      opacity: 0.85;
    }

    8% {
      opacity: 0.95;
    }

    10% {
      opacity: 0.9;
    }

    11% {
      opacity: 0.922;
    }

    12% {
      opacity: 0.9;
    }

    14% {
      opacity: 0.95;
    }

    16% {
      opacity: 0.98;
    }

    17% {
      opacity: 0.9;
    }

    19% {
      opacity: 0.93;
    }

    20% {
      opacity: 0.99;
    }

    24% {
      opacity: 1;
    }

    26% {
      opacity: 0.94;
    }

    28% {
      opacity: 0.98;
    }

    37% {
      opacity: 0.93;
    }

    38% {
      opacity: 0.5;
    }

    39% {
      opacity: 0.96;
    }

    42% {
      opacity: 1;
    }

    44% {
      opacity: 0.97;
    }

    46% {
      opacity: 0.94;
    }

    56% {
      opacity: 0.9;
    }

    58% {
      opacity: 0.9;
    }

    60% {
      opacity: 0.99;
    }

    68% {
      opacity: 1;
    }

    70% {
      opacity: 0.9;
    }

    72% {
      opacity: 0.95;
    }

    93% {
      opacity: 0.93;
    }

    95% {
      opacity: 0.95;
    }

    97% {
      opacity: 0.93;
    }

    to {
      opacity: 1;
    }
  }
</style>
