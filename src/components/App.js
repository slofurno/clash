import React, { Component, PropTypes } from 'react'
import { connect } from 'react-redux'
import Clash from './Clash'
import RoomList from './RoomList'

import {
  setInput,
  joinRoom
} from '../actions'


const mapDispatchToProps = (dispatch) => {
  return {
		setInput: (e) => {
      dispatch(setInput(e))
		},
    joinRoom: (x) => {
      dispatch(joinRoom(x))
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
      setInput,
      currentClash,
      joinRoom
    } = this.props

    return (
      <div>
        <RoomList rooms = {rooms} joinRoom = {joinRoom}/>
        <Clash 
          setInput={setInput} 
          value = {currentClash.input}
        />
      </div>
    )
  }
}

let selector = x => x
export default connect(selector, mapDispatchToProps)(App)
