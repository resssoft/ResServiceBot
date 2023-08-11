package weather

import "time"

const (
	checkItemsDelay = time.Second * 2
	//gismeteoUrl     = "https://api.gismeteo.net/v2/weather/forecast/aggregate/?lang=ru&days=3&latitude=40.19&longitude=44.50241"
	gismeteoUrl = "http://127.0.0.1:8080/api/w/"
)

type item struct {
	chatId int64
	hour   int
	min    int
}

type gismeteoData struct {
	Meta struct {
		Message string `json:"message,omitempty"`
		Code    string `json:"code,omitempty"`
	}
	Response []struct {
		Storm bool   `json:"storm,omitempty"`
		City  int    `json:"city,omitempty"`
		Gm    int    `json:"gm,omitempty"`
		Icon  string `json:"icon,omitempty"`

		Date struct {
			UTC   string `json:"UTC,omitempty"`
			Local string `json:"local,omitempty"`
			TZ    int    `json:"time_zone_offset,omitempty"`
			Unix  int64  `json:"unix,omitempty"`
		} `json:"date,omitempty"`
		Temperature struct {
			Comfort TemperatureData `json:"comfort,omitempty"`
			Water   TemperatureData `json:"water,omitempty"`
			Air     TemperatureData `json:"air,omitempty"`
		} `json:"temperature,omitempty"`
		Description struct {
			Full string `json:"full,omitempty"`
		} `json:"description,omitempty"`
	} `json:"response,omitempty"`
}

type TemperatureData struct {
	Max Temperatures `json:"max,omitempty"`
	Min Temperatures `json:"min,omitempty"`
}

type Temperatures struct {
	C float64 `json:"C,omitempty"`
	F float64 `json:"F,omitempty"`
}

var (
	Services          = []string{"yandex", "gismeteo"}
	gismeteoIconTypes = map[string]string{
		"d":   "День",
		"n":   "Ночь	",
		"c1":  "малооблачно",
		"c2":  "облачно",
		"c3":  "пасмурно",
		"r1":  "слабый дождь",
		"r2":  "дождь",
		"r3":  "проливной дождь",
		"s1":  "слабый снег",
		"s2":  "снег",
		"rs3": "сильный снег",
		"rs1": "слабый снег с дождем",
		"rs2": "снег с дождем",
		"s3":  "сильный снег с дождем",
		"st":  "Гроза",
	}
	gismeteoTods = map[int]string{
		0: "Ночь",
		1: "Утро",
		2: "День",
		3: "Вечер",
	}

	gismeteoIcons = map[string]string{
		"c3_r1_st":    "🌧⚡️",
		"d_c1_rs3_st": "🌧",
		"d_st":        "⚡️",
		"n_c2_rs1":    "",
		"c3_r1":       "❄️",
		"d_c1_rs3":    "",
		"d":           "",
		"n_c2_rs2_st": "⚡️",
		"c3_r2_st":    "⚡️",
		"d_c1_s1_st":  "⚡️",
		"mist":        "⚡️",
		"n_c2_rs2":    "⚡️",
		"c3_r2":       "",
		"d_c1_s1":     "",
		"n_c1_r1_st":  "⚡️",
		"n_c2_rs3_st": "⚡️",
		"c3_r3_st":    "⚡️",
		"d_c1_s2_st":  "⚡️",
		"n_c1_r1":     "",
		"n_c2_rs3":    "",
		"c3_r3":       "",
		"d_c1_s2":     "",
		"n_c1_r2_st":  "⚡️",
		"n_c2_s1_st":  "⚡️",
		"c3_rs1_st":   "⚡️",
		"d_c1_s3_st":  "⚡️",
		"n_c1_r2":     "",
		"n_c2_s1":     "",
		"c3_rs1":      "",
		"d_c1_s3":     "",
		"n_c1_r3_st":  "⚡️",
		"n_c2_s2_st":  "⚡️",
		"c3_rs2_st":   "⚡️",
		"d_c1_st":     "⚡️",
		"n_c1_r3":     "",
		"n_c2_s2":     "",
		"c3_rs2":      "",
		"d_c1":        "",
		"n_c1_rs1_st": "⚡️",
		"n_c2_s3_st":  "⚡️",
		"c3_rs3_st":   "⚡️",
		"d_c2_r1_st":  "⚡️",
		"n_c1_rs1":    "",
		"n_c2_s3":     "",
		"c3_rs3":      "",
		"d_c2_r1":     "",
		"n_c1_rs2_st": "⚡️",
		"n_c2_st":     "⚡️",
		"c3_s1_st":    "⚡️",
		"d_c2_r2_st":  "🌦⚡️",
		"n_c1_rs2":    "",
		"n_c2":        "",
		"c3_s1":       "",
		"d_c2_r2":     "🌦",
		"n_c1_rs3_st": "🌦⚡️",
		"n_st":        "⚡️",
		"c3_s2_st":    "⚡️",
		"d_c2_r3_st":  "⚡️",
		"n_c1_rs3":    "",
		"n":           "",
		"c3_s2":       "",
		"d_c2_r3":     "",
		"n_c1_s1_st":  "⚡️",
		"r1_mist":     "⚡️",
		"c3_s3_st":    "⚡️",
		"d_c2_rs1_st": "⚡️",
		"n_c1_s1":     "",
		"r1_st_mist":  "⚡️",
		"c3_s3":       "",
		"d_c2_rs1":    "",
		"n_c1_s2_st":  "⚡️",
		"r2_mist":     "⚡️",
		"c3_st":       "⚡️",
		"d_c2_rs2_st": "⚡️",
		"n_c1_s2":     "",
		"r2_st_mist":  "⚡️",
		"c3":          "",
		"d_c2_rs2":    "",
		"n_c1_s3_st":  "⚡️",
		"r3_mist":     "⚡️",
		"d_c1_r1_st":  "⚡️",
		"d_c2_rs3_st": "⚡️",
		"n_c1_s3":     "",
		"r3_st_mist":  "⚡️",
		"d_c1_r1":     "",
		"d_c2_rs3":    "",
		"n_c1_st":     "⚡️",
		"s1_mist":     "⚡️",
		"d_c1_r2_st":  "⚡️",
		"d_c2_s1_st":  "⚡️",
		"n_c1":        "",
		"s1_st_mist":  "⚡️",
		"d_c1_r2":     "",
		"d_c2_s1":     "",
		"n_c2_r1_st":  "⚡️",
		"s2_mist":     "⚡️",
		"d_c1_r3_st":  "⚡️",
		"d_c2_s2_st":  "⚡️",
		"n_c2_r1":     "",
		"s2_st_mist":  "⚡️",
		"d_c1_r3":     "",
		"d_c2_s2":     "",
		"n_c2_r2_st":  "⚡️",
		"s3_mist":     "⚡️",
		"d_c1_rs1_st": "⚡️",
		"d_c2_s3_st":  "⚡️",
		"n_c2_r2":     "",
		"s3_st_mist":  "⚡️",
		"d_c1_rs1":    "",
		"d_c2_s3":     "",
		"n_c2_r3_st":  "⚡️",
		"d_c1_rs2_st": "⚡️",
		"d_c2_st":     "⛅️⚡️",
		"n_c2_r3":     "",
		"d_c1_rs2":    "",
		"d_c2":        "",
		"n_c2_rs1_st": "⚡️",
		"":            "",
	}
)
