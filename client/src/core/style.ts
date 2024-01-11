import 'element-plus/theme-chalk/dark/css-vars.css';
import 'element-plus/dist/index.css';
import '@/assets/scss/tailwind.scss';
import '@/assets/scss/auth.scss';
import 'remixicon/fonts/remixicon.css';
import { useThemeSetting } from '@/store/themeSetting.store';
import { storeToRefs } from 'pinia';
import '@/assets/vue-multiselect.scss';
import ElementPlus from 'element-plus';
import type { App } from 'vue';
export function InitStyle(app: App<Element>) {
	app.use(ElementPlus);
	const themeSettingStore = useThemeSetting();
	const { theme } = storeToRefs(themeSettingStore);
	document.documentElement.classList.add(theme.value);
	document.body.classList.add(theme.value);
}
