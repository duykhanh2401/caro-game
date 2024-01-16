import { defineStore } from 'pinia';
import { ref } from 'vue';

export const useThemeSetting = defineStore(
	'themeSetting',
	() => {
		{
			const isDark = ref(true);
			const theme = ref('dark');
			const sidebarCollapse = ref(false);
			const sidebarHidden = ref(false);
			const isMouseHovered = ref(false);
			const sidebarMobile = ref(false);
			function toggleDarkTheme() {
				isDark.value = !isDark.value;
				document.documentElement.classList.remove(theme.value);
				document.body.classList.remove(theme.value);
				theme.value = theme.value === 'dark' ? 'light' : 'dark';
				document.documentElement.classList.add(theme.value);
				document.body.classList.add(theme.value);
			}

			function enterMouseHover() {
				isMouseHovered.value = true;
			}

			function leaveMouseHover() {
				isMouseHovered.value = false;
			}

			function toggleMobileSidebar() {
				sidebarMobile.value = !sidebarMobile.value;
			}

			return {
				isDark,
				theme,
				sidebarCollapse,
				sidebarHidden,
				sidebarMobile,
				isMouseHovered,
				toggleDarkTheme,
				enterMouseHover,
				leaveMouseHover,
				toggleMobileSidebar,
			};
		}
	},
	{
		persist: true,
	},
);
