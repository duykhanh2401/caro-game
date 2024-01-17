import { ref } from 'vue';
import useWebSocket from './useWebSocket';

const isConnecting = ref(false);
const url = ref('ws://localhost:8080/ws/chat');
const ws = ref<WebSocket>();

enum ResponseEvents {
	ERROR,
	CONNECTED,
	TOPIC_ROOMS,
	ME_CHANGED_USERNAME,
	OTHER_CHANGED_USERNAME,
	ME_JOINED_CHAT,
	OTHER_JOINED_CHAT,
	ME_LEFT_CHAT,
	OTHER_LEFT_CHAT,
	ME_MESSAGE_SEND,
	OTHER_MESSAGE_SEND,
	OLD_MESSAGES,
	ME_CREATED_ROOM,
}

enum RequestEvents {
	GET_ROOMS,
	CHANGE_USERNAME,
	JOIN_ROOM,
	LEFT_CHAT,
	SEND_MESSAGE,
	GET_OLD_MESSAGES,
	CREATE_ROOM,
}

export default function UseConnectGlobal() {
	async function connectServer(): Promise<void> {
		if (isConnecting.value) return;

		const { connect, connections } = useWebSocket();
		isConnecting.value = true;
		try {
			await connect(url.value);
		} catch (err: any) {
			console.log(err);
			connectServer();
		} finally {
			isConnecting.value = false;
		}

		ws.value = connections.value.get(url.value);
	}

	async function addListener(ws: WebSocket | undefined) {
		if (!ws) return;

		ws.addEventListener('close', (e) => {
			console.log('wsChat(close):', e);
			if (e.code !== 1000) connectServer();
		});

		ws.addEventListener('error', (e) => {
			console.error('wsChat(error):', e);
		});

		ws.addEventListener('message', (e) => {
			const res = JSON.parse(e.data);
			switch (res.type) {
				case ResponseEvents.ERROR:
					console.log('wsChat(message): ResponseEvents.ERROR:');
					console.log(JSON.stringify(res, null, 2));
					break;
				case ResponseEvents.CONNECTED:
			}
		});
	}
}
