import React from "react";
import ReactDOM from "react-dom";
import { App } from "./components/App";

ReactDOM.render(
  <React.StrictMode>
    <App webhookUrl="http://localhost:8080/receipts/receipt" />
  </React.StrictMode>,

  document.getElementById("root")
);
