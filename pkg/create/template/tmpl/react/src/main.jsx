import React from 'react'
import {createRoot} from 'react-dom/client'
import './style.css'
import App from './App'

const container = document.getElementById('root')

const root = createRoot(container)

root.render(
    // strict mode makes the component App call twice
    // <React.StrictMode>
        <App/>
    // </React.StrictMode>
)
