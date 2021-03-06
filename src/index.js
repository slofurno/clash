import React from 'react'
import { createStore, applyMiddleware} from 'redux'
import { Provider } from 'react-redux'
import thunk from 'redux-thunk'
import { render } from 'react-dom'

import rootReducer from './reducers'
import { 
  getRooms,
  getProblems,
  dial,
  setToken,
  createAccount
} from './actions'
import App from './components/App'

let store = createStore(
  rootReducer,
  applyMiddleware(thunk)
)

let unsubscribe = store.subscribe(() => 
  console.log(store.getState())
)

store.dispatch(getRooms())
store.dispatch(getProblems())
dial(store.dispatch)

let token = localStorage.getItem("Clash.token")
if (!token) {
  store.dispatch(createAccount({email:"asdf", password:"asdf"}))
} else {
  store.dispatch(setToken(token))
}

render(
  <Provider store={store}>
    <App/>
  </Provider>,
  document.getElementById('root')
)

