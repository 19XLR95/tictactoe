package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"math/rand"
	controller "net/http"
	"time"
)

func newGame(res controller.ResponseWriter, req *controller.Request) {
	res.Header().Set("Access-Control-Allow-Origin", "*")
	res.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, DELETE, PUT")
	res.Header().Set("Access-Control-Allow-Headers", "*")

	if req.Method == controller.MethodPost {
		res.Header().Set("Content-Type", "application/json")

		var newGame NewGame

		err := json.NewDecoder(req.Body).Decode(&newGame)

		if err != nil {
			log.Println("Body decoding error!")
			log.Println(err)

			res.WriteHeader(controller.StatusBadRequest)
		} else {
			if newGame.Username == "" || newGame.Difficulty == 0 {
				log.Println("Username or difficulty is empty!")

				res.WriteHeader(controller.StatusBadRequest)
			} else {
				stmt, errStmt := dbConn.Prepare("insert into games(game_key, username, difficulty) values(?, ?, ?)")

				if errStmt != nil {
					log.Println("DB prepare statement error!")
					log.Println(errStmt)

					res.WriteHeader(controller.StatusInternalServerError)
				} else {
					gameKey := convertBase64(generateRandomString(8))

					queryRes, errQueryExec := stmt.Exec(gameKey, newGame.Username, newGame.Difficulty)

					if errQueryExec != nil {
						log.Println("DB query error!")
						log.Println(errQueryExec)

						res.WriteHeader(controller.StatusInternalServerError)
					} else {
						affRows, affRowsErr := queryRes.RowsAffected()

						if affRowsErr != nil || affRows == 0 {
							log.Println("DB query error, game not saved!")
							log.Println(affRowsErr)

							res.WriteHeader(controller.StatusInternalServerError)
						} else {
							game := Game{gameKey, newGame.Username, newGame.Difficulty, initialGameState, false, -1}

							gameJSON, gameJSONErr := json.Marshal(game)

							if gameJSONErr != nil {
								log.Println("Converting game struct to json error!")
								log.Println(gameJSONErr)

								res.WriteHeader(controller.StatusInternalServerError)
							} else {
								res.WriteHeader(controller.StatusCreated)

								res.Write([]byte(string(gameJSON)))
							}
						}
					}
				}
			}
		}
	} else if req.Method == controller.MethodOptions {
		res.WriteHeader(controller.StatusOK)
	} else {
		res.WriteHeader(controller.StatusMethodNotAllowed)
	}
}

func makeMove(res controller.ResponseWriter, req *controller.Request) {
	res.Header().Set("Access-Control-Allow-Origin", "*")
	res.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, DELETE, PUT")
	res.Header().Set("Access-Control-Allow-Headers", "*")

	if req.Method == controller.MethodPost {
		res.Header().Set("Content-Type", "application/json")

		var game Game

		err := json.NewDecoder(req.Body).Decode(&game)

		if err != nil {
			log.Println("Body decoding error!")
			log.Println(err)

			res.WriteHeader(controller.StatusBadRequest)
		} else {
			if game.GameKey == "" || game.GameState == "" {
				log.Println("Game key or game state is empty!")

				res.WriteHeader(controller.StatusBadRequest)
			} else {
				var gameState [][]int

				gameStateJSONErr := json.Unmarshal([]byte(game.GameState), &gameState)

				if gameStateJSONErr != nil {
					log.Println("Invalid game state!")

					res.WriteHeader(controller.StatusBadRequest)
				} else {
					if len(gameState) != 3 || len(gameState[0]) != 3 ||
						len(gameState[1]) != 3 || len(gameState[2]) != 3 {
						log.Println("Invalid game state, array length is not valid!")

						res.WriteHeader(controller.StatusBadRequest)
					} else {
						whoWin := -1

						playerWin := isPlayerWin(gameState, 0, 0, clientPlayer, true)

						if playerWin {
							whoWin = clientWin
						} else {
							playerWin = isPlayerWin(gameState, 1, 1, clientPlayer, true)

							if playerWin {
								whoWin = clientWin
							} else {
								playerWin = isPlayerWin(gameState, 2, 2, clientPlayer, true)

								if playerWin {
									whoWin = clientWin
								} else {
									playerWin = tie(gameState)

									if playerWin {
										whoWin = tieGame
									}
								}
							}
						}

						var stmt *sql.Stmt
						var errStmt error

						query := ""

						if whoWin == -1 {
							query = "update games set game_state = ?, is_cpu_turn = ?, cpu_turn_started_at = ? where game_key = ?"
						} else {
							query = "update games set game_state = ?, finished_at = ?, who_win = ? where game_key = ?"
						}

						stmt, errStmt = dbConn.Prepare(query)

						if errStmt != nil {
							log.Println("DB prepare statement error!")
							log.Println(errStmt)

							res.WriteHeader(controller.StatusInternalServerError)
						} else {
							var queryRes sql.Result
							var errQueryExec error

							if whoWin == -1 {
								queryRes, errQueryExec = stmt.Exec(game.GameState, true, time.Now(), game.GameKey)
							} else {
								queryRes, errQueryExec = stmt.Exec(game.GameState, time.Now(), whoWin, game.GameKey)
							}

							if errQueryExec != nil {
								log.Println("DB query error!")
								log.Println(errQueryExec)

								res.WriteHeader(controller.StatusInternalServerError)
							} else {
								affRows, affRowsErr := queryRes.RowsAffected()

								if affRowsErr != nil || affRows == 0 {
									log.Println("DB query error, game state not updated!")
									log.Println(affRowsErr)

									res.WriteHeader(controller.StatusInternalServerError)
								} else {
									res.WriteHeader(controller.StatusOK)

									if whoWin == -1 {
										go cpuMove(game)
									}
								}
							}
						}
					}
				}
			}
		}
	} else if req.Method == controller.MethodOptions {
		res.WriteHeader(controller.StatusOK)
	} else {
		res.WriteHeader(controller.StatusMethodNotAllowed)
	}
}

func getGame(res controller.ResponseWriter, req *controller.Request) {
	res.Header().Set("Access-Control-Allow-Origin", "*")
	res.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, DELETE, PUT")
	res.Header().Set("Access-Control-Allow-Headers", "*")

	if req.Method == controller.MethodGet {
		res.Header().Set("Content-Type", "application/json")

		gameKey, exist := req.URL.Query()["game_key"]

		if !exist || len(gameKey[0]) != 12 {
			log.Println("Game key empty or not exist!")

			res.WriteHeader(controller.StatusBadRequest)
		} else {
			selRows, errSelRows := dbConn.Query("select game_key, username, difficulty, game_state, is_cpu_turn, cpu_turn_started_at, finished_at, who_win from games where game_key = ?", gameKey[0])

			if errSelRows != nil {
				log.Println("DB query error!")
				log.Println(errSelRows)

				res.WriteHeader(controller.StatusInternalServerError)
			} else {
				var game Game
				var scanErr error

				for selRows.Next() {
					var gameKey string
					var username string
					var difficulty int
					var gameState string
					var isCPUTurn bool
					var cpuTurnStartedAt sql.NullTime
					var finishedAt sql.NullTime
					var whoWin sql.NullInt32

					scanErr = selRows.Scan(&gameKey, &username, &difficulty, &gameState, &isCPUTurn, &cpuTurnStartedAt, &finishedAt, &whoWin)

					if scanErr != nil {
						log.Println("Query result scan error!")
						log.Println(scanErr)

						res.WriteHeader(controller.StatusInternalServerError)
					} else {
						if !isCPUTurn {
							game.GameKey = gameKey
							game.Username = username
							game.Difficulty = difficulty
							game.GameState = gameState
							game.GameFinished = (finishedAt.Valid)
							if whoWin.Valid {
								game.WhoWin = whoWin.Int32
							} else {
								game.WhoWin = -1
							}

							gameJSON, gameJSONErr := json.Marshal(game)

							if gameJSONErr != nil {
								log.Println("Converting game struct to json error!")
								log.Println(gameJSONErr)

								res.WriteHeader(controller.StatusInternalServerError)
							} else {
								res.WriteHeader(controller.StatusOK)

								res.Write([]byte(string(gameJSON)))
							}
						} else if time.Now().Sub(cpuTurnStartedAt.Time).Seconds() < cpuTurnTimeLimit {
							res.WriteHeader(controller.StatusOK)

							res.Write([]byte(`{"msg": "CPU makes move. Please wait."}`))
						} else {
							go cpuMove(game)

							res.WriteHeader(controller.StatusOK)

							res.Write([]byte(`{"msg": "CPU makes move. Please wait."}`))
						}
					}
				}
			}
		}
	} else {
		res.WriteHeader(controller.StatusMethodNotAllowed)
	}
}

func cpuMove(game Game) {
	time.Sleep(time.Duration(threadSleepTime) * time.Second)

	var gameState [][]int

	json.Unmarshal([]byte(game.GameState), &gameState)

	rows, cols := getEmptySlots(gameState)

	var r, c int

	var gameFinish bool

	if game.Difficulty == difficultyEasy {
		ns := rand.NewSource(time.Now().UnixNano())
		nr := rand.New(ns)
		r = rows[nr.Intn(len(rows))]

		ns = rand.NewSource(time.Now().UnixNano())
		nr = rand.New(ns)
		c = cols[nr.Intn(len(cols))]

		gameState[r][c] = cpuPlayer

		gameFinish = isPlayerWin(gameState, r, c, cpuPlayer, false)
	} else if game.Difficulty == difficultyImpossible {
		r, c, gameFinish = impossibleMove(gameState, rows, cols)

		gameState[r][c] = cpuPlayer
	}

	jsonGameState, _ := json.Marshal(gameState)

	query := ""

	if gameFinish {
		query = "update games set game_state = ?, is_cpu_turn = ?, finished_at = ?, who_win = ? where game_key = ?"
	} else {
		query = "update games set game_state = ?, is_cpu_turn = ? where game_key = ?"
	}

	stmt, errStmt := dbConn.Prepare(query)

	if errStmt != nil {
		log.Println("DB prepare statement error!")
		log.Println(errStmt)
	} else {
		var queryRes sql.Result
		var errQueryExec error

		if gameFinish {
			queryRes, errQueryExec = stmt.Exec(string(jsonGameState), false, time.Now(), cpuWin, game.GameKey)
		} else {
			queryRes, errQueryExec = stmt.Exec(string(jsonGameState), false, game.GameKey)
		}

		if errQueryExec != nil {
			log.Println("DB query error!")
			log.Println(errQueryExec)

		} else {
			affRows, affRowsErr := queryRes.RowsAffected()

			if affRowsErr != nil || affRows == 0 {
				log.Println("DB query error, game state not updated!")
				log.Println(affRowsErr)
			}
		}
	}
}

func isPlayerWin(gameState [][]int, r, c, p int, slotFilled bool) bool {
	if r == 0 && c == 0 {
		if (gameState[r][c+1] == p && gameState[r][c+2] == p) ||
			(gameState[r+1][c+1] == p && gameState[r+2][c+2] == p) ||
			(gameState[r+1][c] == p && gameState[r+2][c] == p) {
			if slotFilled && gameState[r][c] != p {
				return false
			} else {
				return true
			}
		}
	} else if r == 0 && c == 1 {
		if (gameState[r][c-1] == p && gameState[r][c+1] == p) ||
			(gameState[r+1][c] == p && gameState[r+2][c] == p) {
			if slotFilled && gameState[r][c] != p {
				return false
			} else {
				return true
			}
		}
	} else if r == 0 && c == 2 {
		if (gameState[r][c-1] == p && gameState[r][c-2] == p) ||
			(gameState[r+1][c-1] == p && gameState[r+2][c-p] == p) ||
			(gameState[r+1][c] == p && gameState[r+2][c] == p) {
			if slotFilled && gameState[r][c] != p {
				return false
			} else {
				return true
			}
		}
	} else if r == 1 && c == 0 {
		if (gameState[r][c+1] == p && gameState[r][c+2] == p) ||
			(gameState[r-1][c] == p && gameState[r+1][c] == p) {
			if slotFilled && gameState[r][c] != p {
				return false
			} else {
				return true
			}
		}
	} else if r == 1 && c == 1 {
		if (gameState[r][c-1] == p && gameState[r][c+1] == p) ||
			(gameState[r-1][c-1] == p && gameState[r+1][c+1] == p) ||
			(gameState[r-1][c] == p && gameState[r+1][c] == p) ||
			(gameState[r-1][c+1] == p && gameState[r+1][c-1] == p) {
			if slotFilled && gameState[r][c] != p {
				return false
			} else {
				return true
			}
		}
	} else if r == 1 && c == 2 {
		if (gameState[r][c-1] == p && gameState[r][c-2] == p) ||
			(gameState[r-1][c] == p && gameState[r+1][c] == p) {
			if slotFilled && gameState[r][c] != p {
				return false
			} else {
				return true
			}
		}
	} else if r == 2 && c == 0 {
		if (gameState[r][c+1] == p && gameState[r][c+2] == p) ||
			(gameState[r-1][c] == p && gameState[r-2][c] == p) ||
			(gameState[r-1][c+1] == p && gameState[r-2][c+2] == p) {
			if slotFilled && gameState[r][c] != p {
				return false
			} else {
				return true
			}
		}
	} else if r == 2 && c == 1 {
		if (gameState[r][c-1] == p && gameState[r][c+1] == p) ||
			(gameState[r-1][c] == p && gameState[r-2][c] == p) {
			if slotFilled && gameState[r][c] != p {
				return false
			} else {
				return true
			}
		}
	} else {
		if (gameState[r][c-1] == p && gameState[r][c-2] == p) ||
			(gameState[r-1][c-1] == p && gameState[r-2][c-2] == p) ||
			(gameState[r-1][c] == p && gameState[r-2][c] == p) {
			if slotFilled && gameState[r][c] != p {
				return false
			} else {
				return true
			}
		}
	}

	return false
}

func impossibleMove(gameState [][]int, rows, cols []int) (int, int, bool) {
	var s []int
	p := -1

	for i := 0; i < len(rows); i++ {
		r := rows[i]
		c := cols[i]

		if r == 0 && c == 0 {
			if isPlayerWin(gameState, r, c, cpuPlayer, false) {
				return r, c, true
			} else if (gameState[r][c+1] == clientPlayer && gameState[r][c+2] == clientPlayer) ||
				(gameState[r+1][c+1] == clientPlayer && gameState[r+2][c+2] == clientPlayer) ||
				(gameState[r+1][c] == clientPlayer && gameState[r+2][c] == clientPlayer) {
				p = i
			} else if gameState[r][c+1] == clientPlayer || gameState[r+1][c+1] == clientPlayer || gameState[r+1][c] == clientPlayer {
				s = append(s, i)
			}
		} else if r == 0 && c == 1 {
			if isPlayerWin(gameState, r, c, cpuPlayer, false) {
				return r, c, true
			} else if (gameState[r][c-1] == clientPlayer && gameState[r][c+1] == clientPlayer) ||
				(gameState[r+1][c] == clientPlayer && gameState[r+2][c] == clientPlayer) {
				p = i
			} else if gameState[r][c-1] == clientPlayer || gameState[r][c+1] == clientPlayer || gameState[r+1][c] == clientPlayer {
				s = append(s, i)
			}
		} else if r == 0 && c == 2 {
			if isPlayerWin(gameState, r, c, cpuPlayer, false) {
				return r, c, true
			} else if (gameState[r][c-1] == clientPlayer && gameState[r][c-2] == clientPlayer) ||
				(gameState[r+1][c-1] == clientPlayer && gameState[r+2][c-2] == clientPlayer) ||
				(gameState[r+1][c] == clientPlayer && gameState[r+2][c] == clientPlayer) {
				p = i
			} else if gameState[r][c-1] == clientPlayer || gameState[r+1][c-1] == clientPlayer || gameState[r+1][c] == clientPlayer {
				s = append(s, i)
			}
		} else if r == 1 && c == 0 {
			if isPlayerWin(gameState, r, c, cpuPlayer, false) {
				return r, c, true
			} else if (gameState[r][c+1] == clientPlayer && gameState[r][c+2] == clientPlayer) ||
				(gameState[r-1][c] == clientPlayer && gameState[r+1][c] == clientPlayer) {
				p = i
			} else if gameState[r][c+1] == clientPlayer || gameState[r-1][c] == clientPlayer || gameState[r+1][c] == clientPlayer {
				s = append(s, i)
			}
		} else if r == 1 && c == 1 {
			if isPlayerWin(gameState, r, c, cpuPlayer, false) {
				return r, c, true
			} else if (gameState[r][c-1] == clientPlayer && gameState[r][c+1] == clientPlayer) ||
				(gameState[r-1][c-1] == clientPlayer && gameState[r+1][c+1] == clientPlayer) ||
				(gameState[r-1][c] == clientPlayer && gameState[r+1][c] == clientPlayer) ||
				(gameState[r-1][c+1] == clientPlayer && gameState[r+1][c-1] == clientPlayer) {
				p = i
			} else if gameState[r][c-1] == clientPlayer || gameState[r][c+1] == clientPlayer || gameState[r-1][c-1] == clientPlayer ||
				gameState[r+1][c+1] == clientPlayer || gameState[r-1][c] == clientPlayer || gameState[r+1][c] == clientPlayer ||
				gameState[r-1][c+1] == clientPlayer || gameState[r+1][c-1] == clientPlayer {
				s = append(s, i)
			}
		} else if r == 1 && c == 2 {
			if isPlayerWin(gameState, r, c, cpuPlayer, false) {
				return r, c, true
			} else if (gameState[r][c-1] == clientPlayer && gameState[r][c-2] == clientPlayer) ||
				(gameState[r-1][c] == clientPlayer && gameState[r+1][c] == clientPlayer) {
				p = i
			} else if gameState[r][c-1] == clientPlayer || gameState[r-1][c] == clientPlayer || gameState[r+1][c] == clientPlayer {
				s = append(s, i)
			}
		} else if r == 2 && c == 0 {
			if isPlayerWin(gameState, r, c, cpuPlayer, false) {
				return r, c, true
			} else if (gameState[r][c+1] == clientPlayer && gameState[r][c+2] == clientPlayer) ||
				(gameState[r-1][c] == clientPlayer && gameState[r-2][c] == clientPlayer) ||
				(gameState[r-1][c+1] == clientPlayer && gameState[r-2][c+2] == clientPlayer) {
				p = i
			} else if gameState[r][c+1] == clientPlayer || gameState[r-1][c] == clientPlayer || gameState[r-1][c+1] == clientPlayer {
				s = append(s, i)
			}
		} else if r == 2 && c == 1 {
			if isPlayerWin(gameState, r, c, cpuPlayer, false) {
				return r, c, true
			} else if (gameState[r][c-1] == clientPlayer && gameState[r][c+1] == clientPlayer) ||
				(gameState[r-1][c] == clientPlayer && gameState[r-2][c] == clientPlayer) {
				p = i
			} else if gameState[r][c-1] == clientPlayer || gameState[r][c+1] == clientPlayer || gameState[r-1][c] == clientPlayer {
				s = append(s, i)
			}
		} else {
			if isPlayerWin(gameState, r, c, cpuPlayer, false) {
				return r, c, true
			} else if (gameState[r][c-1] == clientPlayer && gameState[r][c-2] == clientPlayer) ||
				(gameState[r-1][c-1] == clientPlayer && gameState[r-2][c-2] == clientPlayer) ||
				(gameState[r-1][c] == clientPlayer && gameState[r-2][c] == clientPlayer) {
				p = i
			} else if gameState[r][c-1] == clientPlayer || gameState[r-1][c-1] == clientPlayer || gameState[r-1][c] == clientPlayer {
				s = append(s, i)
			}
		}
	}

	if p != -1 {
		return rows[p], cols[p], false
	}

	ns := rand.NewSource(time.Now().UnixNano())
	nr := rand.New(ns)
	i := s[nr.Intn(len(s))]

	return rows[i], cols[i], false
}
