import Vue, { type App } from 'vue';
import Toast, { POSITION, type PluginOptions } from 'vue-toastification';
import 'vue-toastification/dist/index.css';

const options: PluginOptions = {
	timeout: 2000,
	closeOnClick: true,
	pauseOnFocusLoss: false,
	pauseOnHover: true,
	draggable: true,
	draggablePercent: 0.6,
	showCloseButtonOnHover: false,
	hideProgressBar: false,
	closeButton: 'button',
	icon: true,
	rtl: false,
	position: POSITION.TOP_CENTER,
};

export function InitToast(app: App<Element>) {
	app.use(Toast, options);
}
