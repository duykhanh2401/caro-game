import type { App } from 'vue';
import Icon from '@/components/Icon/index.vue';
import Modal from '@/components/Modal/index.vue';
import Button from '@/components/Button/index.vue';
import TextInput from '@/components/TextInput/index.vue';
import Card from '@/components/Card/index.vue';
export default function InitGlobalComponent(app: App<Element>) {
	app.component('DKIcon', Icon);
	app.component('DKModal', Modal);
	app.component('DKButton', Button);
	app.component('DKTextInput', TextInput);
	app.component('DKCard', Card);
}
