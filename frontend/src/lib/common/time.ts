export const formatUnixTime = (time: number, timeZone?: string) =>
	new Intl.DateTimeFormat('en-US', {
		timeZone,
		hour: 'numeric',
		hour12: true
	})
		.format(time * 1000)
		.replace(' ', '')
		.toLowerCase();

export const formatWeekdayDate = (time: number, timeZone?: string) =>
	new Intl.DateTimeFormat(undefined, {
		timeZone,
		weekday: 'long',
		month: 'long',
		day: 'numeric',
		year: 'numeric'
	}).format(time * 1000);

export const formatHourMinute = (time: number, timeZone?: string) =>
	new Intl.DateTimeFormat(undefined, {
		timeZone,
		hour: 'numeric',
		minute: '2-digit'
	}).format(time * 1000);

export const formatDateKey = (time: number, timeZone?: string) => {
	const parts = new Intl.DateTimeFormat('en-US', {
		timeZone,
		year: 'numeric',
		month: '2-digit',
		day: '2-digit'
	}).formatToParts(time * 1000);

	const values = Object.fromEntries(parts.map((part) => [part.type, part.value]));
	return `${values.year}-${values.month}-${values.day}`;
};
