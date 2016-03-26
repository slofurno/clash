import React, { Component, PropTypes } from 'react'

const Lobby = ({room, events}) => {
  
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
    <div className = "card" style = {{
      width: "100%"
    }}>
      <h2> {room.name} </h2>
      <ul>
        { Object.keys(here).map((x, i) => <div className = "user" key = {i}>{x}</div>) }
      </ul> 
    </div>
  ) 
}

export default Lobby
