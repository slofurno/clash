import React, { Component, PropTypes } from 'react'
import { connect } from 'react-redux'
import Clash from './Clash'
import RoomList from './RoomList'
import Lobby from './Lobby'

import {
  setInput,
  joinRoom,
  postCode
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
      setInput,
      currentClash,
      joinRoom,
      postCode

    } = this.props

    return (
      <div>
        <RoomList rooms = {rooms} joinRoom = {joinRoom}/>
        <Lobby room = {currentRoom} events = {events}/>
        <Clash 
          setInput = {setInput} 
          clash = {currentClash}
          postCode = {postCode}
        />
      </div>
    )
  }
}

let selector = x => x
export default connect(selector, mapDispatchToProps)(App)
