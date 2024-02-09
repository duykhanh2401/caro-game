import { ref } from 'vue';
import type { IUser } from './useUserState';

export interface IRoom {
	id: string;
	name: string;
	master: string;
	masterWin: number;
	masterFirst: boolean;
	guest: string;
	guestWin: number;
	isMasterTurn: boolean;
	guestReady: boolean;
	masterReady: boolean;
	data: string[];
	winnerRow: number[];
	tempData: string[];
	gameEnd: boolean;
	turnClosest: number | null;
}

export interface IMessage {
	id: string;
	message: string;
	timestamp: number;
	user: IUser;
	type: MessageType;
}

export enum MessageType {
	NEW_MESSAGE,
	ME_CHANGED_USERNAME,
	OTHER_CHANGED_USERNAME,
	OTHER_LEFT,
	ME_JOINED,
	OTHER_JOINED,
}

const messages = ref<IMessage[]>([]);
const rooms = ref<IRoom[]>();
const currentRoom = ref<IRoom>({
	id: '',
	name: '',
	master: '',
	masterWin: 0,
	masterReady: false,
	masterFirst: true,
	guest: '',
	guestWin: 0,
	isMasterTurn: true,
	guestReady: false,
	data: [],
	winnerRow: [],
	tempData: [],
	gameEnd: false,
	turnClosest: null,
});
const users = ref<IUser[]>();

export default function useRoomState() {
	return {
		rooms,
		users,
		currentRoom,
		messages,
	};
}
