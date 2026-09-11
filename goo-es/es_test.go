package gooes

import (
	"testing"

	goo_utils "github.com/liqiongtao/googo.io/goo-utils"
)

func TestClient(t *testing.T) {
	Init(Config{
		Addresses: []string{""},
		User:      "",
		Password:  "",
	})

	date := "2024-12-08"

	index := []string{"smartcard-http-sca-*"}

	query := goo_utils.M{
		"sort": []goo_utils.M{
			{"log_datetime": "asc"},
		},
		"query": goo_utils.M{
			"bool": goo_utils.M{
				"filter": []goo_utils.M{
					{"match_phrase": goo_utils.M{"log_tags": "/sca/device/reportinfo"}},
					{
						"range": goo_utils.M{
							"log_datetime": goo_utils.M{
								"gte": date + " 00:00:00",
								"lte": date + " 23:59:59",
							},
						},
					},
				},
			},
		},
	}

	Client().PageSearch(index, query.Json(), func(p goo_utils.Params) error {
		//fmt.Println(p)
		return nil
	})
}
