import React, { Component, PropTypes } from 'react'

const ClashResults = ({results}) => {
  return (
    <div>
      {results.map(x => <pre key={x.id}>{x.code}</pre>)}
    </div>
  )
}

export default ClashResults
