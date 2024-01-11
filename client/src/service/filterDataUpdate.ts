export const filterDataUpdate = (oldData, newData) => {
	const dataUpdate = {};
	for (const [key, value] of Object.entries(newData)) {
		if (JSON.stringify(oldData[key]) != JSON.stringify(value)) {
			dataUpdate[key] = value;
		}
	}

	return dataUpdate;
};
