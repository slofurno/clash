import request from './request'

export const SET_INPUT = 'SET_INPUT'
export const ADD_ROOMS = 'ADD_ROOMS'
export const JOINED_ROOM = 'JOINED_ROOM'
export const ADD_EVENTS = 'ADD_EVENTS'
export const ADD_PROBLEMS = 'ADD_PROBLEMS'
export const SHOW_CLASH = 'SHOW_CLASH'
export const JOIN_CLASH ='JOIN_CLASH'
export const ADD_RESULT = 'ADD_RESULT'
export const ADD_USER = 'ADD_USER'
export const ADD_CLASH_RESULT = 'ADD_CLASH_RESULT'
export const SET_TOKEN = 'SET_TOKEN'

let host = location.host
let baseurl = `http://${host}/api`
let origin = location.origin.replace(/^http/, "ws")
let sock = null

let seen = {}

export function dial (dispatch) {
  sock = new WebSocket(origin + "/api/ws")

  sock.onmessage = function (e) {
    let d = JSON.parse(e.data)
    handleEvent(d, dispatch)
  }
}

function handleEvent (d, dispatch) {
  switch (d.verb) {
  case "JOINED_LOBBY":
    dispatch(getUser(d.noun))
    dispatch(addEvents([d]))
    break;
  case "STARTED_CLASH":
    dispatch(getClash(d.noun))
    break;
  case "ran":
    dispatch(getResult(d.subject))
    break;
  case "POSTED_RESULT":
    dispatch(getClashResult(d.noun))
    break;
    console.log("someone posted a result")
  default:
  }
}

export function setSlide (slide) {
  return {
    type: "SET_SLIDE",
    slide
  }
}

function addClashResults (result) {
  return {
    type: "SET_CLASH_RESULT",
    result
  }
}

function addClashResult (result) {
  return {
    type: ADD_CLASH_RESULT,
    result
  }
}

export function setInput (e) {
  return {
    type: SET_INPUT,
    value: e.target.value
  }
}

export function setToken (token) {
  return {
    type: SET_TOKEN,
    token
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

function addResult (result) {
  return {
    type: ADD_RESULT,
    result
  }
}

function setUser (user) {
  return {
    type: ADD_USER,
    user
  }
}

export function getUser (id) {
  return function (dispatch, getState) {
    if (seen[id]) {
      return
    } else {
      seen[id] = true
    }

    if (getState().users[id]) {
      return
    }
    return request({
      method: "GET",
      url: `/api/accounts/${id}` 
    })
    .then(parse)
    .then(x => {
      let m = {}
      m[id] = x.username || "user"
      dispatch(setUser(m))
    })
    .catch(error)    
  }
}

export function postResult (clash, code) {
  return function (dispatch, getState) {
    let token = getState().token || ""
    return request({
      method: "POST",
      url: `/api/clash/${clash}/code/${code}`,
      headers: {Authorization: token}
    })
//    .then(x => dispatch(showResults(clash)))
    .catch(error)
  }
}


function getResult (id) {
  return function (dispatch, getState) {
    let matches = getState().results.filter(x => x.id === id)
    if (matches.length > 0) {
      console.log("already have this result")
      return
    }
    return request({
      method: "GET",
      url: `/api/code/${id}`
    })
    .then(parse)
    .then(x => {
      dispatch(addResult(x))
      dispatch(getUser(x.user))
    })
    .catch(error)
  }
}

export function joinRoom(room) {
  return function (dispatch, getState) {
    const token = getState().token
    subscribe(room.id)  
    return request({
      method: "POST",
      url: `/api/rooms/${room.id}`,
      headers: {Authorization:token}
    })
    .then(parse)
    .then(xs => {
      //this is kind of fd
      //joinroom clears clashes + then we add our rsults
      dispatch(joinRoomSuccess(room))
      xs.map(x => handleEvent(x, dispatch))
    })
    .catch(error)
  }
}

export function getRooms() {
  return function (dispatch) {
    return request({
      method: "GET",
      url: `/api/rooms`
    })
    .then(parse)
    .then(x => dispatch(getRoomsSuccess(x)))
    .catch(error)
  }
}

function addClash (clash) {
  return {
    type: "ADD_CLASH",
    clash
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

//FIXME: have to set both clash + problemid
//
function setProblem2 (problem, clash) {
  return {
    type: "SET_PROBLEM",
    problem,
    clash
  }
}

function addCode (code) {
  return {
    type: "ADD_CODE",
    code
  } 
}

export function setCode (code) {
  return {
    type: "SET_CODE",
    code
  } 
}

//TODO: make this set the code + the clashResult
function getClashResult (id) {
  return function (dispatch, getState) {
    return request({
      method: "GET",
      url: `/api/code/${id}`
    })
    .then(parse)
    .then(x => {
      dispatch(addCode(x))
      dispatch(addClashResult([
        {
          code: x.id,
          user: x.user,
          time: x.time,
          clash: x.clash,
          status: x.status
        }
      ]))
    })
    .catch(error)
  }
}

export function setProblem (problem, clash) {
  return function (dispatch, getState) {
    dispatch(setProblem2(problem,clash))
    return request({
      method: "GET",
      url: `/api/clash/${clash}/code`
    })
    .then(parse)
    .then(x => dispatch(addClashResults(x)))
    .catch(error)
  }
}

export function getCode (id) {
  return function (dispatch, getState) {
    dispatch(setCode(id))
    if (getState().codes[id]) {
      return
    }
    return request({
      method: "GET",
      url: `/api/code/${id}`
    })
    .then(parse)
    .then(x => dispatch(addCode(x)))
    .catch(error)
  }
}

export function postCode (clash, code) {
  return function (dispatch, getState) {
    let token = getState().token || ""
    return request({
      method: "POST",
      url: `/api/clash/${clash}`,
      body: JSON.stringify(code),
      headers: {Authorization:token}
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
    subscribe(id)
    //dispatch(joinClash(id))
    return request({
      method: "GET",
      url: `${baseurl}/clash/${id}`
    })
    .then(parse)
    .then(x => { 
      dispatch(addClash(x))
      //dispatch(getProblem(x.problem))
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

export function createAccount (account) {
  return function (dispatch) {
    return request({
      method: "POST",
      url: "/api/accounts",
      body: JSON.stringify(account)
    })
    .then(parse)
    .then(x => {
      localStorage.setItem("Clash.token", x.token)
      dispatch(setToken(x.token))
    }) 
    .catch(error)
  }
}
