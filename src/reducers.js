import { combineReducers } from 'redux'
import {
  SET_INPUT,
  ADD_ROOMS,
  JOINED_ROOM,
  ADD_EVENTS,
  ADD_PROBLEMS,
  SHOW_CLASH,
  JOIN_CLASH,
  SET_TOKEN,
  ADD_RESULT,
  ADD_USER,
  ADD_CLASH_RESULT
} from './actions'

function users (state = {}, action) {
  switch (action.type) {
  case ADD_USER:
    return Object.assign({}, state, action.user)
  default:
    return state
  }
}

function token (state = null, action) {
  switch (action.type) {
  case SET_TOKEN:
    return action.token
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
  case "SET_PROBLEM":
    return Object.assign({}, state, {id: action.clash})
  case JOIN_CLASH:
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

function problems (state = {}, action) {
  switch (action.type) {
  case ADD_PROBLEMS:
    let p = {}
    action.problems.forEach(x => p[x.id] = x)
    return p 
    //return state.concat(action.problems)
  default:
    return state
  }
}

function results (state = [], action) {
  switch (action.type) {
  case ADD_RESULT:
    return state.concat([action.result])
  default:
    return state
  }
}

function clashResults (state = [], action) {
  switch (action.type) {
  case "SET_CLASH_RESULT":
    return action.result
  case ADD_CLASH_RESULT:
    return state.concat(action.result)
  default:
    return state
  }
}

function clashes (state = [], action) {
  switch (action.type) {
  case JOINED_ROOM:
    return []
  case "ADD_CLASH":
    return state.concat([action.clash])
  default:
    return state
  }
}

function slide (state = "LOBBY", action) {
  switch (action.type) {
  case "SET_SLIDE":
    return action.slide
  default:
    return state
  }
}

function currentProblem (state = "", action) {
  switch (action.type) {
  case "SET_PROBLEM":
    return action.problem
  default:
    return state
  }
}

function codes (state = {}, action) {
  switch (action.type) {
  case "ADD_CODE":
    let c = {}
    c[action.code.id] = action.code
    return Object.assign({}, state, c)
  default:
    return state
  }
}

function currentCode (state = "", action) {
  switch (action.type) {
  case "SET_CODE":
    return action.code
  default:
    return state
  }
}


const rootReducer = combineReducers({
  users,
  currentClash,
  rooms,
  subscriptions,
  currentRoom,
  events,
  problems,
  token,
  results,
  clashResults,
  clashes,
  slide,
  currentProblem,
  currentCode,
  codes
})

export default rootReducer
