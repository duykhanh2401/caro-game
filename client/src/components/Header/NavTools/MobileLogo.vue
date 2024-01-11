<template>
	<router-link :to="{ name: 'home' }">
		<img src="@/assets/images/logo/logo-c.svg" alt="" v-if="!isDark" />

		<img src="@/assets/images/logo/logo-c-white.svg" alt="" v-if="isDark" />
	</router-link>
</template>
<script setup lang="ts">
import { useThemeSetting } from '@/store/themeSetting.store';
import { storeToRefs } from 'pinia';
const themeSettingStore = useThemeSetting();
const { isDark } = storeToRefs(themeSettingStore);
</script>
