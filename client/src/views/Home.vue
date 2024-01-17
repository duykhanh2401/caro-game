<template>
	<FirstPage v-if="!me.isUpdateUserName" @changeUsername="changeUsername"></FirstPage>
	<Menu v-else></Menu>
</template>
<script setup lang="ts">
import Menu from './Menu.vue';
import Game from './Game.vue';
import FirstPage from './FirstPage.vue';

import useConnectGlobal from '../composable/useConnectGlobal';
import { onBeforeUnmount } from 'vue';
import useUserState from '../composable/useUserState';

const { connectServer, ws, changeUsername } = useConnectGlobal();
const { me } = useUserState();
connectServer();

onBeforeUnmount(() => {
	ws.value?.close(1000);
});
</script>
