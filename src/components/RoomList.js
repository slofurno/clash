import React, { Component, PropTypes } from 'react'

const Room = ({room, joinRoom}) => {
  return (
    <li onClick = {(e) => joinRoom(room)}>
      {room.name}
    </li>
  ) 
}

const RoomList = ({rooms, joinRoom}) => {
  return (
    <ul>
      {rooms.map((x, i) => <Room key = {i} room = {x} joinRoom = {joinRoom}/>)}
    </ul>
  )
}

export default RoomList
