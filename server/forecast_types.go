package server

type ForecastResponse struct {
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	Timezone  string    `json:"timezone"`
	Offset    float64   `json:"offset"`
	Elevation int       `json:"elevation"`
	Currently Currently `json:"currently"`
	Minutely  Minutely  `json:"minutely"`
	Hourly    Hourly    `json:"hourly"`
	Daily     Daily     `json:"daily"`
	Alerts    []Alerts  `json:"alerts"`
	Flags     Flags     `json:"flags"`
}
type Currently struct {
	Time                 int     `json:"time"`
	Summary              string  `json:"summary"`
	Icon                 string  `json:"icon"`
	NearestStormDistance float64 `json:"nearestStormDistance"`
	NearestStormBearing  float64 `json:"nearestStormBearing"`
	PrecipIntensity      float64 `json:"precipIntensity"`
	PrecipProbability    float64 `json:"precipProbability"`
	PrecipIntensityError float64 `json:"precipIntensityError"`
	PrecipType           string  `json:"precipType"`
	Temperature          float64 `json:"temperature"`
	ApparentTemperature  float64 `json:"apparentTemperature"`
	DewPoint             float64 `json:"dewPoint"`
	Humidity             float64 `json:"humidity"`
	Pressure             float64 `json:"pressure"`
	WindSpeed            float64 `json:"windSpeed"`
	WindGust             float64 `json:"windGust"`
	WindBearing          float64 `json:"windBearing"`
	CloudCover           float64 `json:"cloudCover"`
	UvIndex              float64 `json:"uvIndex"`
	Visibility           float64 `json:"visibility"`
	Ozone                float64 `json:"ozone"`
	Smoke                float64 `json:"smoke"`
	FireIndex            float64 `json:"fireIndex"`
	FeelsLike            float64 `json:"feelsLike"`
	CurrentDayIce        float64 `json:"currentDayIce"`
	CurrentDayLiquid     float64 `json:"currentDayLiquid"`
	CurrentDaySnow       float64 `json:"currentDaySnow"`
}
type MinutelyData struct {
	Time                 int     `json:"time"`
	PrecipIntensity      float64 `json:"precipIntensity"`
	PrecipProbability    float64 `json:"precipProbability"`
	PrecipIntensityError float64 `json:"precipIntensityError"`
	PrecipType           string  `json:"precipType"`
}
type Minutely struct {
	Summary string         `json:"summary"`
	Icon    string         `json:"icon"`
	Data    []MinutelyData `json:"data"`
}
type HourlyData struct {
	Time                 int     `json:"time"`
	Icon                 string  `json:"icon"`
	Summary              string  `json:"summary"`
	PrecipIntensity      float64 `json:"precipIntensity"`
	PrecipProbability    float64 `json:"precipProbability"`
	PrecipIntensityError float64 `json:"precipIntensityError"`
	PrecipAccumulation   float64 `json:"precipAccumulation"`
	PrecipType           string  `json:"precipType"`
	Temperature          float64 `json:"temperature"`
	ApparentTemperature  float64 `json:"apparentTemperature"`
	DewPoint             float64 `json:"dewPoint"`
	Humidity             float64 `json:"humidity"`
	Pressure             float64 `json:"pressure"`
	WindSpeed            float64 `json:"windSpeed"`
	WindGust             float64 `json:"windGust"`
	WindBearing          float64 `json:"windBearing"`
	CloudCover           float64 `json:"cloudCover"`
	UvIndex              float64 `json:"uvIndex"`
	Visibility           float64 `json:"visibility"`
	Ozone                float64 `json:"ozone"`
	Smoke                float64 `json:"smoke"`
	LiquidAccumulation   float64 `json:"liquidAccumulation"`
	SnowAccumulation     float64 `json:"snowAccumulation"`
	IceAccumulation      float64 `json:"iceAccumulation"`
	NearestStormDistance float64 `json:"nearestStormDistance"`
	NearestStormBearing  float64 `json:"nearestStormBearing"`
	FireIndex            float64 `json:"fireIndex"`
	FeelsLike            float64 `json:"feelsLike"`
}
type Hourly struct {
	Summary string       `json:"summary"`
	Icon    string       `json:"icon"`
	Data    []HourlyData `json:"data"`
}
type DailyData struct {
	Time                        int     `json:"time"`
	Icon                        string  `json:"icon"`
	Summary                     string  `json:"summary"`
	DawnTime                    int     `json:"dawnTime"`
	SunriseTime                 int     `json:"sunriseTime"`
	SunsetTime                  int     `json:"sunsetTime"`
	DuskTime                    int     `json:"duskTime"`
	MoonPhase                   float64 `json:"moonPhase"`
	PrecipIntensity             float64 `json:"precipIntensity"`
	PrecipIntensityMax          float64 `json:"precipIntensityMax"`
	PrecipIntensityMaxTime      int     `json:"precipIntensityMaxTime"`
	PrecipProbability           float64 `json:"precipProbability"`
	PrecipAccumulation          float64 `json:"precipAccumulation"`
	PrecipType                  string  `json:"precipType"`
	TemperatureHigh             float64 `json:"temperatureHigh"`
	TemperatureHighTime         int     `json:"temperatureHighTime"`
	TemperatureLow              float64 `json:"temperatureLow"`
	TemperatureLowTime          int     `json:"temperatureLowTime"`
	ApparentTemperatureHigh     float64 `json:"apparentTemperatureHigh"`
	ApparentTemperatureHighTime int     `json:"apparentTemperatureHighTime"`
	ApparentTemperatureLow      float64 `json:"apparentTemperatureLow"`
	ApparentTemperatureLowTime  int     `json:"apparentTemperatureLowTime"`
	DewPoint                    float64 `json:"dewPoint"`
	Humidity                    float64 `json:"humidity"`
	Pressure                    float64 `json:"pressure"`
	WindSpeed                   float64 `json:"windSpeed"`
	WindGust                    float64 `json:"windGust"`
	WindGustTime                int     `json:"windGustTime"`
	WindBearing                 float64 `json:"windBearing"`
	CloudCover                  float64 `json:"cloudCover"`
	UvIndex                     float64 `json:"uvIndex"`
	UvIndexTime                 int     `json:"uvIndexTime"`
	Visibility                  float64 `json:"visibility"`
	TemperatureMin              float64 `json:"temperatureMin"`
	TemperatureMinTime          int     `json:"temperatureMinTime"`
	TemperatureMax              float64 `json:"temperatureMax"`
	TemperatureMaxTime          int     `json:"temperatureMaxTime"`
	ApparentTemperatureMin      float64 `json:"apparentTemperatureMin"`
	ApparentTemperatureMinTime  int     `json:"apparentTemperatureMinTime"`
	ApparentTemperatureMax      float64 `json:"apparentTemperatureMax"`
	ApparentTemperatureMaxTime  int     `json:"apparentTemperatureMaxTime"`
	SmokeMax                    float64 `json:"smokeMax"`
	SmokeMaxTime                int     `json:"smokeMaxTime"`
	LiquidAccumulation          float64 `json:"liquidAccumulation"`
	SnowAccumulation            float64 `json:"snowAccumulation"`
	IceAccumulation             float64 `json:"iceAccumulation"`
	FireIndexMax                float64 `json:"fireIndexMax"`
	FireIndexMaxTime            int     `json:"fireIndexMaxTime"`
}
type Daily struct {
	Summary string      `json:"summary"`
	Icon    string      `json:"icon"`
	Data    []DailyData `json:"data"`
}
type Alerts struct {
	Title       string   `json:"title"`
	Regions     []string `json:"regions"`
	Severity    string   `json:"severity"`
	Time        int      `json:"time"`
	Expires     int      `json:"expires"`
	Description string   `json:"description"`
	URI         string   `json:"uri"`
}
type SourceTimes struct {
	Hrrr018  string `json:"hrrr_0-18"`
	HrrrSubh string `json:"hrrr_subh"`
	Nbm      string `json:"nbm"`
	NbmFire  string `json:"nbm_fire"`
	Hrrr1848 string `json:"hrrr_18-48"`
	Gfs      string `json:"gfs"`
	Gefs     string `json:"gefs"`
}
type Hrrr struct {
	X    int     `json:"x"`
	Y    int     `json:"y"`
	Lat  float64 `json:"lat"`
	Long float64 `json:"long"`
}
type Nbm struct {
	X    int     `json:"x"`
	Y    int     `json:"y"`
	Lat  float64 `json:"lat"`
	Long float64 `json:"long"`
}
type Gfs struct {
	X    int     `json:"x"`
	Y    int     `json:"y"`
	Lat  float64 `json:"lat"`
	Long float64 `json:"long"`
}
type Etopo struct {
	X    int     `json:"x"`
	Y    int     `json:"y"`
	Lat  float64 `json:"lat"`
	Long float64 `json:"long"`
}
type SourceIDX struct {
	Hrrr  Hrrr  `json:"hrrr"`
	Nbm   Nbm   `json:"nbm"`
	Gfs   Gfs   `json:"gfs"`
	Etopo Etopo `json:"etopo"`
}
type Flags struct {
	Sources        []string    `json:"sources"`
	SourceTimes    SourceTimes `json:"sourceTimes"`
	NearestStation int         `json:"nearest-station"`
	Units          string      `json:"units"`
	Version        string      `json:"version"`
	SourceIDX      SourceIDX   `json:"sourceIDX"`
	ProcessTime    int         `json:"processTime"`
}
