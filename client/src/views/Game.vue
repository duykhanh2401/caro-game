<template>
	<div class="grid grid-cols-12 h-full">
		<div class="col-span-2 hidden md:block"></div>
		<div class="md:col-span-8 col-span-12 flex flex-col">
			<div class="py-2">
				<div class="px-3 container md:mx-auto">
					<div class="flex mx-2 flex-wrap">
						<div
							class="flex-shrink-0 flex-grow-0 flex items-center basis-auto w-1/2"
						>
							<div class="w-[40px]">
								<IConX v-if="currentRoom?.roomMasterFirst"></IConX>
								<IConO v-else></IConO>
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
							<div class="mr-2">0</div>
						</div>
						<div
							class="flex-shrink-0 flex-grow-0 basis-auto flex flex-row-reverse items-center w-1/2"
						>
							<div class="w-[40px]">
								<IConX v-if="!currentRoom?.roomMasterFirst"></IConX>
								<IConO v-else></IConO>
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
							<div class="ml-2">0</div>
						</div>
					</div>
				</div>
			</div>
			<div class="flex py-4 flex-grow justify-center items-center">
				<div class="game">
					<div class="board-wrapper">
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
			currentRoom?.roomMaster != '' &&
			(currentRoom?.guestReady == false || currentRoom?.masterReady == false)
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
import { ref, getCurrentInstance, onMounted } from 'vue';
import IConX from '@/assets/icons/x_icon.svg';
import IConO from '@/assets/icons/o_icon.svg';
import useRoomState from '../composable/useRoomState';
import useUserState, { type IUser } from '../composable/useUserState';
import useConnectGlobal from '../composable/useConnectGlobal';
const data = ref<Array<String>>([]);
const hoverCell = ref<number | null>();
const gridCount = 16;

const { currentRoom } = useRoomState();
const { users, me } = useUserState();
const { leaveRoom, guestReady, masterReady } = useConnectGlobal();

const currentClick = ref<string | null>();

function isMaster(): boolean {
	if (currentRoom.value?.roomMaster == me.value.id) {
		return true;
	}

	return false;
}

function GetMasterName(): IUser | undefined {
	return users.value?.find((user) => {
		return user.id == currentRoom.value?.roomMaster;
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
	data.value = new Array(gridCount * gridCount).fill(null);
}

function getSplitDataArr() {
	return new Array(Math.ceil(data.value.length / gridCount))
		.fill(null)
		.map((item, i) => data.value.slice(gridCount * i, gridCount * (i + 1)));
}

const getCellStyle = (index) => {
	const styles = {
		height: `${100 / gridCount}%`,
		width: `${100 / gridCount}%`,
	};

	// if (index == currentClick.value) {
	// 	styles.backgroundColor = '#ccc';
	// }

	// if (winnerRow) {
	//     const winIndex = winnerRow.findIndex(i => i === index);
	//     if (winIndex >= 0) styles.animationDelay = `${winIndex * (1 / gridCount)}s`;
	// }

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
		// victorious: winnerRow?.includes(index),
		clicked: String(index) == currentClick.value,
	};
};

const onClick = (e) => {
	console.log(e.target.dataset);

	if (e.target.classList.contains('cell') && e.target.dataset.value == 'empty') {
		if (currentClick.value == e.target.dataset.key) {
			data.value[Number(currentClick.value)] = 'x';
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
		const cell = data.value[index];
		if (index && cell === null) {
			// handleCellSet(index);
			// vibrate(25);
		}
		hoverCell.value = null;
	}
};

// const handleCellSet = index => {
//         setBoard(prevBoard => [...set(prevBoard, index, currentPlayer)]);
//         togglePlayer();
//     };

const getBoardIndex = (rowIndex, cellIndex) => rowIndex * gridCount + cellIndex;

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

.clicked {
	background-color: #eceff2 !important;
}

.dark .clicked {
	background-color: #384455 !important;
}

.game {
	height: 100%;
	aspect-ratio: 1 / 1;
	max-height: 80vw;
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
	padding: 1em 0;
	display: flex;
	justify-content: space-between;
	width: 100%;
	max-width: 40em;

	.panel-x-block,
	.panel-o-block,
	.info-content {
		display: flex;
		align-items: flex-end;
		font-size: 2.6em;

		> img {
			width: 1.2em;
			transition: transform 0.2s ease, opacity 0.2s ease;
		}
	}

	.panel-x-block,
	.panel-o-block {
		font-family: 'Roboto Mono', monospace;

		> img {
			opacity: 0.4;
			margin-left: 0.2em;
			transform-origin: left center;
		}
	}

	.panel-o-block > img {
		margin-left: 0;
		margin-right: 0.2em;
		transform-origin: right center;
	}

	&.x-move .panel-x-block,
	&.o-move .panel-o-block {
		> img {
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
