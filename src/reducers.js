import { combineReducers } from 'redux'
import {
  SET_INPUT,
  ADD_ROOMS,
  JOINED_ROOM,
  ADD_EVENTS,
  ADD_PROBLEMS,
  SHOW_CLASH,
  JOIN_CLASH
} from './actions'

function user (state = {}, action) {
  switch (action.type) {
  default:
    return state
  }
}

const initialClash = {
  id: -1,
  input: ""
}


function currentClash (state = {}, action) {
  switch (action.type) {
  case JOIN_CLASH:
    console.log("join", action)
    return Object.assign({}, state, {id: action.id})
  case SHOW_CLASH:
    return Object.assign({}, state, {value: state.value, problem: action.problem})
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
    return Object.assign({}, initialRoom, action.room) 
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

function events (state = [], action) {
  switch (action.type) {
  case ADD_EVENTS:
    return state.concat(action.events)
  default:
    return state
  }
} 

function problems (state = [], action) {
  switch (action.type) {
  case ADD_PROBLEMS:
    return state.concat(action.problems)
  default:
    return state
  }
}



const rootReducer = combineReducers({
  user,
  currentClash,
  rooms,
  subscriptions,
  currentRoom,
  events,
  problems
})

export default rootReducer