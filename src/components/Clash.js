import React, { Component, PropTypes } from 'react'

const Clash = ({text, setInput}) => {

  return (
		<div>
			
			<textarea 
				rows="8"
				onChange = {setInput}
        value = {text}
			></textarea>
		</div>
  )

}

export default Clash
