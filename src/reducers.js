import { combineReducers } from 'redux'
import {
  SET_INPUT,
  ADD_ROOMS,
  JOINED_ROOM
} from './actions'

function user (state = {}, action) {
  switch (action.type) {
  default:
    return state
  }
}

function currentClash (state = {}, action) {

  switch (action.type) {
  case SET_INPUT:
    return Object.assign({}, state, {value:action.value})
  default:
    return state
  }
}

const initialRoom = {
  id: -1,
  name: "room",
  people: []
}

function currentRoom (state = initialRoom, action) {
  switch (action.type) {
  case JOINED_ROOM:
    return Object.assign({}, action.room) 
  default:
    return state
  }
}

function rooms (state = [], action) {

  switch (action.type) {
  case ADD_ROOMS:
    return action.rooms
  default:
    return state
  }
}

function subscriptions (state = [], action) {

  switch (action.type) {
  default:
    return state
  }
}

const rootReducer = combineReducers({
  user,
  currentClash,
  rooms,
  subscriptions
})

export default rootReducer
