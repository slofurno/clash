import request from './request'

export const SET_INPUT = 'SET_INPUT'
export const ADD_ROOMS = 'ADD_ROOMS'
export const JOINED_ROOM = 'JOINED_ROOM'
export const ADD_EVENTS = 'ADD_EVENTS'
export const ADD_PROBLEMS = 'ADD_PROBLEMS'
export const SHOW_CLASH = 'SHOW_CLASH'
export const JOIN_CLASH ='JOIN_CLASH'

let host = location.host
let baseurl = `http://${host}/api`
let origin = location.origin.replace(/^http/, "ws")
let sock = null

export function dial (dispatch) {
  sock = new WebSocket(origin + "/api/ws")

  sock.onmessage = function (e) {
    let d = JSON.parse(e.data)

      console.log(d)
    switch (d.verb) {
    case "JOINED_LOBBY":
      dispatch(addEvents([d]))
    case "STARTED_CLASH":
      dispatch(getClash(d.noun))
    case "ran":
      dispatch(getResult(d.subject))
    default:
      console.log(d)
    }
  }
}


export function setInput (e) {
  return {
    type: SET_INPUT,
    value: e.target.value
  }
}

function parse (x) {
  return JSON.parse(x)
}

function error (x) {
  console.error(x)  
}

function getProblemSuccess (problems) {
  return {
    type: ADD_PROBLEMS,
    problems
  }
}

function getRoomsSuccess (rooms) {
  return {
    type: ADD_ROOMS,
    rooms
  }
}

function joinRoomSuccess (room) {
  return {
    type: JOINED_ROOM,
    room
  }
}

function addEvents (events) {
  return {
    type: ADD_EVENTS,
    events
  }
}

function showClash (problem) {
  return {
    type: SHOW_CLASH,
    problem 
  }
}

function joinClash (id) {
  return {
    type: JOIN_CLASH,
    id
  }
}

function waitForResult (id) {
  //
  subscribe(id) 
}

function getResult (id) {
  return function (dispatch) {
    return request({
      method: "GET",
      url: `/api/code/${id}`
    })
    .then(parse)
    .then(x => console.log(x))
    .catch(error)
  }
}

export function joinRoom(room) {
  return function (dispatch) {
    subscribe(room.id)  
    return request({
      method: "GET",
      url: `/api/events/${room.id}`
    })
    .then(parse)
    .then(x => {
      dispatch(addEvents(x))
      dispatch(joinRoomSuccess(room))
    })
    .catch(error)
  }
}

export function getRooms() {
  return function (dispatch) {
    return request({
      method: "GET",
      url: `${baseurl}/rooms`
    })
    .then(parse)
    .then(x => dispatch(getRoomsSuccess(x)))
    .catch(error)
  }
}

function unsubscribe (topic) {
  sock.send(JSON.stringify({
    type: "UNSUB",
    subject: topic 
  }))
}

function subscribe (topic) {
  sock.send(JSON.stringify({
    type: "SUB",
    subject: topic 
  }))
}

function getProblem (problem) {
  return function (dispatch) {
    return request({
      method: "GET",
      url: `/api/problems/${problem}`
    })
    .then(parse)
    .then(x => dispatch(showClash(x)))
    .catch(error)
  }
}

export function postCode (clash, code) {
  return function (dispatch, getState) {
    return request({
      method: "POST",
      url: `/api/clash/${clash}`,
      body: JSON.stringify(code),
      headers: {Authorization:"9fd80e71-c5a7-427d-a80b-1ccee4ed5c9e"}
    })
    .then(parse)
    .then(x => {
      waitForResult(x.id)  
    })
    .catch(error)
  }
}

export function getClash (id) {
  return function (dispatch) {
    dispatch(joinClash(id))
    return request({
      method: "GET",
      url: `${baseurl}/clash/${id}`
    })
    .then(parse)
    .then(x => { 
      dispatch(getProblem(x.problem))
    })
  }
}

export function getProblems () {
  return function (dispatch) {
    return request({
      method: "GET",
      url: "/api/problems"
    })
    .then(parse)
    .then(x => dispatch(getProblemSuccess(x)))
    .catch(error)
  }
}
