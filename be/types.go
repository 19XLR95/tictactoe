package main

// NewGame struct
type NewGame struct {
	Username   string `json:"username"`
	Difficulty int    `json:"difficulty"`
}

// Game struct
type Game struct {
	GameKey      string `json:"game_key"`
	Username     string `json:"username"`
	Difficulty   int    `json:"difficulty"`
	GameState    string `json:"game_state"`
	GameFinished bool   `json:"game_finished"`
	WhoWin       int32  `json:"who_win"`
}

// Score struct
type Score struct {
	Username   string `json:"username"`
	FinishedIn int64  `json:"finished_in"`
	WhoWin     int    `json:"who_win"`
}
