export const formatUnixTime = (time: number) => new Intl.DateTimeFormat("en-US", {
  hour: "numeric",
  hour12: true,
}).format(time * 1000).replace(" ", "").toLowerCase();