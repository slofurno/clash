import React, { Component, PropTypes } from 'react'

const Result = ({result}) => {

  const {code, output, diff, status} = result

  let color = status === "0" ? "green" : "red"

  return (
    <div>
      <div>{code}</div>
      <div>{output}</div>
      <div style = {{
        backgroundColor: color, 
      }}>{diff}</div>
    </div> 
  )
}

const Clash = ({setInput, clash, postCode, results}) => {
  //TODO: fix this
  console.log(clash)
  let { text, input } = clash.problem || {}
  let { value } = clash
  let onSubmit = (e) => postCode(clash.id, {code:value, runner:"js"})

  return (
		<div style = {{
      width: "500px",
      height: "400px",
      padding: "20px",
      backgroundColor: "silver",
    }}>
      <p>{ text }</p>
      <pre>{ input }</pre>
			<textarea 
				rows="8"
				onChange = {setInput}
        value = {value}
			></textarea>
      <input type="button" value="submit" onClick = {onSubmit} />
   
      { results.map((x, i) => <Result key = {i} result = {x}/>) }
      

		</div>
  )

}

export default Clash
