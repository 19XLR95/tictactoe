import React from "react";
import "../styles/high-scores-page.css";
import HighScore from "../types/HighScore";
import * as TicTacToeAPI from "../api/TicTacToeAPI";

type DefaultProps = {
    openStartPage: () => void,
    onResumeClick: () => void,
    resumeGame: boolean
};

type DefaultState = {
    prevHighScores: Array<HighScore> | null,
    currentHighScores: Array<HighScore> | null,
    nextHighScores: Array<HighScore> | null,
    offset: number,
    msg: string | null
};

export default class HighScoresPage extends React.Component<DefaultProps, DefaultState> {
    private readonly LIMIT: number = 10;

    constructor(props: DefaultProps) {
        super(props);

        this.state = {
            prevHighScores: null,
            currentHighScores: null,
            nextHighScores: null,
            offset: 0,
            msg: "Please wait..."
        };

        this.onPrevClick = this.onPrevClick.bind(this);
        this.onNextClick = this.onNextClick.bind(this);
    }

    componentDidMount() {
        TicTacToeAPI.highScores(
            this.state.offset, this.LIMIT,
            (currentHighScores: Array<HighScore> | null) => {
                if(currentHighScores !== null) {
                    TicTacToeAPI.highScores(
                        this.state.offset + this.LIMIT, this.LIMIT,
                        (nextHighScores: Array<HighScore> | null) => {
                            if(nextHighScores !== null) {
                                if(nextHighScores.length > 0) {
                                    this.setState({
                                        currentHighScores: currentHighScores,
                                        nextHighScores: nextHighScores,
                                        msg: null
                                    });
                                } else {
                                    this.setState({
                                        currentHighScores: currentHighScores,
                                        msg: null
                                    });
                                }
                            } else {
                                this.setState({
                                    msg: "Cannot retrieve high scores after retries!"
                                });
                            }
                        }
                    );
                } else {
                    this.setState({
                        msg: "Cannot retrieve high scores after retries!"
                    });
                }
            }
        );
    }

    onPrevClick() {
        if(this.state.prevHighScores) {
            this.setState({
                msg: "Please wait..."
            });

            if(this.state.offset - (2 * this.LIMIT) >= 0) {
                TicTacToeAPI.highScores(
                    this.state.offset - (2 * this.LIMIT), this.LIMIT,
                    (prevHighScores: Array<HighScore> | null) => {
                        if(prevHighScores !== null && prevHighScores.length > 0) {
                            this.setState(
                                (prevState: DefaultState) => {
                                    return {
                                        nextHighScores: prevState.currentHighScores,
                                        currentHighScores: prevState.prevHighScores,
                                        prevHighScores: prevHighScores,
                                        offset: prevState.offset - this.LIMIT,
                                        msg: null
                                    };
                                }
                            );
                        } else {
                            this.setState({
                                msg: "Cannot retrieve high scores after retries!"
                            });
                        }
                    }
                );
            } else {
                this.setState(
                    (prevState: DefaultState) => {
                        return {
                            nextHighScores: prevState.currentHighScores,
                            currentHighScores: prevState.prevHighScores,
                            prevHighScores: null,
                            offset: prevState.offset - this.LIMIT,
                            msg: null
                        };
                    }
                );
            }
        }
    }

    onNextClick() {
        if(this.state.nextHighScores) {
            this.setState({
                msg: "Please wait..."
            });

            TicTacToeAPI.highScores(
                this.state.offset + (2 * this.LIMIT), this.LIMIT,
                (nextHighScores: Array<HighScore> | null) => {
                    if(nextHighScores !== null) {
                        if(nextHighScores.length > 0) {
                            this.setState(
                                (prevState: DefaultState) => {
                                    return {
                                        prevHighScores: prevState.currentHighScores,
                                        currentHighScores: prevState.nextHighScores,
                                        nextHighScores: nextHighScores,
                                        offset: prevState.offset + this.LIMIT,
                                        msg: null
                                    };
                                }
                            );
                        } else {
                            this.setState(
                                (prevState: DefaultState) => {
                                    return {
                                        prevHighScores: prevState.currentHighScores,
                                        currentHighScores: prevState.nextHighScores,
                                        nextHighScores: null,
                                        offset: prevState.offset + this.LIMIT,
                                        msg: null
                                    };
                                }
                            );
                        }
                    } else {
                        this.setState({
                            msg: "Cannot retrieve high scores after retries!"
                        });
                    }
                }
            );
        }
    }

    render() {
        return (
            <div id="high-scores-page">
                <div>
                    <div>
                        <h1>High Scores</h1>
                        <h2>Wins</h2>
                    </div>
                    <div>
                        <div><b>Username</b></div>
                        <div><b>Finished In (Seconds)</b></div>
                        {
                            this.state.currentHighScores && this.state.currentHighScores.map(
                                (e: HighScore, i: number) => {
                                    return (
                                        <div key={i} className="high-score-cell">
                                            <div >{e.username}</div>
                                            <div >{e.finishedIn}</div>
                                        </div>
                                    );
                                }
                            )
                        }
                    </div>
                    <div>
                        <button onClick={this.onPrevClick} disabled={this.state.prevHighScores === null}>&lt;</button>
                        <button onClick={this.onNextClick} disabled={this.state.nextHighScores === null}>&gt;</button>
                        {this.props.resumeGame && <button onClick={this.props.onResumeClick}>Resume</button>}
                        <button onClick={this.props.openStartPage}>New Game</button>
                    </div>
                    <div style={{color: "red", textAlign: "center", height: "22px"}}>{this.state.msg}</div>
                </div>
            </div>
        );
    }
}
