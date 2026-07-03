export const number = (value: number | null | undefined, digits = 0) =>
	value === null || value === undefined || !Number.isFinite(value) ? 'N/A' : value.toFixed(digits);

export const percent = (value: number | null | undefined) =>
	value === null || value === undefined || !Number.isFinite(value)
		? 'N/A'
		: `${Math.round(value * 100)}%`;
