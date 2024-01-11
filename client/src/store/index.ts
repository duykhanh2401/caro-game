import piniaPluginPersistedstate from 'pinia-plugin-persistedstate';
import { createPinia, storeToRefs } from 'pinia';
import type { App } from 'vue';
const pinia = createPinia();
pinia.use(piniaPluginPersistedstate);

export function InitStore(app: App<Element>) {
	app.use(pinia);
}
