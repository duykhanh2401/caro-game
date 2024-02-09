<template>
	<div class="grid grid-cols-12 h-full">
		<div class="col-span-2 hidden md:block"></div>
		<div class="md:col-span-8 col-span-12 flex flex-col">
			<div class="py-2">
				<div class="px-3 container md:mx-auto">
					<div
						class="flex mx-2 flex-wrap status-panel"
						:class="getClassCurrentUserTurn()"
					>
						<div
							class="flex-shrink-0 flex-grow-0 flex items-center basis-auto w-1/2"
						>
							<div
								class="w-[40px] panel-x-block"
								v-if="currentRoom?.masterFirst"
							>
								<IConX></IConX>
							</div>
							<div class="w-[40px] panel-o-block" v-else>
								<IConO></IConO>
							</div>

							<div
								class="text-right min-w-0 flex-shrink text-lg flex flex-grow flex-col"
							>
								<span class="truncate">{{
									GetMasterName()?.username
								}}</span>
							</div>
							<div
								class="relative inline-block h-12 w-12 mx-2 !rounded-full object-cover object-center"
							>
								<img
									alt="John's Animal Avatar"
									:src="getRandomAvatar()"
								/>
							</div>
							<div class="mr-2">{{ currentRoom.masterWin }}</div>
						</div>
						<div
							class="flex-shrink-0 flex-grow-0 basis-auto flex flex-row-reverse items-center w-1/2"
						>
							<div
								class="w-[40px] panel-x-block"
								v-if="!currentRoom?.masterFirst"
							>
								<IConX></IConX>
							</div>
							<div class="w-[40px] panel-o-block" v-else>
								<IConO></IConO>
							</div>

							<div
								class="min-w-0 flex-shrink text-lg flex flex-grow flex-col"
							>
								<span class="truncate" v-if="GetGuestName()">{{
									GetGuestName()?.username
								}}</span>
								<span v-else class="loading truncate">
									Đang chờ đối thủ
								</span>
							</div>
							<div
								class="relative inline-block h-12 w-12 mx-2 !rounded-full object-cover object-center"
							>
								<img
									alt="John's Animal Avatar"
									:src="getRandomAvatar()"
								/>
							</div>
							<div class="ml-2">{{ currentRoom.guestWin }}</div>
						</div>
					</div>
				</div>
			</div>
			<div class="flex py-4 flex-grow justify-center items-center">
				<div class="game">
					<div class="board-wrapper" :class="currentRoom.gameEnd && 'ended'">
						<div class="reset-block" @click="reset">
							<span class="win-caption">Ấn để tiếp tục</span>
						</div>
						<table
							class="board"
							@click="onClick"
							@touchmove="onTouchMove"
							@touchstart="onTouchMove"
							@touchend="onTouchEnd"
						>
							<tbody>
								<tr
									v-for="(row, rowIndex) in getSplitDataArr()"
									:key="rowIndex"
								>
									<td
										v-for="(cell, cellIndex) in row"
										:key="cellIndex"
										:style="
											getCellStyle(
												getBoardIndex(rowIndex, cellIndex),
											)
										"
										:data-key="getBoardIndex(rowIndex, cellIndex)"
										:data-value="cell || 'empty'"
										class="cell"
										:class="
											getCellClassNames(
												getBoardIndex(rowIndex, cellIndex),
											)
										"
									>
										<IConX v-if="cell == 'x'" />
										<IConO v-else-if="cell == 'o'" />
									</td>
								</tr>
							</tbody>
						</table>
					</div>
				</div>
			</div>
			<div class="h-[60px] flex px-4 flex-row items-center whitespace-nowrap">
				<div>
					<DKButton
						btnClass="btn-primary btn-sm"
						@click="leaveRoom(currentRoom?.id)"
						>Huỷ bỏ ván chơi</DKButton
					>
				</div>
			</div>
		</div>
		<div class="col-span-2 hidden md:block"></div>
	</div>
	<DKModal
		:isOpen="
			currentRoom?.guest != '' &&
			currentRoom?.master != '' &&
			(currentRoom?.guestReady == false || currentRoom?.masterReady == false) &&
			!currentRoom.gameEnd
		"
		:centered="true"
		title="Sẵn Sàng"
	>
		<div class="flex justify-around">
			<div class="w-1/2 flex flex-col items-center">
				<div>
					<span class="truncate">{{ GetMasterName()?.username }}</span>
				</div>
				<DKButton
					v-if="isMaster()"
					class="mt-2"
					btnClass="btn-primary btn-sm"
					@click="masterReady(!currentRoom?.masterReady)"
					>{{
						currentRoom?.masterReady ? 'Huỷ sẵn sàng' : 'Sẵn sàng'
					}}</DKButton
				>
				<DKButton
					v-else
					isDisabled
					class="mt-2"
					:class="!currentRoom?.masterReady && `loading`"
					btnClass="btn-primary btn-sm"
					>{{
						currentRoom?.masterReady
							? 'Người chơi đã sẵn sàng'
							: 'Đang đợi người chơi sẵn sàng'
					}}</DKButton
				>
			</div>
			<div class="w-1/2 flex flex-col items-center">
				<div>
					<span class="truncate">{{ GetGuestName()?.username }}</span>
				</div>
				<DKButton
					v-if="!isMaster()"
					class="mt-2"
					btnClass="btn-primary btn-sm"
					@click="guestReady(!currentRoom?.guestReady)"
					>{{ currentRoom?.guestReady ? 'Huỷ sẵn sàng' : 'Sẵn sàng' }}</DKButton
				>
				<DKButton
					v-else
					isDisabled
					class="mt-2"
					:class="!currentRoom?.guestReady && `loading`"
					btnClass="btn-primary btn-sm "
					>{{
						currentRoom?.guestReady
							? 'Người chơi đã sẵn sàng'
							: 'Đang đợi người chơi sẵn sàng'
					}}</DKButton
				>
			</div>
		</div>
		<div class="flex justify-center mt-5">
			<DKButton
				class="mt-2"
				@click="leaveRoom(currentRoom?.id)"
				btnClass="btn-primary btn-sm"
				>Rời phòng</DKButton
			>
		</div>
	</DKModal>
</template>
<script setup lang="ts">
import { ref, getCurrentInstance, onMounted, watch } from 'vue';
import IConX from '@/assets/icons/x_icon.svg';
import IConO from '@/assets/icons/o_icon.svg';
import useRoomState from '../composable/useRoomState';
import useUserState, { type IUser } from '../composable/useUserState';
import useConnectGlobal from '../composable/useConnectGlobal';
import { fireConfetti } from '../utils/fireConfetti';
const hoverCell = ref<number | null>();
const gridCount = 15;

const { currentRoom } = useRoomState();
const { users, me } = useUserState();
const { leaveRoom, guestReady, masterReady, gameHandle } = useConnectGlobal();

const currentClick = ref<string | null>();

function isMaster(): boolean {
	if (currentRoom.value?.master == me.value.id) {
		return true;
	}

	return false;
}

function GetMasterName(): IUser | undefined {
	return users.value?.find((user) => {
		return user.id == currentRoom.value?.master;
	});
}

function GetGuestName(): IUser | undefined {
	return users.value?.find((user) => {
		return user.id == currentRoom.value?.guest;
	});
}

function getRandomAvatar(): string {
	console.log('Get Random avt');

	return (
		'https://anonymous-animals.azurewebsites.net/avatar/' +
		listAvatarRandom[Math.floor(Math.random() * listAvatarRandom.length)]
	);
}

function reset() {
	currentRoom.value.data = new Array(gridCount * gridCount).fill(null);
	currentRoom.value.gameEnd = false;
	currentRoom.value.winnerRow = [];
	currentRoom.value.turnClosest = null;
}

function getSplitDataArr() {
	return new Array(Math.ceil(currentRoom.value.data.length / gridCount))
		.fill(null)
		.map((item, i) =>
			currentRoom.value.data.slice(gridCount * i, gridCount * (i + 1)),
		);
}

const getCellStyle = (index) => {
	const styles = {
		height: `${100 / gridCount}%`,
		width: `${100 / gridCount}%`,
		animationDelay: '',
	};

	if (currentRoom.value.winnerRow.length > 0) {
		const winIndex = currentRoom.value.winnerRow.findIndex((i) => i === index);
		if (winIndex >= 0) styles.animationDelay = `${winIndex * (1 / gridCount)}s`;
	}

	return styles;
};

const getCellClassNames = (index: number) => {
	var sameRow = false;
	var sameColumn = false;
	if (hoverCell.value) {
		sameRow =
			Math.floor(index / gridCount) === Math.floor(hoverCell.value / gridCount);
		sameColumn = index % gridCount === hoverCell.value % gridCount;
	}

	return {
		hovered: index === hoverCell.value,
		highlighted: Number.isInteger(hoverCell.value) && (sameRow || sameColumn),
		clicked: String(index) == currentClick.value,
		victorious: currentRoom.value.winnerRow.includes(index),
		turnClosest: currentRoom.value.turnClosest == index,
	};
};

const onClick = (e) => {
	if (e.target.classList.contains('cell') && e.target.dataset.value == 'empty') {
		if (currentClick.value == e.target.dataset.key) {
			gameHandle(Number(currentClick.value));
		} else {
			currentClick.value = e.target.dataset.key;
		}
	}
};

const onTouchMove = (e) => {
	const { clientX, clientY } = e.touches[0];
	const el = document.elementFromPoint(clientX, clientY) as HTMLElement;
	if (el?.dataset.key) {
		hoverCell.value = el?.classList.contains('cell') ? +el?.dataset.key : null;
	}
};

const onTouchEnd = (e) => {
	const { clientX, clientY } = e.changedTouches[0];
	const el = document.elementFromPoint(clientX, clientY) as HTMLElement;
	const index = el?.dataset.key;
	if (index) {
		const cell = currentRoom.value.data[index];
		if (index && cell === null) {
			// handleCellSet(index);
			// vibrate(25);
		}
		hoverCell.value = null;
	}
};

watch(currentRoom, (newData, oldData) => {
	console.log('Watch', currentRoom.value.winnerRow);

	if (currentRoom.value.winnerRow.length) {
		fireConfetti(120, { y: 0.8, x: 1 });
		fireConfetti(60, { y: 0.8, x: 0 });
	}
});

const getBoardIndex = (rowIndex, cellIndex) => rowIndex * gridCount + cellIndex;

function getClassCurrentUserTurn() {
	if (
		(currentRoom.value?.masterFirst && currentRoom.value.isMasterTurn) ||
		(!currentRoom.value?.masterFirst && !currentRoom.value?.isMasterTurn)
	) {
		return 'x-move';
	} else {
		return 'o-move';
	}
}

reset();
const listAvatarRandom = [
	'Alligator',
	'Chipmunk',
	'Gopher',
	'Liger',
	'Quagga',
	'Anteater',
	'Chupacabra',
	'Grizzly',
	'Llama',
	'Rabbit',
	'Armadillo',
	'Cormorant',
	'Hedgehog',
	'Manatee',
	'Raccoon',
	'Auroch',
	'Coyote',
	'Hippo',
	'Mink',
	'Rhino',
	'Axolotl',
	'Crow',
	'Hyena',
	'Monkey',
	'Sheep',
	'Badger',
	'Dingo',
	'Ibex',
	'Moose',
	'Shrew',
	'Bat',
	'Dinosaur',
	'Ifrit',
	'Narwhal',
	'Skunk',
	'Beaver',
	'Dolphin',
	'Iguana',
	'Orangutan',
	'Squirrel',
	'Buffalo',
	'Duck',
	'Jackal',
	'Otter',
	'Tiger',
	'Camel',
	'Elephant',
	'Kangaroo',
	'Panda',
	'Turtle',
	'Capybara',
	'Ferret',
	'Koala',
	'Penguin',
	'Walrus',
	'Chameleon',
	'Fox',
	'Kraken',
	'Platypus',
	'Wolf',
	'Cheetah',
	'Frog',
	'Lemur',
	'Pumpkin',
	'Wolverine',
	'Chinchilla',
	'Giraffe',
	'Leopard',
	'Python',
	'Wombat',
];
</script>
<style lang="scss">
table {
	border-collapse: separate;
}

.dark .cell {
	background-color: #1e293b;
	border: 1.5px solid #0f172a;
}

.clicked,
.turnClosest {
	background-color: #ffeaa7 !important;
}

.dark .clicked,
.dark .turnClosest {
	background-color: #384455 !important;
}

.game {
	height: 100%;
	aspect-ratio: 1 / 1;
	max-height: 95vw;
	display: flex;
	flex-direction: column;
	align-items: center;
	justify-content: space-evenly;
	flex-wrap: wrap;
	font-size: 2vmin;
}

.board-wrapper {
	position: relative;
	aspect-ratio: 1 / 1;
	height: 100%;
	max-width: 100%;

	&.ended {
		.board {
			pointer-events: none;
			opacity: 0.4;
		}

		.reset-block {
			opacity: 1;
			transform: scale(1);
			pointer-events: auto;
		}
	}

	.board {
		width: 100%;
		height: 100%;
		margin: 0;
		flex-wrap: wrap;
		border-spacing: 0;
		touch-action: none;
	}

	.reset-block {
		position: absolute;
		width: 100%;
		height: 100%;
		display: flex;
		justify-content: center;
		align-items: center;
		opacity: 0;
		transform: scale(0.6);
		transition: transform 0.2s ease, opacity 0.2s ease;
		pointer-events: none;
		z-index: 1;

		.win-caption {
			font-size: 4em;
			text-shadow: 1px -1px 4px rgba(0, 0, 0, 0.3);
			font-weight: bold;
			user-select: none;
		}
	}
}

.cell {
	border: 1px solid rgba(0, 0, 0, 0.2);
	justify-content: center;
	padding: 0;
	align-items: center;
	background-color: #fff;
	transition: transform 0.2s ease, box-shadow 0.2s ease;
	user-select: none;

	> img {
		display: block;
		pointer-events: none;
	}

	&.victorious {
		animation: bounce 2s infinite;
	}
}

.status-panel {
	.panel-x-block,
	.panel-o-block,
	.info-content {
		display: flex;
		align-items: flex-end;
		font-size: 2.6em;

		> svg {
			width: 1.2em;
			transition: transform 0.2s ease, opacity 0.2s ease;
		}
	}

	.panel-x-block,
	.panel-o-block {
		font-family: 'Roboto Mono', monospace;

		> svg {
			opacity: 0.4;
			margin-left: 0.2em;
			transform-origin: left center;
		}
	}

	.panel-o-block > svg {
		margin-left: 0;
		margin-right: 0.2em;
		transform-origin: right center;
	}

	&.x-move .panel-x-block,
	&.o-move .panel-o-block {
		> svg {
			opacity: 1;
			transform: scale(1.6);
		}
	}
}

.menu {
	height: 100%;
	max-height: 100vw;
	max-width: 100vh;
	margin: 0 6em;
	display: flex;
	flex-direction: column;
	align-items: center;
	justify-content: space-evenly;
	font-size: 2vmin;

	.menu-items {
		width: 100%;
		flex-grow: 1;
		display: flex;
		flex-direction: column;
		justify-content: center;

		.menu-btn {
			padding: 0.6em;
			font-size: 2em;
			font-family: inherit;
			font-weight: bold;
		}
	}

	.game-configs {
		width: 100%;
		flex-grow: 2;
		display: flex;
		flex-direction: column;
		justify-content: space-evenly;

		.config-item {
			font-size: 2em;
		}
	}
}

.config-item {
	.menu-tab {
		font-family: inherit;
		font-size: 0.8em;

		.MuiTab-iconWrapper {
			font-size: 1em;

			> svg {
				width: 100%;
				max-width: 50%;
			}
		}
	}
}

@media (hover: none) {
	.cell {
		&:not(.victorious).hovered {
			position: relative;
			transform: scale(1.15);
			box-shadow: 0 0 4px rgba(0, 0, 0, 0.6);
			z-index: 1;
		}
		&:not(.hovered).highlighted {
			background-color: rgba(0, 0, 0, 0.1);
		}
	}
}

@media (hover: hover) {
	.cell {
		&:not(.victorious):hover {
			position: relative;
			transform: scale(1.15);
			box-shadow: 0 0 4px rgba(0, 0, 0, 0.6);
			z-index: 1;
		}
	}
}

@keyframes bounce {
	0% {
		transform: scale(1);
	}
	10% {
		transform: scale(1.2);
	}
	20% {
		transform: scale(1);
	}
	100% {
		transform: scale(1);
	}
}

.loading:after {
	overflow: hidden;
	display: inline-block;
	vertical-align: bottom;
	-webkit-animation: ellipsis steps(4, end) 900ms infinite;
	animation: ellipsis steps(4, end) 900ms infinite;
	content: '\2026';
	/* ascii code for the ellipsis character */
	width: 0px;
}

@keyframes ellipsis {
	to {
		width: 40px;
	}
}

@-webkit-keyframes ellipsis {
	to {
		width: 40px;
	}
}
</style>
