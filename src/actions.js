import request from './request'

export const SET_INPUT = 'SET_INPUT'
export const ADD_ROOMS = 'ADD_ROOMS'
export const JOINED_ROOM = 'JOINED_ROOM'

let host = location.host
let baseurl = `http://${host}/api`
let origin = location.origin.replace(/^http/, "ws")
let sock = new WebSocket(origin + "/api/ws")

sock.onmessage = function (e) {
  console.log(e)
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

function getRoomsSuccess (rooms) {
  return {
    type: ADD_ROOMS,
    rooms
  }
}

export function joinRoom(room) {
  subscribe(room.id)  
  return {
    type: JOINED_ROOM,
    room
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

export function getClash (id) {
  return function (dispatch) {
    
    return request({
      method: "GET",
      url: `${baseurl}/clash/${id}`
    })
    .then(parse)
    .then(x => {
      
    })
  }
}
