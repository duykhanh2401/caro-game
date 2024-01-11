<template>
	<ul>
		<li
			v-for="(item, i) in items"
			:key="i"
			:class="`
			${item.children && item.children.length > 0 ? 'item-has-children' : ''}
			${activeSubmenu === i ? 'open' : ''}
			${router.currentRoute.value.path === item.path ? 'menu-item-active' : ''}
      `"
			class="single-sidebar-menu"
		>
			<!-- ?? single menu with no childrenred !!  -->

			<router-link
				:to="`${item.path}`"
				class="menu-link"
				v-if="item.children?.length == 0 && !item.isHidden"
			>
				<span class="menu-icon flex-grow-0" v-if="item.icon">
					<DKIcon :icon="item.icon"
				/></span>
				<div class="text-box flex-grow" v-if="item.name">
					{{ item.name }}
				</div>
			</router-link>

			<!-- !!sub menu parent li !! -->
			<div
				class="menu-link chekc"
				v-else-if="!item.isHidden"
				:class="activeSubmenu === i ? 'parent_active not-collapsed' : 'collapsed'"
				@click="toggleSubmenu(i)"
			>
				<div class="flex-1 flex items-start">
					<span class="menu-icon" v-show="item.icon">
						<DKIcon :icon="item.icon"
					/></span>
					<div class="text-box" v-if="item.name">{{ item.name }}</div>
				</div>
				<div class="flex-0">
					<div
						class="menu-arrow transform transition-all duration-300"
						:class="
							activeSubmenu === i
								? ' ltr:rotate-90 rtl:rotate-90'
								: 'rtl:rotate-180'
						"
					>
						<DKIcon icon="heroicons-outline:chevron-right" />
					</div>
				</div>
			</div>
			<Transition
				v-if="!item.isHidden"
				enter-active-class="submenu_enter-active"
				leave-active-class="submenu_leave-active"
				@before-enter="beforeEnter"
				@enter="enter"
				@after-enter="afterEnter"
				@before-leave="beforeLeave"
				@leave="leave"
				@after-leave="afterLeave"
			>
				<!-- !! SubMenu !! -->
				<ul class="sub-menu" v-if="i === activeSubmenu">
					<li
						v-for="(ci, index) in item.children"
						:key="index"
						class="block pl-4 pr-1 mb-2 first:mt-2"
					>
						<router-link
							v-if="!ci.isHidden"
							:to="`${item.path}/${ci.path}`"
							v-slot="{ isActive }"
						>
							<span
								class="text-sm flex space-x-3 items-center transition-all duration-150"
								:class="
									isActive
										? ' text-slate-900 dark:text-white font-medium'
										: 'text-slate-600 dark:text-slate-300'
								"
							>
								<span
									class="ml-3 text-lg flex-grow-0 m-0"
									v-if="item.icon"
								>
									<DKIcon :icon="ci.icon"
								/></span>
								<span class="flex-1">
									{{ ci.name }}
								</span>
							</span>
						</router-link>
					</li>
				</ul>
			</Transition>
		</li>
	</ul>
</template>
<script setup lang="ts">
import { type Menu } from '@/types';
import { ref, onMounted, onUpdated } from 'vue';
import { useRouter } from 'vue-router';
const router = useRouter();
const activeSubmenu = ref<Number | null>(null);
const props = defineProps({
	name: {
		type: String,
		default: '',
	},
	icon: {
		type: String,
		default: '',
	},
	link: {
		type: String,
		default: '',
	},
	items: { type: Array<Menu>, required: true },
});

function beforeEnter(element) {
	requestAnimationFrame(() => {
		if (!element.style.height) {
			element.style.height = '0px';
		}

		element.style.display = null;
	});
}
function enter(element) {
	requestAnimationFrame(() => {
		requestAnimationFrame(() => {
			element.style.height = `${element.scrollHeight}px`;
		});
	});
}
function afterEnter(element) {
	element.style.height = null;
}
function beforeLeave(element) {
	requestAnimationFrame(() => {
		if (!element.style.height) {
			element.style.height = `${element.offsetHeight}px`;
		}
	});
}
function leave(element) {
	requestAnimationFrame(() => {
		requestAnimationFrame(() => {
			element.style.height = '0px';
		});
	});
}
function afterLeave(element) {
	element.style.height = null;
}
function toggleSubmenu(index) {
	console.log(index, activeSubmenu.value);

	if (activeSubmenu.value === index) {
		activeSubmenu.value = null;
	} else {
		activeSubmenu.value = index;
	}
}

onMounted(() => {
	props.items.map((item) => {
		if (item.path === router.currentRoute.value.name) {
			activeSubmenu.value = null;
		}
	});

	props.items.map((item, i) => {
		item.children?.map((ci) => {
			if (ci.name === router.currentRoute.value.name) {
				activeSubmenu.value = i;
			}
		});
	});
});

onUpdated(() => {});
// watch: {
// 	$route() {
// 		if (this.$store.themeSettingsStore.mobielSidebar) {
// 			this.$store.themeSettingsStore.mobielSidebar = false;
// 		}

// 		this.items.map((item) => {
// 			if (item.path === this.$route.name) {
// 				this.activeSubmenu = null;
// 			}
// 		});
// 	},
// },

// created() {
// 	this.items.map((item, i) => {
// 		item.children?.map((ci) => {
// 			console.log(ci.path, router.currentRoute.value.name);

// 			if (ci.name === router.currentRoute.value.name) {
// 				this.activeSubmenu = i;
// 			}
// 		});
// 	});
// },
// // update if route chnage then activesubmenu null

// updated() {},
</script>
<style lang="scss">
.submenu_enter-active,
.submenu_leave-active {
	overflow: hidden;
	transition: all 0.34s linear;
}

.not-collapsed .has-icon {
	transition: all 0.34s linear;
}
.not-collapsed .has-icon {
	@apply transform rotate-180;
}

// single sidebar menu css
.single-sidebar-menu {
	@apply relative;
	.menulabel {
		@apply text-slate-800 dark:text-slate-300 text-xs font-semibold uppercase mb-4 mt-4;
	}
	> .menu-link {
		@apply flex text-slate-600 font-medium dark:text-slate-300 text-sm capitalize px-[10px] py-3 rounded-[4px] cursor-pointer;
	}
	.menu-icon {
		@apply icon-box inline-flex items-center text-slate-600 dark:text-slate-300 text-lg ltr:mr-3 rtl:ml-3;
	}
}
// menu item has chilren
.item-has-children {
	.menu-arrow {
		@apply h-5 w-5 text-base text-slate-300 bg-slate-100 dark:bg-[#334155] dark:text-slate-300 rounded-full flex justify-center items-center;
	}
}

// close sidebar css
.close_sidebar .menulabel {
	@apply hidden;
}

.close_sidebar:not(.sidebar-hovered) {
	.menu-arrow {
		@apply hidden;
	}
	.single-sidebar-menu {
		.text-box {
			@apply absolute  left-full ml-5 w-[180px] top-0 px-4 py-3 bg-white shadow-dropdown rounded-[4px] dark:bg-slate-800 z-[999] invisible opacity-0 transition-all duration-150;
		}
		&:hover {
			.text-box {
				@apply visible opacity-100;
			}
		}
	}
	.item-has-children {
		.text-box {
			@apply hidden;
		}

		> ul {
			@apply ml-4 absolute left-full top-0 w-[230px] bg-white shadow-dropdown rounded-[4px] dark:bg-slate-800 z-[999] px-4 pt-3 transition-all duration-150 invisible opacity-0;
			display: block !important;
		}
		&:hover {
			> ul {
				@apply visible opacity-100;
			}
		}
	}
}
.menu-badge {
	@apply py-1 px-2 text-xs capitalize font-semibold rounded-[.358rem] whitespace-nowrap align-baseline inline-flex bg-slate-900 text-slate-50 dark:bg-slate-700 dark:text-slate-300;
}
// active menu
.item-has-children {
	.parent_active {
		@apply bg-secondary-500 bg-opacity-20;
		.icon-box,
		.menu-icon,
		.text-box {
			@apply text-slate-700 dark:text-slate-200;
		}
		.menu-arrow {
			@apply bg-secondary-500 text-slate-600 text-opacity-70 bg-opacity-30 dark:text-white;
		}
	}
}
.menu-item-active {
	.menu-link {
		@apply bg-slate-800 dark:bg-slate-700;
		.icon-box,
		.menu-icon,
		.text-box {
			@apply text-white dark:text-slate-300;
		}
	}
	.menu-badge {
		@apply bg-slate-100  text-slate-900;
	}
}
</style>
