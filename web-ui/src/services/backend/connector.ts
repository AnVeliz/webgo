export const ConnectToBackendWebSocket = (onMessage: (msg: string) => void) => {
    if ("WebSocket" in window) {
        const ws = new WebSocket("ws://localhost:8089/");

        ws.onopen = () => {
            console.log("websocket connected");
        };

        ws.onmessage = (evt) => {
            console.log("message received");
            const receivedMsg = evt.data;
            onMessage(receivedMsg)
          };

        ws.onclose = () => {
            window.close();
            console.log("websocket disconnected");
        };

        ws.onerror = () => {
            window.close();
        };
    } else {
        alert("WebSocket is NOT supported by your Browser!");
    }
}