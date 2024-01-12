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
		></DKTextInput>
		<DKCheckBox label="Remember me" checked />
		<div class="flex justify-end">
			<DKButton btnClass="btn-primary btn-sm">Gửi</DKButton>
		</div>
	</DKModal>
</template>
<script setup lang="ts"></script>
