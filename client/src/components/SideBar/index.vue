<template>
	<div>
		<div
			:class="`sidebar-wrapper bg-white dark:bg-slate-800 shadow-base ${
				sidebarCollapse ? closeClass : openClass
			} ${isMouseHovered ? 'sidebar-hovered' : ''}`"
			@mouseenter="themeSettingStore.enterMouseHover()"
			@mouseleave="themeSettingStore.leaveMouseHover()"
		>
			<div
				:class="`logo-segment  border-none flex justify-between items-center bg-white dark:bg-slate-800 z-[9] py-6  sticky top-0   px-4  ${
					sidebarCollapse ? closeClass : openClass
				} 
        			${isMouseHovered ? 'logo-hovered' : ''} `"
			>
				<router-link :to="{ name: 'home' }">
					<img src="@/assets/images/logo/logo-c.svg" alt="" v-if="!isDark" />
					<img
						src="@/assets/images/logo/logo-c-white.svg"
						alt=""
						v-if="isDark"
					/>
				</router-link>
				<span
					class="cursor-pointer text-slate-900 dark:text-white text-2xl"
					v-if="!sidebarCollapse || isMouseHovered"
					@click="sidebarCollapse = !sidebarCollapse"
				>
					<div
						class="h-4 w-4 border-[1.5px] border-slate-900 dark:border-slate-700 rounded-full transition-all duration-150"
						:class="
							sidebarCollapse
								? ''
								: 'ring-2 ring-inset ring-offset-4 ring-black-900 dark:ring-slate-400 bg-slate-900 dark:bg-slate-400 dark:ring-offset-slate-700'
						"
					></div>
				</span>
			</div>
			<div
				class="h-[60px] absolute top-[80px] nav-shadow z-[1] w-full transition-all duration-200 pointer-events-none"
				:class="[shadowBase ? ' opacity-100' : ' opacity-0']"
			></div>
			<Simplebar class="sidebar-menu px-4 h-[calc(100%-80px)]">
				<Navmenu :items="routerStore.routerSidebar" />
			</Simplebar>
		</div>
	</div>
</template>
<script setup lang="js">
import { useThemeSetting } from '@/store/themeSetting.store';
import { storeToRefs } from 'pinia';
import { ref, onUnmounted } from 'vue';
const themeSettingStore = useThemeSetting();
const openClass = ref('w-[248px]');
const closeClass = ref('w-[72px] close_sidebar');
const shadowBase = ref(false);
import Simplebar from 'simplebar-vue';
import Navmenu from './Navmenu.vue';
import { useRouterStore } from '@/store/router.store';
const routerStore = useRouterStore()

const { isMouseHovered, sidebarCollapse, sidebarHidden, isDark } = storeToRefs(themeSettingStore);
onUnmounted(() => {
	console.log("Unmounted Navbar")
})
</script>
<style lang="scss">
.sidebar-wrapper {
	@apply fixed ltr:left-0 rtl:right-0 top-0   h-screen   z-[999];
	transition: width 0.2s cubic-bezier(0.39, 0.575, 0.565, 1);
	will-change: width;
}

.nav-shadow {
	background: linear-gradient(
		rgb(255, 255, 255) 5%,
		rgba(255, 255, 255, 75%) 45%,
		rgba(255, 255, 255, 20%) 80%,
		transparent
	);
}
.dark {
	.nav-shadow {
		background: linear-gradient(
			rgba(#1e293b, 100%) 5%,
			rgba(#1e293b, 75%) 45%,
			rgba(#1e293b, 20%) 80%,
			transparent
		);
	}
}
.sidebar-wrapper.sidebar-hovered {
	width: 248px !important;
}
.logo-segment.logo-hovered {
	width: 248px !important;
}
</style>
