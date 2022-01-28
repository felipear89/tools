package display

import (
	"get-service-version/entity"
)

type Urls struct {
	Live                      string
	Preprodlive               string
	ProdBlue, ProdGreen       string
	PreprodBlue, PreprodGreen string
}

type servicesMap map[string]Urls

var Services servicesMap

func init() {
	Services = servicesMap{
		entity.BootstrapV4: Urls{
			Live:         "URL",
			Preprodlive:  "URL",
			ProdBlue:     "URL",
			ProdGreen:    "URL",
			PreprodBlue:  "URL",
			PreprodGreen: "URL",
		},
		entity.BootstrapV3: Urls{
			Live:         "URL",
			Preprodlive:  "URL",
			ProdBlue:     "URL",
			ProdGreen:    "URL",
			PreprodBlue:  "URL",
			PreprodGreen: "URL",
		},
	}
}
