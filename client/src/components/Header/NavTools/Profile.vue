<template>
	<Dropdown class-menu-items=" w-[180px] top-[58px] ">
		<div class="flex items-center">
			<div class="flex-1 mr-[10px]">
				<div class="lg:h-8 lg:w-8 h-7 w-7 rounded-full">
					<img
						:src="profileImg"
						alt=""
						class="block w-full h-full object-cover rounded-full"
					/>
				</div>
			</div>
			<div
				class="flex-none text-slate-600 dark:text-white text-sm font-normal items-center lg:flex hidden overflow-hidden text-ellipsis whitespace-nowrap"
			>
				<span
					class="overflow-hidden text-ellipsis whitespace-nowrap w-[85px] block"
					>Duy Khánh</span
				>
				<span class="text-base inline-block ml-[10px]"
					><DKIcon icon="heroicons-outline:chevron-down"></DKIcon
				></span>
			</div>
		</div>
		<template #menus>
			<MenuItem v-slot="{ active }" v-for="(item, i) in ProfileMenu" :key="i">
				<div
					type="button"
					:class="`${
						active
							? 'bg-slate-100 dark:bg-slate-700 dark:bg-opacity-70 text-slate-900 dark:text-slate-300'
							: 'text-slate-600 dark:text-slate-300'
					} `"
					class="inline-flex items-center space-x-2 rtl:space-x-reverse w-full px-4 py-2 first:rounded-t last:rounded-b font-normal cursor-pointer"
					@click="item.link()"
				>
					<div class="flex-none text-lg">
						<DKIcon :icon="item.icon" />
					</div>
					<div class="flex-1 text-sm">
						{{ item.label }}
					</div>
				</div>
			</MenuItem>
		</template>
	</Dropdown>
</template>
<script setup lang="ts">
import { MenuItem } from '@headlessui/vue';
import profileImg from '@/assets/images/all-img/user.jpg';
import Dropdown from '@/components/Dropdown/index.vue';
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import { useAuthStore } from '@/store/auth.store';

const authStore = useAuthStore();

const router = useRouter();
const ProfileMenu = ref<ItemDropMenuProfile[]>([
	{
		icon: 'heroicons-outline:login',
		label: 'Logout',
		link: () => {
			authStore.logout();
			router.go(0);
		},
	},
]);

interface ItemDropMenuProfile {
	label: string;
	icon: string;
	link: Function;
}
</script>
