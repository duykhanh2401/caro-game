<template>
	<div class="grid grid-cols-12 h-full">
		<div class="col-span-2 hidden md:block"></div>
		<div class="md:col-span-8 col-span-12">
			<div class="pb-[64px] relative h-[calc(100%-64px)]">
				<div
					class="flex items-center justify-center absolute top-1/4 left-0 right-0"
				>
					<div class="flex w-[280px] h-[140px]"><IConX /><IConO /></div>
				</div>
				<div
					class="flex flex-col items-center absolute bottom-[40px] left-0 right-0"
				>
					<span
						class="select-none btn-primary mt-2 w-[60%] rounded-[999px] btn inline-flex justify-center cursor-pointer"
						@click="openCreateRoom = true"
						>Tạo Phòng</span
					>
					<span
						class="select-none btn-primary mt-2 w-[60%] rounded-[999px] btn inline-flex justify-center cursor-pointer"
						@click="openJoinRoom = true"
						>Vào Phòng</span
					>
					<span
						class="select-none btn-primary mt-2 w-[60%] rounded-[999px] btn inline-flex justify-center cursor-pointer"
						@click="openModalGetRooms"
						>Danh Sách Phòng</span
					>
				</div>
			</div>
			<div
				class="nav-bottom fixed bottom-0 z-50 md:w-[50vh] md:left-1/2 left-0 h-16 w-full border-gray-200 dark:border-gray-600"
			>
				<div class="grid h-full grid-cols-3 font-medium">
					<button
						type="button"
						class="inline-flex flex-col items-center justify-center px-5 hover:bg-gray-50 dark:hover:bg-gray-800 group"
					>
						<svg
							class="w-5 h-5 mb-2 text-gray-500 dark:text-gray-400 group-hover:text-blue-600 dark:group-hover:text-blue-500"
							aria-hidden="true"
							xmlns="http://www.w3.org/2000/svg"
							fill="currentColor"
							viewBox="0 0 20 20"
						>
							<path
								d="m19.707 9.293-2-2-7-7a1 1 0 0 0-1.414 0l-7 7-2 2a1 1 0 0 0 1.414 1.414L2 10.414V18a2 2 0 0 0 2 2h3a1 1 0 0 0 1-1v-4a1 1 0 0 1 1-1h2a1 1 0 0 1 1 1v4a1 1 0 0 0 1 1h3a2 2 0 0 0 2-2v-7.586l.293.293a1 1 0 0 0 1.414-1.414Z"
							/>
						</svg>
						<span
							class="text-xs text-gray-500 dark:text-gray-400 group-hover:text-blue-600 dark:group-hover:text-blue-500"
							>Trang chủ</span
						>
					</button>

					<button
						type="button"
						class="inline-flex flex-col items-center justify-center px-5 hover:bg-gray-50 dark:hover:bg-gray-800 group"
					>
						<svg
							class="w-5 h-5 mb-2 text-gray-500 dark:text-gray-400 group-hover:text-blue-600 dark:group-hover:text-blue-500"
							aria-hidden="true"
							xmlns="http://www.w3.org/2000/svg"
							fill="none"
							viewBox="0 0 20 20"
						>
							<path
								stroke="currentColor"
								stroke-linecap="round"
								stroke-linejoin="round"
								stroke-width="2"
								d="M4 12.25V1m0 11.25a2.25 2.25 0 0 0 0 4.5m0-4.5a2.25 2.25 0 0 1 0 4.5M4 19v-2.25m6-13.5V1m0 2.25a2.25 2.25 0 0 0 0 4.5m0-4.5a2.25 2.25 0 0 1 0 4.5M10 19V7.75m6 4.5V1m0 11.25a2.25 2.25 0 1 0 0 4.5 2.25 2.25 0 0 0 0-4.5ZM16 19v-2"
							/>
						</svg>
						<span
							class="text-xs text-gray-500 dark:text-gray-400 group-hover:text-blue-600 dark:group-hover:text-blue-500"
							>Cài Đặt</span
						>
					</button>
					<button
						type="button"
						class="inline-flex flex-col items-center justify-center px-5 hover:bg-gray-50 dark:hover:bg-gray-800 group"
					>
						<SwitchDark />
						<span
							class="text-xs text-gray-500 dark:text-gray-400 group-hover:text-blue-600 dark:group-hover:text-blue-500"
							>Chế độ tối</span
						>
					</button>
				</div>
			</div>
		</div>
		<div class="col-span-2 hidden md:block"></div>
	</div>
	<DKModal v-model:isOpen="openCreateRoom" :centered="true" title="Tạo phòng mới">
		<DKTextInput
			label="Tên phòng:"
			placeholder="Tên của bạn..."
			v-model="roomName"
		></DKTextInput>
		<div class="flex justify-end">
			<DKButton btnClass="btn-primary btn-sm" @click="createRoom(roomName)"
				>Gửi</DKButton
			>
		</div>
	</DKModal>
	<DKModal v-model:isOpen="openJoinRoom" :centered="true" title="Vào phòng">
		<span>Vùi lòng nhập ID của phòng:</span>
		<PincodeInput
			class="justify-center mt-4"
			v-model="pincode"
			autofocus
			:digits="6"
			input-class="block w-9 h-9 py-3 text-sm font-extrabold text-center text-gray-900 bg-white border border-gray-300 rounded-lg focus:ring-primary-500 focus:border-primary-500 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-primary-500 dark:focus:border-primary-500 mx-1"
		/>
	</DKModal>
	<DKModal v-model:isOpen="openListRoom" title="Danh sách phòng">
		<ul
			v-for="room in rooms"
			:key="room.id"
			class="flex flex-col justify-end text-start -space-y-px"
		>
			<li
				class="flex mb-2 items-center gap-x-2 p-3 text-sm bg-white border text-gray-800 first:rounded-t-lg first:mt-0 last:rounded-b-lg dark:bg-slate-900 dark:border-gray-700 dark:text-gray-200"
			>
				<div class="w-full flex justify-between truncate">
					<div class="flex-1">
						<div class="me-3 flex-1 text-xl font-bold">
							{{ room.name }}
						</div>
						<div class="me-3 flex-1 mt-1">
							{{ room.id }}
						</div>
					</div>
					<button
						type="button"
						class="flex items-center gap-x-2 text-gray-500 hover:text-blue-500 whitespace-nowrap"
						@click="joinRoom(room.id)"
					>
						Vào phòng
					</button>
				</div>
			</li>
		</ul>
	</DKModal>
</template>
<script setup lang="ts">
import IConX from '@/assets/icons/x_icon.svg';
import IConO from '@/assets/icons/o_icon.svg';
import SwitchDark from '@/components/Header/NavTools/SwitchDark.vue';
import { ref } from 'vue';
import PincodeInput from '@/components/PincodeInput/index.vue';
import useRoomState from '../composable/useRoomState';
import useConnectGlobal from '../composable/useConnectGlobal';
const openCreateRoom = ref(false);
const openJoinRoom = ref(false);
const openListRoom = ref(false);
const { createRoom, me, getRooms, joinRoom } = useConnectGlobal();
const { rooms } = useRoomState();

const pincode = ref('');
const roomName = ref('Phòng của ' + me.value.username);
const isLoadingGetRooms = ref(false);

function openModalGetRooms() {
	openListRoom.value = true;
	getRooms();
}

function closeAllModal() {
	openCreateRoom.value = false;
	openJoinRoom.value = false;
	openListRoom.value = false;
}
</script>
<style>
.nav-bottom {
	transform: translate(-50%, 0%);
	left: 50%;
}
</style>
