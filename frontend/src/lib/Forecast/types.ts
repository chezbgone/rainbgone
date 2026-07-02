export interface Geocode {
	formatted_address: string;
	geometry: {
		location: {
			lat: number;
			lng: number;
		};
		viewport: {
			northeast: {
				lat: number;
				lng: number;
			};
			southwest: {
				lat: number;
				lng: number;
			};
		};
	};
}

export interface Forecast {
	alerts: {
		title: string;
		regions: string[];
		severity: string;
		time: number;
		expires: number;
		description: string;
		uri: string;
	}[];
	latitude: number;
	longitude: number;
	timezone: string;
	offset: number;
	elevation: number;
	currently: {
		time: number;
		summary: string;
		icon: string;
		nearestStormDistance: number;
		nearestStormBearing: number;
		precipIntensity: number;
		precipProbability: number;
		precipIntensityError: number;
		precipType: 'rain' | 'snow' | 'sleet' | 'none';
		temperature: number;
		apparentTemperature: number;
		dewPoint: number;
		humidity: number;
		pressure: number;
		windSpeed: number;
		windGust: number;
		windBearing: number;
		cloudCover: number;
		uvIndex: number;
		visibility: number;
		ozone: number;
		smoke?: number;
		fireIndex?: number;
		feelsLike?: number;
		currentDayIce?: number;
		currentDayLiquid?: number;
		currentDaySnow?: number;
	};
	minutely: {
		summary: string;
		icon: string;
		data: {
			time: number;
			precipIntensity: number;
			precipProbability: number;
			precipIntensityError: number;
			precipType: 'rain' | 'snow' | 'sleet' | 'none';
		}[];
	};
	hourly: {
		summary: string;
		icon: string;
		data: {
			time: number;
			summary: string;
			icon: string;
			precipIntensity: number;
			precipProbability: number;
			precipIntensityError: number;
			precipAccumulation: number;
			precipType: 'rain' | 'snow' | 'sleet' | 'none';
			temperature: number;
			apparentTemperature: number;
			dewPoint: number;
			humidity: number;
			pressure: number;
			windSpeed: number;
			windGust: number;
			windBearing: number;
			cloudCover: number;
			uvIndex: number;
			visibility: number;
			ozone: number;
			smoke?: number;
			liquidAccumulation?: number;
			snowAccumulation?: number;
			iceAccumulation?: number;
			nearestStormDistance?: number;
			nearestStormBearing?: number;
			fireIndex?: number;
			feelsLike?: number;
		}[];
	};
	// Flat hourly series anchored at today's local midnight (today's elapsed hours are
	// backfilled by the backend via the Time Machine API). Use this for per-day stripes
	// via contiguous 24-hour slices, instead of `hourly`, whose series starts at the
	// current hour.
	hourlyFromMidnight: Forecast['hourly']['data'];
	daily: {
		summary: string;
		icon: string;
		data: {
			time: number;
			summary: string;
			icon: string;
			dawnTime?: number;
			sunriseTime: number;
			sunsetTime: number;
			duskTime?: number;
			moonPhase: number;
			precipIntensity: number;
			precipIntensityMax: number;
			precipIntensityMaxTime: number;
			precipProbability: number;
			precipAccumulation: number;
			precipType: 'rain' | 'snow' | 'sleet' | 'none';
			temperatureHigh: number;
			temperatureHighTime: number;
			temperatureLow: number;
			temperatureLowTime: number;
			apparentTemperatureHigh: number;
			apparentTemperatureHighTime: number;
			apparentTemperatureLow: number;
			apparentTemperatureLowTime: number;
			dewPoint: number;
			humidity: number;
			pressure: number;
			windSpeed: number;
			windGust: number;
			windGustTime: number;
			windBearing: number;
			cloudCover: number;
			uvIndex: number;
			uvIndexTime: number;
			visibility: number;
			temperatureMin: number;
			temperatureMinTime: number;
			temperatureMax: number;
			temperatureMaxTime: number;
			apparentTemperatureMin: number;
			apparentTemperatureMinTime: number;
			apparentTemperatureMax: number;
			apparentTemperatureMaxTime: number;
			smokeMax?: number;
			smokeMaxTime?: number;
			liquidAccumulation?: number;
			snowAccumulation?: number;
			iceAccumulation?: number;
			fireIndexMax?: number;
			fireIndexMaxTime?: number;
		}[];
	};
	flags: {
		sources: string[];
		sourceTimes: {
			'hrrr_0-18': string;
			hrrr_subh: string;
			nbm: string;
			nbm_fire: string;
			'hrrr_18-48': string;
			gfs: string;
			gefs: string;
		};
		'nearest-station': number;
		units: string;
		version: string;
		processTime?: number;
	};
}

export type DailyDatum = Forecast['daily']['data'][number];
