import { createRouter, createWebHistory } from 'vue-router';
import Nprogress from 'nprogress';

const router = createRouter({
	history: createWebHistory(import.meta.env.BASE_URL),
	routes: [
		{
			path: '/',
			name: 'Home',
			component: () => import('@/views/Home.vue'),
		},
	],
});

router.beforeEach(async (to, from) => {
	Nprogress.start();
});

router.afterEach(() => {
	Nprogress.done();
});

router.onError(() => {
	Nprogress.remove();
});

export default router;
