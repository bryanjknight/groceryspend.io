import React from 'react';
import ReactDOM from 'react-dom';
import './popup.css';
import { App } from './App';

ReactDOM.render(
  <React.StrictMode>
    <App webhookUrl="http://localhost:8080/webhook/receipt"/>
  </React.StrictMode>,
  document.getElementById('root')
);