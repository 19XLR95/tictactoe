import React from "react";
import "../styles/game-page.css"
import Game from "../types/Game";
import * as TicTacToeAPI from "../api/TicTacToeAPI";

type DefaultProps = {
    game: Game,
    openStartPage: () => void,
    onHighScoresClick: () => void
};

type DefaultState = {
    game: Game,
    playerTurn: boolean,
    msg: string | null
};

export default class GamePage extends React.Component<DefaultProps, DefaultState> {
    constructor(props: DefaultProps) {
        super(props);

        this.state = {
            game: props.game,
            playerTurn: true,
            msg: null
        };

        this.onCellClick = this.onCellClick.bind(this);
    }

    onCellClick(i: number, j: number) {
        if(this.state.playerTurn && !this.state.game.gameFinished && this.state.game.gameState[i][j] === 0) {
            let game = this.state.game;
            game.gameState[i][j] = 1;

            this.setState({
                game: game,
                playerTurn: false,
                msg: "Please wait..."
            });

            TicTacToeAPI.makeMove(
                this.state.game,
                (result: boolean) => {
                    if(result) {
                        TicTacToeAPI.getGame(
                            this.state.game.gameKey,
                            (gameResponse: Game | null) => {
                                if(gameResponse !== null) {
                                    this.setState({
                                        game: gameResponse,
                                        playerTurn: true
                                    });

                                    if(gameResponse.gameFinished) {
                                        process.env.REACT_APP_GAME_KEY_LOCAL_STORAGE && localStorage.removeItem(process.env.REACT_APP_GAME_KEY_LOCAL_STORAGE);
                                    }

                                    this.setState({
                                        msg: null
                                    });
                                } else {
                                    this.setState({
                                        msg: "Cannot retrieve the game after retries!"
                                    });
                                }
                            }
                        );
                    } else {
                        game.gameState[i][j] = 0;

                        this.setState({
                            game: game,
                            playerTurn: true,
                            msg: "Cannot make move after retries!"
                        });
                    }
                }
            );
        }
    }

    render() {
        let color: string = !this.state.game.gameFinished ? 
            (this.state.playerTurn ? "green" : "red") :
            (this.state.game.whoWin === (process.env.REACT_APP_CLIENT_WIN && parseInt(process.env.REACT_APP_CLIENT_WIN)) ? "green" : 
            this.state.game.whoWin === (process.env.REACT_APP_CPU_WIN && parseInt(process.env.REACT_APP_CPU_WIN)) ? "red" : "#FFB900");

        let text: string = !this.state.game.gameFinished ? 
            (this.state.playerTurn ? "Your Turn" : "CPU's Turn") :
            (this.state.game.whoWin === (process.env.REACT_APP_CLIENT_WIN && parseInt(process.env.REACT_APP_CLIENT_WIN)) ? "You won!" : 
            this.state.game.whoWin === (process.env.REACT_APP_CPU_WIN && parseInt(process.env.REACT_APP_CPU_WIN)) ? "CPU Won!" : "Tie!");

        return (
            <div id="game-page">
                <div>
                    <div>
                        <h1>{this.state.game.username}</h1>
                        <h2 style={{color: color}}>{text}</h2>
                    </div>
                    <div>
                        <div onClick={() => {this.onCellClick(0, 0)}}>
                            {this.state.game.gameState[0][0] !== 0 && (this.state.game.gameState[0][0] === 1 ? "X" : "O")}
                        </div>
                        <div onClick={() => {this.onCellClick(0, 1)}}>
                            {this.state.game.gameState[0][1] !== 0 && (this.state.game.gameState[0][1] === 1 ? "X" : "O")}
                        </div>
                        <div onClick={() => {this.onCellClick(0, 2)}}>
                            {this.state.game.gameState[0][2] !== 0 && (this.state.game.gameState[0][2] === 1 ? "X" : "O")}
                        </div>
                        <div onClick={() => {this.onCellClick(1, 0)}}>
                            {this.state.game.gameState[1][0] !== 0 && (this.state.game.gameState[1][0] === 1 ? "X" : "O")}
                        </div>
                        <div onClick={() => {this.onCellClick(1, 1)}}>
                            {this.state.game.gameState[1][1] !== 0 && (this.state.game.gameState[1][1] === 1 ? "X" : "O")}
                        </div>
                        <div onClick={() => {this.onCellClick(1, 2)}}>
                            {this.state.game.gameState[1][2] !== 0 && (this.state.game.gameState[1][2] === 1 ? "X" : "O")}
                        </div>
                        <div onClick={() => {this.onCellClick(2, 0)}}>
                            {this.state.game.gameState[2][0] !== 0 && (this.state.game.gameState[2][0] === 1 ? "X" : "O")}
                        </div>
                        <div onClick={() => {this.onCellClick(2, 1)}}>
                            {this.state.game.gameState[2][1] !== 0 && (this.state.game.gameState[2][1] === 1 ? "X" : "O")}
                        </div>
                        <div onClick={() => {this.onCellClick(2, 2)}}>
                            {this.state.game.gameState[2][2] !== 0 && (this.state.game.gameState[2][2] === 1 ? "X" : "O")}
                        </div>
                    </div>
                    <div>
                        <button onClick={this.props.openStartPage}>New Game</button>
                        <button onClick={this.props.onHighScoresClick}>High Scores</button>
                    </div>
                    <div style={{color: "red", textAlign: "center", height: "22px"}}>{this.state.msg}</div>
                </div>
            </div>
        );
    }
}
