import { ref } from 'vue';
import type { IUser } from './useUserState';

export interface IRoom {
	id: string;
	name: string;
	master: string;
	masterWin: number;
	Guest: string;
	GuestWin: number;
	MasterFirst: boolean;
	IsMasterTurn: boolean;
	GuestReady: boolean;
	MasterReady: boolean;
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
const currentRoom = ref<IRoom>();
const users = ref<IUser[]>();

export default function useRoomState() {
	return {
		rooms,
		users,
		currentRoom,
		messages,
	};
}