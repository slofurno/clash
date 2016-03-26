import React, { Component, PropTypes } from 'react'

const Lobby = ({room, events, users}) => {
  
  let xs = events.filter(x => x.subject === room.id).sort((a,b) => a.time - b.time)
  let here = {}

  xs.forEach(x => {
    if (x.verb === "JOINED_LOBBY") {
      here[x.noun] = true
    } else {
      delete here[x.noun]
    }
  })

  return (
    <div style = {{
      width: "100%",
      height: "100%"
    }}>
      <h2> {room.name} </h2>
      <ul>
        { Object.keys(here).map((x, i) => <div className = "user" key = {i}>{users[x]}</div>) }
      </ul> 
    </div>
  ) 
}

export default Lobby
