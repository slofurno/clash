import React, { Component, PropTypes } from 'react'
import { connect } from 'react-redux'
import Clash from './Clash'
import RoomList from './RoomList'
import Lobby from './Lobby'
import ClashResults from './ClashResults'

import {
  setInput,
  joinRoom,
  postCode,
  postResult,
  setSlide,
  setProblem
} from '../actions'


const mapDispatchToProps = (dispatch) => {
  return {
		setInput: (e) => {
      dispatch(setInput(e))
		},
    joinRoom: (x) => {
      dispatch(joinRoom(x))
    },
    postCode: (clash, code) => {
      dispatch(postCode(clash, code))
    },
    postResult: (clash, code) => {
      dispatch(postResult(clash, code))
    },
    setSlide: (slide) => {
      dispatch(setSlide(slide))
    },
    setProblem: (problem, clash) => {
      dispatch(setProblem(problem, clash))
    },
    /*onTodoClick: (id) => {
      dispatch(toggleTodo(id))
    }*/
  }
}

class App extends Component {
  render () {
    const { 
      rooms,
      currentRoom,
      events,
      users,
      results,
      clashes,
      problems,
      currentProblem,
      clashResults,
      slide,
      setInput,
      currentClash,
      joinRoom,
      postCode,
      postResult,
      setSlide,
      setProblem

    } = this.props

    let visibleResults = clashResults.filter(x => x.subject === currentClash.id)
    let visibleClashResults = clashResults.filter(x => x.clash === currentClash.id)
    let visibleProblem = problems[currentProblem] || {}

    let visibleSlide = (function(){
      switch(slide) {
      case "LOBBY":
        return (
            <Lobby 
              room = {currentRoom} 
              events = {events}
              users = {users}
              clashes = {clashes}
              setProblem = {setProblem}
            />
        )
      case "CLASH":
        return (
          <Clash 
            setInput = {setInput} 
            clash = {currentClash}
            postCode = {postCode}
            postResult = {postResult}
            results = {results}
            visibleResults = {visibleResults}
            users = {users}
            visibleProblem = {visibleProblem}
          />
        )
      case "RESULTS":
        return (
          <ClashResults
            results = {visibleClashResults}             
          />
        )
      default:
        return (<div></div>)
      }
    }())

    return (
      <div>
        <div style = {{display: "inline-block"}}>
          <div style = {{
            width: "800px",
            height: "400px",
            padding: "6px",
            backgroundColor: "silver",
          }}>
            {visibleSlide}
          </div>
          <div>
            <ul>
              <li className="tab" onClick = {e => setSlide("LOBBY")}>Lobby</li>
              <li className="tab" onClick = {e => setSlide("CLASH")}>Clash</li>
              <li className="tab" onClick = {e => setSlide("RESULTS")}>Results</li>
            </ul>
          </div>
        </div>
        <div style = {{
          display: "inline-block",
          verticalAlign: "top",
          width: "150px"
        }}>
          <RoomList rooms = {rooms} joinRoom = {joinRoom}/>
        </div>
      </div>
    )
  }
}

let selector = x => x
export default connect(selector, mapDispatchToProps)(App)
