import React, { Component, PropTypes } from 'react'

const Result = ({result, postResult}) => {

  const {code, output, diff, status} = result

  let color = status === 0 ? "green" : "red"
  let onClick = e => postResult(result.clash, result.id)

  return (
    <div>
      <input type="button" value="submit" onClick={onClick}/>
      <pre>{code}</pre>
      <pre>{output}</pre>
      <pre style = {{
        backgroundColor: color, 
      }}>{diff}</pre>
    </div> 
  )
}

const Clash = ({users, setInput, clash, postCode, results, postResult, visibleProblem}) => {
  let { text, input } = visibleProblem 
  let { value } = clash
  let onSubmit = (e) => postCode(clash.id, {code:value, runner:"js"})

  let style = {
    display: "inline-block",
    width: "50%",
    padding: "5px",
    verticalAlign: "top"
  }

  let currentResults = results.filter(x => x.clash === clash.id)

  return (
		<div style = {{
      width: "100%",
      height: "100%",
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
        { currentResults.slice().reverse().map((x, i) => 
            <Result key = {i} result = {x} postResult={postResult}/>
          )
        }
      </div>

		</div>
  )

}

export default Clash
