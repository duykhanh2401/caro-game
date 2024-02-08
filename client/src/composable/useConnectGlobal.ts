import { ref } from 'vue';
import useWebSocket from './useWebSocket';
import useUserState from './useUserState';
import { useToast } from 'vue-toastification';
import useRoomState, { MessageType } from './useRoomState';
const isConnecting = ref(false);
const url = ref('ws://localhost:2401/ws/caro');
const ws = ref<WebSocket>();

enum ResponseEvents {
	ERROR,
	CONNECTED,
	ME_GET_ROOMS,
	ME_CHANGED_USERNAME,
	OTHER_CHANGED_USERNAME,
	ME_JOINED_CHAT,
	OTHER_JOINED_CHAT,
	ME_LEFT_ROOM,
	OTHER_LEFT_ROOM,
	ME_MESSAGE_SEND,
	OTHER_MESSAGE_SEND,
	OLD_MESSAGES,
	ME_CREATED_ROOM,
	GUEST_LEAVE_ROOM,
	ME_TO_ROOM_MASTER,
	GUEST_READY_RESPONSE,
	MASTER_READY_RESPONSE,
	GAME_HANDLE_RESPONSE,
	GAME_END,
}

enum RequestEvents {
	GET_ROOMS,
	CHANGE_USERNAME,
	JOIN_ROOM,
	LEFT_ROOM,
	SEND_MESSAGE,
	GET_OLD_MESSAGES,
	CREATE_ROOM,
	GUEST_READY,
	MASTER_READY,
	GAME_HANDLE_REQUEST,
}

export default function useConnectGlobal() {
	const { me, users } = useUserState();
	const { currentRoom, rooms, messages } = useRoomState();
	const toast = useToast();
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
		addListener(ws.value);
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
			console.log('Current Room: ', currentRoom.value);
			const res = JSON.parse(e.data);
			switch (res.type) {
				case ResponseEvents.ERROR:
					console.log('wsChat(message): ResponseEvents.ERROR:');
					console.log(JSON.stringify(res, null, 2));

					toast.error(res.error.message);
					break;
				case ResponseEvents.CONNECTED:
					console.log('ws Connected');
					console.log(res);
					me.value = res.body.data;
					currentRoom.value = {
						id: '',
						name: '',
						master: '',
						masterWin: 0,
						guest: '',
						guestWin: 0,
						roomMasterFirst: true,
						isMasterTurn: true,
						guestReady: false,
						masterReady: false,
						roomMaster: '',
						data: [],
						winnerRow: [],
						tempData: [],
						gameEnd: false,
					};
					const username = window.localStorage.getItem('username');
					if (username) {
						changeUsername(username);
					}
					changeUsername();
					break;
				case ResponseEvents.ME_CHANGED_USERNAME:
					me.value = { ...res.body.data };
					me.value.isUpdateUserName = true;
					console.log('Me:', me.value, res.body.data);
					toast.success(`Chào mừng: ${me.value.username}`);
					const isSaveUsername = window.localStorage.getItem('saveUsername');
					if (isSaveUsername) {
						window.localStorage.setItem('username', me.value.username);
					}
					break;
				case ResponseEvents.ME_CREATED_ROOM:
					currentRoom.value = { ...currentRoom.value, ...res.body.room };

					users.value.push(res.body.user);
					messages.value = [
						...messages.value,
						{
							message: res.body.message,
							timestamp: Date.now(),
							id: '123',
							type: MessageType.ME_JOINED,
							user: me.value,
						},
					];
					toast.success('Bạn đã tạo phòng thành công !!!');
					break;
				case ResponseEvents.ME_GET_ROOMS:
					rooms.value = res.body.data;
					break;
				case ResponseEvents.ME_JOINED_CHAT:
					currentRoom.value = { ...currentRoom.value, ...res.body.data.room };
					users.value.push(...res.body.data.user);

					break;
				case ResponseEvents.OTHER_JOINED_CHAT:
					currentRoom.value = { ...currentRoom.value, ...res.body.data.room };
					users.value.push(res.body.data.user);
					toast.info(res.body.message);

					break;
				case ResponseEvents.GUEST_LEAVE_ROOM:
					toast.info(res.body.message);
					currentRoom.value = { ...currentRoom.value, ...res.body.room };
					break;
				case ResponseEvents.ME_LEFT_ROOM:
					toast.info(res.body.message);
					currentRoom.value = {
						id: '',
						name: '',
						master: '',
						masterWin: 0,
						guest: '',
						guestWin: 0,
						roomMasterFirst: true,
						isMasterTurn: true,
						guestReady: false,
						masterReady: false,
						roomMaster: '',
						data: [],
						winnerRow: [],
						tempData: [],
						gameEnd: false,
					};
					break;
				case ResponseEvents.GUEST_READY_RESPONSE:
					console.log('GUEST READY !!!!!');

					if (currentRoom.value) {
						currentRoom.value.guestReady = res.body.isReady || false;
					}
					break;
				case ResponseEvents.MASTER_READY_RESPONSE:
					console.log('MASTER READY !!!!!');

					if (currentRoom.value) {
						currentRoom.value.masterReady = res.body.isReady || false;
					}
					break;
				case ResponseEvents.GAME_HANDLE_RESPONSE:
					if (res.body.data.isXTurn) {
						currentRoom.value!.data[res.body.data.index] = 'x';
					} else {
						currentRoom.value!.data[res.body.data.index] = 'o';
					}

					currentRoom.value!.isMasterTurn = res.body.data.isMasterTurn;
					break;
				case ResponseEvents.GAME_END:
					currentRoom.value = { ...currentRoom.value, ...res.body.data.room };
					currentRoom.value.winnerRow = res.body.data.winnerRow;
			}
		});
	}

	function changeUsername(username?: string) {
		if (!username) return;

		try {
			ws.value?.send(
				JSON.stringify({
					type: RequestEvents.CHANGE_USERNAME,
					body: {
						username,
					},
				}),
			);
		} catch (error) {
			console.log(error);
		}
	}

	function createRoom(roomName?: string) {
		if (!roomName) return;
		try {
			ws.value?.send(
				JSON.stringify({
					type: RequestEvents.CREATE_ROOM,
					body: {
						roomName,
					},
				}),
			);
		} catch (error) {
			console.log(error);
		}
	}

	function getRooms() {
		try {
			ws.value?.send(
				JSON.stringify({
					type: RequestEvents.GET_ROOMS,
				}),
			);
		} catch (error) {
			console.log(error);
		}
	}

	function leaveRoom(roomID?: string) {
		if (!roomID) return;
		try {
			ws.value?.send(
				JSON.stringify({
					type: RequestEvents.LEFT_ROOM,
					body: {
						roomID,
					},
				}),
			);
		} catch (error) {}
	}

	function joinRoom(roomID: string) {
		if (!roomID) return;
		try {
			ws.value?.send(
				JSON.stringify({
					type: RequestEvents.JOIN_ROOM,
					body: {
						roomID,
					},
				}),
			);
		} catch (error) {}
	}

	function gameHandle(index: number) {
		try {
			ws.value?.send(
				JSON.stringify({
					type: RequestEvents.GAME_HANDLE_REQUEST,
					body: {
						roomID: currentRoom.value?.id,
						index: index,
					},
				}),
			);
		} catch (error) {}
	}

	function guestReady(ready: boolean) {
		try {
			ws.value?.send(
				JSON.stringify({
					type: RequestEvents.GUEST_READY,
					body: {
						roomID: currentRoom.value?.id,
						isReady: ready,
					},
				}),
			);
		} catch (error) {}
	}

	function masterReady(ready: boolean) {
		try {
			ws.value?.send(
				JSON.stringify({
					type: RequestEvents.MASTER_READY,
					body: {
						roomID: currentRoom.value?.id,
						isReady: ready,
					},
				}),
			);
		} catch (error) {}
	}

	return {
		connectServer,
		changeUsername,
		createRoom,
		getRooms,
		joinRoom,
		leaveRoom,
		guestReady,
		masterReady,
		gameHandle,
		ws,
		url,
		me,
	};
}
