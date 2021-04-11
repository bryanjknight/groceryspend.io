import React from "react";
import ReactDOM from "react-dom";
import "./popup.css";
import { App } from "./App";
import { AppState, Auth0Provider } from "@auth0/auth0-react";

import { getBrowserInstance } from "./browser";

const getRedirectURL = () => {
  return `${getBrowserInstance().identity.getRedirectURL()}`;
};

ReactDOM.render(
  <Auth0Provider
    domain="groceryspend-dev.us.auth0.com"
    clientId="tonoXWFW9VLF9FHkzNxiUULKtibDkTuf"
    redirectUri={getRedirectURL()}
    onRedirectCallback={(appState: AppState) => {
      alert("On Redirect Callback" + JSON.stringify(appState));
    }}
  >
    <React.StrictMode>
      <App webhookUrl="http://localhost:8080/receipts/receipt" />
    </React.StrictMode>
  </Auth0Provider>,
  document.getElementById("root")
);
