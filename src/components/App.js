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
  postResult
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
    }
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
      clashResults,
      setInput,
      currentClash,
      joinRoom,
      postCode,
      postResult

    } = this.props

    let visibleResults = clashResults.filter(x => x.subject === currentClash.id)
    let visibleClashResults = clashResults.filter(x => x.clash === currentClash.id)

    return (
      <div>
        <div style = {{display: "inline-block"}}>
          <Lobby room = {currentRoom} events = {events} users = {users} />
          <Clash 
            setInput = {setInput} 
            clash = {currentClash}
            postCode = {postCode}
            postResult = {postResult}
            results = {results}
            visibleResults = {visibleResults}
            users = {users}
          />
          <ClashResults
            results = {visibleClashResults}             
          />
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
