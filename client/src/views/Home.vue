<template>
	<div class="flex items-center justify-center pt-[5vh]">
		<div class="flex w-[280px] h-[140px]"><IConX /><IConO /></div>
	</div>
	<div class="flex flex-col items-center mt-[7vh]">
		<span
			class="select-none btn-primary mt-2 w-[60%] rounded-[999px] btn inline-flex justify-center"
			>Tạo Phòng</span
		>
		<span
			class="select-none btn-primary mt-2 w-[60%] rounded-[999px] btn inline-flex justify-center"
			>Vào Phòng</span
		>
		<span
			class="select-none btn-primary mt-2 w-[60%] rounded-[999px] btn inline-flex justify-center"
			>Danh Sách Phòng</span
		>
	</div>
</template>
<script setup lang="ts">
import IConX from '@/assets/icons/x_icon.svg';
import IConO from '@/assets/icons/o_icon.svg';
</script>
