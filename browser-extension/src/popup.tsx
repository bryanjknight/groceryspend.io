import React from "react";
import ReactDOM from "react-dom";
import { App } from "./components/App";

ReactDOM.render(
  <React.StrictMode>
    <App webhookUrl={`${process.env.API_URL}/receipts/receipt`} />
  </React.StrictMode>,

  document.getElementById("root")
);
