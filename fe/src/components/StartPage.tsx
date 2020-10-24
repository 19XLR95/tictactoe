import React from "react";
import "../styles/start-page.css";
import Game from "../types/Game";
import NewGame from "../types/NewGame";
import * as TicTacToeAPI from "../api/TicTacToeAPI";
import GamePage from "./GamePage";
import HighScoresPage from "./HighScoresPage";

type DefaultProps = {};

type DefaultState = {
    existGameKey: string | null,
    game: Game | null,
    highScoresPage: boolean,
    msg: string | null
};

export default class StartPage extends React.Component<DefaultProps, DefaultState> {
    private inputUsernameRef: React.RefObject<HTMLInputElement>;
    private inputDifficultyRef: React.RefObject<HTMLSelectElement>;

    constructor(props: DefaultProps) {
        super(props);

        this.state = {
            existGameKey: process.env.REACT_APP_GAME_KEY_LOCAL_STORAGE !== undefined ? 
                localStorage.getItem(process.env.REACT_APP_GAME_KEY_LOCAL_STORAGE) : null,
            game: null,
            highScoresPage: false,
            msg: null
        }

        this.inputUsernameRef = React.createRef();
        this.inputDifficultyRef = React.createRef();

        this.onStartClick = this.onStartClick.bind(this);
        this.openStartPage = this.openStartPage.bind(this);
        this.onResumeClick = this.onResumeClick.bind(this);
        this.onHighScoresClick = this.onHighScoresClick.bind(this);
    }

    onStartClick() {
        if(this.inputUsernameRef.current?.value && this.inputDifficultyRef.current?.value && this.inputUsernameRef.current.value.length > 0) {
            this.setState({
                msg: "Please wait..."
            });
            
            let newGame : NewGame = {
                username: this.inputUsernameRef.current.value,
                difficulty: parseInt(this.inputDifficultyRef.current.value)
            };

            TicTacToeAPI.createNewGame(
                newGame,
                (game: Game | null) => {
                    if(game !== null) {
                        this.setState({
                            game: game,
                            existGameKey: game.gameKey
                        });

                        process.env.REACT_APP_GAME_KEY_LOCAL_STORAGE && localStorage.setItem(process.env.REACT_APP_GAME_KEY_LOCAL_STORAGE, game.gameKey);
                    } else {
                        this.setState({
                            msg: "Cannot create a new game after retries!"
                        });
                    }
                }
            );
        } else {
            this.setState({
                msg: "Please enter a username!"
            });
        }
    }

    openStartPage() {
        if(process.env.REACT_APP_GAME_KEY_LOCAL_STORAGE && localStorage.getItem(process.env.REACT_APP_GAME_KEY_LOCAL_STORAGE)) {
            this.setState({
                game: null,
                highScoresPage: false,
                msg: null
            });
        } else {
            this.setState({
                game: null,
                existGameKey: null,
                highScoresPage: false,
                msg: null
            });
        }
    }

    onHighScoresClick() {
        if(process.env.REACT_APP_GAME_KEY_LOCAL_STORAGE && localStorage.getItem(process.env.REACT_APP_GAME_KEY_LOCAL_STORAGE)) {
            this.setState({
                highScoresPage: true,
                game: null,
                msg: null
            });
        } else {
            this.setState({
                existGameKey: null,
                highScoresPage: true,
                game: null,
                msg: null
            });
        }
    }

    onResumeClick() {
        if(this.state.existGameKey) {
            this.setState({
                msg: "Please wait..."
            });

            TicTacToeAPI.getGame(
                this.state.existGameKey,
                (game: Game | null) => {
                    if(game !== null) {
                        this.setState({
                            game: game,
                            highScoresPage: false
                        });
                    } else {
                        this.setState({
                            msg: "Cannot retrieve the game after retries!"
                        });
                    }
                }
            );
        }
    }
    
    render() {
        return (
            <>
                {
                    this.state.game !== null ? 
                    <GamePage game={this.state.game} openStartPage={this.openStartPage} onHighScoresClick={this.onHighScoresClick} /> :
                    this.state.highScoresPage ?
                    <HighScoresPage openStartPage={this.openStartPage} onResumeClick={this.onResumeClick} resumeGame={this.state.existGameKey !== null} /> :
                    <div id="start-page">
                        <div>
                            <h1>Tic Tac Toe</h1>
                            <div>
                                <label htmlFor="username">Username: </label>
                                <input type="text" name="username" id="username" ref={this.inputUsernameRef}/>
                                <label htmlFor="difficulty">Difficulty: </label>
                                <select name="difficulty" id="difficulty" defaultValue="1" ref={this.inputDifficultyRef}>
                                    <option value="1">Easy</option>
                                    <option value="2">Impossible</option>
                                </select>
                                <button onClick={this.onStartClick}>Start</button>
                                {this.state.existGameKey && <button onClick={this.onResumeClick}>Resume</button>}
                                <button onClick={this.onHighScoresClick}>High Scores</button>
                            </div>
                            <div style={{color: "red", height: "22px"}}>{this.state.msg}</div>
                        </div>
                    </div>
                }
            </>
        );
    }
}
