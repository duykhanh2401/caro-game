<template>
	<span>
		<Icon :icon="icon" />
	</span>
</template>
<script lang="ts" setup>
import { Icon } from '@iconify/vue';
const props = defineProps({
	icon: {
		type: String,
		default: 'heroicons-outline:home'
	}
})
</script>
<style lang=""></style>
