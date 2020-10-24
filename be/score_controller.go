package main

import (
	"encoding/json"
	"log"
	controller "net/http"
	"strconv"
)

func highScores(res controller.ResponseWriter, req *controller.Request) {
	res.Header().Set("Access-Control-Allow-Origin", "*")
	res.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, DELETE, PUT")
	res.Header().Set("Access-Control-Allow-Headers", "*")

	if req.Method == controller.MethodGet {
		res.Header().Set("Content-Type", "application/json")

		offsetParam, offsetExist := req.URL.Query()["offset"]
		limitParam, limitExist := req.URL.Query()["limit"]

		var offset int
		var limit int

		if !offsetExist || len(offsetParam[0]) == 0 {
			offset = 0
		} else {
			i, e := strconv.Atoi(offsetParam[0])

			if e != nil {
				offset = 0
			} else {
				offset = i
			}
		}

		if !limitExist || len(limitParam[0]) == 0 {
			limit = 10
		} else {
			i, e := strconv.Atoi(limitParam[0])

			if e != nil {
				limit = 10
			} else {
				limit = i
			}
		}

		selRows, errSelRows := dbConn.Query("select username, timestampdiff(second, created_at, finished_at) seconds, who_win from games where who_win = 1 order by seconds limit ? offset ?", limit, offset)

		if errSelRows != nil {
			log.Println("DB query error!")
			log.Println(errSelRows)

			res.WriteHeader(controller.StatusInternalServerError)
		} else {
			var highScoresArr []Score

			for selRows.Next() {
				var score Score

				scanErr := selRows.Scan(&score.Username, &score.FinishedIn, &score.WhoWin)

				if scanErr != nil {
					log.Println("Query result scan error!")
					log.Println(scanErr)
				} else {
					highScoresArr = append(highScoresArr, score)
				}
			}

			highScoresArrJSON, highScoresArrJSONErr := json.Marshal(highScoresArr)

			if highScoresArrJSONErr != nil {
				log.Println("Converting high scores arr to json error!")
				log.Println(highScoresArrJSONErr)

				res.WriteHeader(controller.StatusInternalServerError)
			} else {
				res.WriteHeader(controller.StatusOK)

				res.Write([]byte(string(highScoresArrJSON)))
			}
		}
	} else {
		res.WriteHeader(controller.StatusMethodNotAllowed)
	}
}
