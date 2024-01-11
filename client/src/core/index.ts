import type { App } from 'vue';
import { InitToast } from './toast';
import { InitStyle } from './style';
import './nprocess';
export function InitCore(app: App<Element>) {
	InitToast(app);
	InitStyle(app);
}
