import { ref, onMounted, onUnmounted, reactive } from 'vue';

export function useWindow() {
	const width = ref(0);

	function handleResize() {
		width.value = window.innerWidth;
	}

	onMounted(() => {
		window.addEventListener('resize', handleResize);
		handleResize();
	});

	onUnmounted(() => {
		window.removeEventListener('resize', handleResize);
	});
	return { width };
}
