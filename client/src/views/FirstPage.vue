<template>
	<DKModal
		:isOpen="true"
		:centered="true"
		:disableClose="true"
		title="Chào mừng bạn đến với game cờ Caro"
	>
		<DKTextInput
			label="Vui lòng nhập tên của bạn:"
			placeholder="Tên của bạn..."
			v-model="username"
		></DKTextInput>
		<DKCheckBox label="Remember me" v-model="isSave" :checked="isSave" />
		<div class="flex justify-end">
			<DKButton
				btnClass="btn-primary btn-sm"
				@click="$emit('changeUsername', username)"
				>Gửi</DKButton
			>
		</div>
	</DKModal>
</template>
<script setup lang="ts">
import { ref, watch } from 'vue';

const isSave = ref(false);
const username = ref('');
const checkCache = window.localStorage.getItem('saveUsername');
if (checkCache) {
	isSave.value = true;
}
watch(isSave, () => {
	if (isSave.value) {
		window.localStorage.setItem('saveUsername', 'true');
	} else {
		window.localStorage.removeItem('saveUsername');
	}
});
</script>
