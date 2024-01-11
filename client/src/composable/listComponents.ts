export function useListComponents() {
	const viewModules = import.meta.glob('@/views/**/*.vue');
	const listComponents = Object.keys(viewModules);

	return { listComponents, viewModules };
}
