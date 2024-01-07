import { io } from 'socket.io-client';

export const socket = io('localhost:2401');

socket.on('connect', () => {
	console.log('Socket Connect');
});
