import type { DailyDatum, Forecast } from '$lib/Forecast/types';

type HourlyDatum = Forecast['hourly']['data'][number];

export function hour(time: number, overrides: Partial<HourlyDatum> = {}): HourlyDatum {
	return {
		time,
		summary: 'Clear',
		icon: 'clear-day',
		precipIntensity: 0,
		precipProbability: 0,
		precipIntensityError: 0,
		precipAccumulation: 0,
		precipType: 'none',
		temperature: 50,
		apparentTemperature: 48,
		dewPoint: 40,
		humidity: 0.5,
		pressure: 1010,
		windSpeed: 5,
		windGust: 8,
		windBearing: 90,
		cloudCover: 0,
		uvIndex: 1,
		visibility: 10,
		ozone: 300,
		...overrides
	};
}

export function hoursBetween(start: string, end: string): HourlyDatum[] {
	const first = Date.parse(start) / 1000;
	const last = Date.parse(end) / 1000;
	return Array.from({ length: (last - first) / 3600 + 1 }, (_, i) => hour(first + i * 3600));
}

export function day(time: number, overrides: Partial<DailyDatum> = {}): DailyDatum {
	return {
		time,
		summary: 'Clear all day',
		icon: 'clear-day',
		sunriseTime: time + 7 * 3600,
		sunsetTime: time + 17 * 3600,
		moonPhase: 0.5,
		precipIntensity: 0,
		precipIntensityMax: 0,
		precipIntensityMaxTime: time,
		precipProbability: 0,
		precipAccumulation: 0,
		precipType: 'none',
		temperatureHigh: 60,
		temperatureHighTime: time + 14 * 3600,
		temperatureLow: 40,
		temperatureLowTime: time + 6 * 3600,
		apparentTemperatureHigh: 58,
		apparentTemperatureHighTime: time + 14 * 3600,
		apparentTemperatureLow: 38,
		apparentTemperatureLowTime: time + 6 * 3600,
		dewPoint: 40,
		humidity: 0.5,
		pressure: 1010,
		windSpeed: 5,
		windGust: 8,
		windGustTime: time,
		windBearing: 90,
		cloudCover: 0,
		uvIndex: 1,
		uvIndexTime: time + 12 * 3600,
		visibility: 10,
		temperatureMin: 40,
		temperatureMinTime: time + 6 * 3600,
		temperatureMax: 60,
		temperatureMaxTime: time + 14 * 3600,
		apparentTemperatureMin: 38,
		apparentTemperatureMinTime: time + 6 * 3600,
		apparentTemperatureMax: 58,
		apparentTemperatureMaxTime: time + 14 * 3600,
		...overrides
	};
}

export function currently(
	time: number,
	overrides: Partial<Forecast['currently']> = {}
): Forecast['currently'] {
	return { ...hour(time), nearestStormDistance: 10, nearestStormBearing: 90, ...overrides };
}
