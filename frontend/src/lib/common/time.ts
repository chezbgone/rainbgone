export const formatUnixTime = (time: number) => new Intl.DateTimeFormat("en-US", {
  hour: "numeric",
  hour12: true,
}).format(time * 1000).replace(" ", "").toLowerCase();

export const formatDateKey = (time: number, timeZone?: string) => {
  const parts = new Intl.DateTimeFormat("en-US", {
    timeZone,
    year: "numeric",
    month: "2-digit",
    day: "2-digit",
  }).formatToParts(time * 1000);

  const values = Object.fromEntries(parts.map(part => [part.type, part.value]));
  return `${values.year}-${values.month}-${values.day}`;
};
