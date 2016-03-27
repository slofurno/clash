import React, { Component, PropTypes } from 'react'

      //onclick should actually get the code + show it..
const ClashResults = ({results, users, setCode}) => {
  return (
    <div>
      {results.map((x, i) => <pre key={i} onClick={e => setCode(x.code)}>{users[x.user]}{"  "}{x.time}</pre>)}
    </div>
  )
}

export default ClashResults
