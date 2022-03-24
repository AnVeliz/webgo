import App from './App.svelte';
import { ConnectToBackendWebSocket } from "./services/backend/connector";

const app = new App({
	target: document.body,
	props: {
		timeText: 'Hello World'
	}
});

ConnectToBackendWebSocket((msg: string) => {
	app.$set({ timeText: new Date(msg).toLocaleTimeString() });
});

export default app;