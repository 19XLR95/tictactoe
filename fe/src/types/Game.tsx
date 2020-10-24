type Game = {
    gameKey: string,
    username: string,
    difficulty: number,
    gameState: Array<Array<number>>,
    gameFinished: boolean,
    whoWin: number
}

export default Game;
