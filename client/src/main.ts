import { createApp } from 'vue';
import App from './App.vue';

import router from './router';
import InitGlobalComponent from './global';
import { InitCore } from './core';
import { InitStore } from './store';
const app = createApp(App);

InitStore(app);
InitCore(app);

app.use(router);
app.mount('#app');
InitGlobalComponent(app);
