<template>
	<div class="mobile-sidebar bg-white dark:bg-slate-800 shadow-base">
		<div class="logo-segment flex justify-between items-center px-4 py-6">
			<router-link :to="{ name: 'home' }">
				<img src="@/assets/images/logo/logo-c.svg" alt="" v-if="!isDark" />
				<img src="@/assets/images/logo/logo-c-white.svg" alt="" v-if="isDark" />
			</router-link>
			<span
				class="cursor-pointer text-slate-900 dark:text-white text-2xl"
				@click="themeSettingStore.toggleMobileSidebar()"
				><Icon icon="heroicons:x-mark"
			/></span>
		</div>
		<div class="sidebar-menu px-4 h-[calc(100%-100px)]" data-simplebar>
			<Navmenu :items="routerStore.routerSidebar" />
		</div>
	</div>
</template>
<script setup lang="ts">
import { Icon } from '@iconify/vue';
import { useThemeSetting } from '@/store/themeSetting.store';
import { storeToRefs } from 'pinia';
import Navmenu from './Navmenu.vue';
import { useRouterStore } from '@/store/router.store';
const routerStore = useRouterStore();
const themeSettingStore = useThemeSetting();
const { isDark } = storeToRefs(themeSettingStore);
</script>
<style lang="scss" scoped>
.mobile-sidebar {
	@apply fixed ltr:left-0 rtl:right-0 top-0   h-full   z-[9999]  w-[280px];
}
</style>
