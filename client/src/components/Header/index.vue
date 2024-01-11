<template>
	<header class="sticky top-0 z-[999]">
		<div
			class="app-header md:px-6 px-[15px] dark:bg-slate-800 shadow-base dark:shadow-base3 bg-white dark:border-b dark:border-slate-700 dark:border-opacity-60 md:py-6 py-3"
		>
			<div class="flex justify-between items-center h-full">
				<div class="flex items-center md:space-x-4 space-x-2 rtl:space-x-reverse">
					<button
						@click="sidebarCollapse = false"
						class="mr-5 rtl:ml-5 text-xl text-slate-900 dark:text-white"
						v-if="sidebarCollapse && width > 1280"
					>
						<Icon icon="akar-icons:arrow-right" />
					</button>
					<MobileLogo v-if="width < 1280" />
					<HandleMobileMenu v-if="width < 1280 && width > 768" />
				</div>
				<div
					class="nav-tools flex items-center lg:space-x-5 space-x-3 rtl:space-x-reverse"
				>
					<SwitchDark />
					<Profile v-if="width > 768" />
					<HandleMobileMenu v-if="width < 768" />
				</div>
			</div>
		</div>
	</header>
</template>
<script lang="ts" setup>
import SwitchDark from '@/components/Header/NavTools/SwitchDark.vue';
import Profile from '@/components/Header/NavTools/Profile.vue';
import MobileLogo from './NavTools/MobileLogo.vue';
import HandleMobileMenu from './NavTools/HandleMobileMenu.vue';
import { useThemeSetting } from '@/store/themeSetting.store';
import { storeToRefs } from 'pinia';
import { useWindow } from '@/composable/window';
const { width } = useWindow();
const themeSettingStore = useThemeSetting();
const { sidebarCollapse } = storeToRefs(themeSettingStore);
</script>
<style lang="scss" scoped>
.floating .app-header {
	@apply md:mx-6 md:my-8 mx-[15px] my-[15px] rounded-md;
}
</style>
