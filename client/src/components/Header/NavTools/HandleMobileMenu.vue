<template>
	<div
		class="cursor-pointer text-slate-900 dark:text-white text-2xl"
		@click="themeSettingStore.toggleMobileSidebar()"
	>
		<DKIcon icon="heroicons-outline:menu-alt-3" />
	</div>
</template>
<script setup lang="ts">
import { useThemeSetting } from '@/store/themeSetting.store';
const themeSettingStore = useThemeSetting();
</script>
