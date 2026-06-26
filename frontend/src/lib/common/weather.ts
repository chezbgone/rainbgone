export enum Weather {
	Clear = 'clear',
	PartlyCloudy = 'partly cloudy',
	MostlyCloudy = 'mostly cloudy',
	Overcast = 'overcast',
	LightRain = 'light rain',
	Rain = 'rain',
	LightSnow = 'light snow',
	Snow = 'snow',
	LightSleet = 'light sleet',
	Sleet = 'sleet'
}

export interface WeatherInfo {
	weather: Weather;
	color: string;
	labelInfo: {
		text: string;
		requiredWidth: number;
		darkText: boolean;
	};
}

export const weatherInfo: Record<Weather, WeatherInfo> = {
	[Weather.Clear]: {
		weather: Weather.Clear,
		color: '#eeeef5',
		labelInfo: {
			text: 'Clear',
			requiredWidth: 2,
			darkText: true
		}
	},
	[Weather.PartlyCloudy]: {
		weather: Weather.PartlyCloudy,
		color: '#d5dae2',
		labelInfo: {
			text: 'Partly Cloudy',
			requiredWidth: 3,
			darkText: true
		}
	},
	[Weather.MostlyCloudy]: {
		weather: Weather.MostlyCloudy,
		color: '#b6bfcb',
		labelInfo: {
			text: 'Mostly Cloudy',
			requiredWidth: 3,
			darkText: true
		}
	},
	[Weather.Overcast]: {
		weather: Weather.Overcast,
		color: '#878f9a',
		labelInfo: {
			text: 'Overcast',
			requiredWidth: 2,
			darkText: false
		}
	},
	[Weather.LightRain]: {
		weather: Weather.LightRain,
		color: '#80a5d6',
		labelInfo: {
			text: 'Light Rain',
			requiredWidth: 3,
			darkText: false
		}
	},
	[Weather.Rain]: {
		weather: Weather.Rain,
		color: '#4a80c7',
		labelInfo: {
			text: 'Rain',
			requiredWidth: 2,
			darkText: false
		}
	},
	[Weather.LightSnow]: {
		weather: Weather.LightSnow,
		color: '#aba4db',
		labelInfo: {
			text: 'Light Snow',
			requiredWidth: 3,
			darkText: false
		}
	},
	[Weather.Snow]: {
		weather: Weather.Snow,
		color: '#8c82ce',
		labelInfo: {
			text: 'Snow',
			requiredWidth: 2,
			darkText: false
		}
	},
	[Weather.LightSleet]: {
		weather: Weather.LightSleet,
		color: '#d696a3',
		labelInfo: {
			text: 'Light Sleet',
			requiredWidth: 3,
			darkText: false
		}
	},
	[Weather.Sleet]: {
		weather: Weather.Sleet,
		color: '#cc7585',
		labelInfo: {
			text: 'Sleet',
			requiredWidth: 2,
			darkText: false
		}
	}
};

export const classifyWeather = (
	precipIntensity: number,
	precipType: 'rain' | 'snow' | 'sleet' | 'none',
	cloudCover: number
): Weather => {
	if (precipType !== 'none') {
		if (precipIntensity >= 0.05) {
			if (precipType === 'rain') return Weather.Rain;
			if (precipType === 'snow') return Weather.Snow;
			if (precipType === 'sleet') return Weather.Sleet;
		}
		if (precipIntensity >= 0.01) {
			if (precipType === 'rain') return Weather.LightRain;
			if (precipType === 'snow') return Weather.LightSnow;
			if (precipType === 'sleet') return Weather.LightSleet;
		}
	}

	if (cloudCover >= 0.9) {
		return Weather.Overcast;
	}
	if (cloudCover >= 0.5) {
		return Weather.MostlyCloudy;
	}
	if (cloudCover >= 0.25) {
		return Weather.PartlyCloudy;
	}
	return Weather.Clear;
};
