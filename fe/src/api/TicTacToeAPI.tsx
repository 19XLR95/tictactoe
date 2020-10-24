import NewGame from "../types/NewGame";
import axios, { AxiosResponse } from "axios";
import Game from "../types/Game";
import HighScore from "../types/HighScore";

function setExponentialBackOffRetry(f: Function, retry: number) {
    setTimeout(f, Math.pow(2, retry) + Math.floor(Math.random() * 100));
}

function createNewGameAPICall(newGame: NewGame, callback: (game: Game | null) => void, retry: number) {
    if(process.env.REACT_APP_TTT_BASE_URL && process.env.REACT_APP_TTT_NEW_GAME_ENDPOINT) {
        axios.post(process.env.REACT_APP_TTT_BASE_URL + process.env.REACT_APP_TTT_NEW_GAME_ENDPOINT,
            {
                username: newGame.username,
                difficulty: newGame.difficulty
            }
        ).then(
            (res: AxiosResponse) => {
                if(res.status === 201) {
                    let game: Game = {
                        gameKey: res.data.game_key,
                        difficulty: res.data.difficulty,
                        gameFinished: res.data.game_finished,
                        gameState: JSON.parse(res.data.game_state),
                        username: res.data.username,
                        whoWin: res.data.who_win
                    };

                    callback(game);
                }
            }
        ).catch(
            (err: Error) => {
                console.error(err);

                let r = retry + 1;

                if(process.env.REACT_APP_MAX_API_RETRY && r <= parseInt(process.env.REACT_APP_MAX_API_RETRY)) {
                    setExponentialBackOffRetry(() => {
                        createNewGameAPICall(newGame, callback, r);
                    }, r);
                } else {
                    callback(null);
                }
            }
        );
    }
}

export function createNewGame(newGame: NewGame, callback: (game: Game | null) => void) {
    createNewGameAPICall(newGame, callback, 0);
}

function makeMoveAPICall(game: Game, callback: (result: boolean) => void, retry: number) {
    if(process.env.REACT_APP_TTT_BASE_URL && process.env.REACT_APP_TTT_MAKE_MOVE_ENDPOINT) {
        axios.post(process.env.REACT_APP_TTT_BASE_URL + process.env.REACT_APP_TTT_MAKE_MOVE_ENDPOINT,
            {
                game_key: game.gameKey,
                username: game.username,
                difficulty: game.difficulty,
                game_state: JSON.stringify(game.gameState),
                game_finished: game.gameFinished,
                who_win: game.whoWin
            }
        ).then(
            (res: AxiosResponse) => {
                if(res.status === 200) {
                    callback(true);
                }
            }
        ).catch(
            (err: Error) => {
                console.error(err);

                let r = retry + 1;

                if(process.env.REACT_APP_MAX_API_RETRY && r <= parseInt(process.env.REACT_APP_MAX_API_RETRY)) {
                    setExponentialBackOffRetry(() => {
                        makeMoveAPICall(game, callback, r)
                    }, r);
                } else {
                    callback(false);
                }
            }
        );
    }
}

export function makeMove(game: Game, callback: (result: boolean) => void) {
    makeMoveAPICall(game, callback, 0);
}

function getGameAPICall(gameKey: string, callback: (game: Game | null) => void, retry: number) {
    if(process.env.REACT_APP_TTT_BASE_URL && process.env.REACT_APP_TTT_GET_GAME_ENDPOINT) {
        axios.get(process.env.REACT_APP_TTT_BASE_URL + process.env.REACT_APP_TTT_GET_GAME_ENDPOINT,
            {
                params: {
                    game_key: gameKey
                }
            }
        ).then(
            (res: AxiosResponse) => {
                if(res.status === 200) {
                    if(!res.data.hasOwnProperty("msg")) {
                        let game: Game = {
                            gameKey: res.data.game_key,
                            difficulty: res.data.difficulty,
                            gameFinished: res.data.game_finished,
                            gameState: JSON.parse(res.data.game_state),
                            username: res.data.username,
                            whoWin: res.data.who_win
                        };
    
                        callback(game);
                    } else {
                        process.env.REACT_APP_GET_GAME_DELAY && setTimeout(
                            () => {
                                getGameAPICall(gameKey, callback, retry);
                            }, parseInt(process.env.REACT_APP_GET_GAME_DELAY)
                        );
                    }
                }
            }
        ).catch(
            (err: Error) => {
                console.error(err);

                let r = retry + 1;

                if(process.env.REACT_APP_MAX_API_RETRY && r <= parseInt(process.env.REACT_APP_MAX_API_RETRY)) {
                    setExponentialBackOffRetry(() => {
                        getGameAPICall(gameKey, callback, r)
                    }, r);
                } else {
                    callback(null);
                }

            }
        );
    }
}

export function getGame(gameKey: string, callback: (game: Game | null) => void) {
    getGameAPICall(gameKey, callback, 0);
}

function highScoresAPICall(offset: number, limit: number, callback: (highScores: Array<HighScore> | null) => void, retry: number) {
    if(process.env.REACT_APP_TTT_BASE_URL && process.env.REACT_APP_TTT_HIGH_SCORES_ENDPOINT) {
        axios.get(process.env.REACT_APP_TTT_BASE_URL + process.env.REACT_APP_TTT_HIGH_SCORES_ENDPOINT,
            {
                params: {
                    offset: offset,
                    limit: limit
                }
            }
        ).then(
            (res: AxiosResponse) => {
                if(res.status === 200 && res.data !== null && res.data.length > 0) {
                    let highScores = Array<HighScore>();
                    res.data.forEach(
                        (e: any) => {
                            highScores.push({
                                finishedIn: e.finished_in,
                                username: e.username,
                                whoWin: e.who_win
                            });
                        }
                    );

                    callback(highScores);
                } else {
                    callback([]);
                }
            }
        ).catch(
            (err: Error) => {
                console.error(err);

                let r = retry + 1;

                if(process.env.REACT_APP_MAX_API_RETRY && r <= parseInt(process.env.REACT_APP_MAX_API_RETRY)) {
                    setExponentialBackOffRetry(() => {
                        highScoresAPICall(offset, limit, callback, r)
                    }, r);
                } else {
                    callback(null);
                }

            }
        );
    }
}

export function highScores(offset: number, limit: number, callback: (highScores: Array<HighScore> | null) => void) {
    highScoresAPICall(offset, limit, callback, 0);
}
