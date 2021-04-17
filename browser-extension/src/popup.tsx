import React from "react";
import ReactDOM from "react-dom";
import "./popup.css";
import { App } from "./App";
import { AppState, Auth0Provider } from "@auth0/auth0-react";

import { getBrowserInstance } from "./browser";

const getRedirectURL = () => {
  // slash is already included in redirect url
  return `${getBrowserInstance().identity.getRedirectURL()}popup.html`;
};

ReactDOM.render(
  <React.StrictMode>
    <Auth0Provider
      domain="groceryspend-dev.us.auth0.com"
      clientId="tonoXWFW9VLF9FHkzNxiUULKtibDkTuf"
      redirectUri={getRedirectURL()}
      skipRedirectCallback={false}
      onRedirectCallback={(appState: AppState) => {
        alert("On Redirect Callback" + JSON.stringify(appState));
      }}
    >
      <App webhookUrl="http://localhost:8080/receipts/receipt" />
    </Auth0Provider>
  </React.StrictMode>,

  document.getElementById("root")
);
