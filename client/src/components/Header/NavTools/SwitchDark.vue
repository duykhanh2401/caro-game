<template>
	<span
		@click="themeSettingStore.toggleDarkTheme()"
		class="h-[28px] w-[28px] lg:h-[32px] lg:w-[32px] lg:bg-gray-500-f7 bg-slate-50 dark:bg-slate-900 lg:dark:bg-slate-900 dark:text-white text-slate-900 cursor-pointer text-[20px] flex flex-col items-center justify-center"
	>
		<DKIcon icon="heroicons-outline:moon" v-if="isDark" />
		<DKIcon icon="heroicons-outline:sun" v-else-if="!isDark" />
	</span>
</template>
<script setup lang="ts">
import { useThemeSetting } from '@/store/themeSetting.store';
import { storeToRefs } from 'pinia';
const themeSettingStore = useThemeSetting();
const { isDark } = storeToRefs(themeSettingStore);
</script>
