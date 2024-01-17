import { ref } from 'vue';

export interface IUser {
	id: string;
	username: string;
	avatar: string;
	doneLoading?: boolean;
	isUpdateUserName: boolean;
}

const users = ref<IUser[]>([]);
const me = ref<IUser>({
	id: '',
	username: '',
	avatar: 'https://picsum.photos/56/56',
	isUpdateUserName: false,
});
export default function useUserState() {
	return {
		users,
		me,
	};
}
