package gooes

import (
	"testing"

	gooutils "github.com/liqiongtao/googo.io/goo-util"
)

func TestClient(t *testing.T) {
	Init(Config{
		Addresses: []string{""},
		User:      "",
		Password:  "",
	})

	date := "2024-12-08"

	index := []string{"smartcard-http-sca-*"}

	query := gooutils.M{
		"sort": []gooutils.M{
			{"log_datetime": "asc"},
		},
		"query": gooutils.M{
			"bool": gooutils.M{
				"filter": []gooutils.M{
					{"match_phrase": gooutils.M{"log_tags": "/sca/device/reportinfo"}},
					{
						"range": gooutils.M{
							"log_datetime": gooutils.M{
								"gte": date + " 00:00:00",
								"lte": date + " 23:59:59",
							},
						},
					},
				},
			},
		},
	}

	Client().PageSearch(index, query.Json(), func(p gooutils.Params) error {
		//fmt.Println(p)
		return nil
	})
}
