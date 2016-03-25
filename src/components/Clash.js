import React, { Component, PropTypes } from 'react'

const Clash = ({setInput, clash, postCode}) => {
  //TODO: fix this
  let { text, input } = clash.problem || {}
  let { value } = clash
  let onSubmit = (e) => postCode(clash.id, {code:value, runner:"js"})

  return (
		<div>
      <p>{ text }</p>
      <pre>{ input }</pre>
			<textarea 
				rows="8"
				onChange = {setInput}
        value = {value}
			></textarea>
      <input type="button" value="submit" onClick = {onSubmit} />
		</div>
  )

}

export default Clash
