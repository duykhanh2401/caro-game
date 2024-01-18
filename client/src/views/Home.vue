<template>
	<FirstPage v-if="!me.isUpdateUserName" @changeUsername="changeUsername"></FirstPage>
	<Game v-else-if="currentRoom"></Game>
	<Menu v-else></Menu>
</template>
<script setup lang="ts">
import Menu from './Menu.vue';
import Game from './Game.vue';
import FirstPage from './FirstPage.vue';

import useConnectGlobal from '../composable/useConnectGlobal';
import { onBeforeUnmount } from 'vue';
import useUserState from '../composable/useUserState';
import useRoomState from '../composable/useRoomState';

const { connectServer, ws, changeUsername } = useConnectGlobal();
const { me } = useUserState();
const { messages, currentRoom, rooms, users } = useRoomState();
connectServer();

onBeforeUnmount(() => {
	ws.value?.close(1000);
});
</script>
