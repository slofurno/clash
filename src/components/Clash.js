import React, { Component, PropTypes } from 'react'

const Result = ({result}) => {

  const {code, output, diff, status} = result

  let color = status === 0 ? "green" : "red"

  return (
    <div>
      <pre>{code}</pre>
      <pre>{output}</pre>
      <pre style = {{
        backgroundColor: color, 
      }}>{diff}</pre>
    </div> 
  )
}

const Clash = ({setInput, clash, postCode, results}) => {
  //TODO: fix this
  console.log(clash)
  let { text, input } = clash.problem || {}
  let { value } = clash
  let onSubmit = (e) => postCode(clash.id, {code:value, runner:"js"})

  let style = {
    display: "inline-block",
    width: "50%",
    padding: "5px",
    verticalAlign: "top"
  }

  return (
		<div style = {{
      width: "800px",
      height: "400px",
      padding: "6px",
      backgroundColor: "silver",
    }}>
      <div style = {style}>
        <p>{ text }</p>
        <pre>{ input }</pre>
        <textarea 
          rows="16"
          onChange = {setInput}
          value = {value}
        ></textarea>
        <input type="button" value="submit" onClick = {onSubmit} />
      </div>
      <div style = {style}>
        { results.slice().reverse().map((x, i) => <Result key = {i} result = {x}/>) }
      </div>

		</div>
  )

}

export default Clash
