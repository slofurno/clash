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

  console.log(here)

  return (
    <div style = {{
      backgroundColor: "lightslategray"
    }}>
      <h2> {room.name} </h2>
      <ul>
        { Object.keys(here).map((x, i) => <li key = {i}> {x} </li>) }
      </ul> 
    </div>
  ) 
}

export default Lobby
