import React from 'react';
import ReactDOM from 'react-dom';
import App, { history } from './App';
import { Auth0Provider, AppState } from '@auth0/auth0-react';
import 'bootstrap/dist/css/bootstrap.min.css';



const onRedirectCallback = (appState: AppState) => {
  // If using a Hash Router, you need to use window.history.replaceState to
  // remove the `code` and `state` query parameters from the callback url.
  // window.history.replaceState({}, document.title, window.location.pathname);
  history.replace((appState && appState.returnTo) || window.location.pathname);
};

const REACT_APP_DOMAIN = 'groceryspend-dev.us.auth0.com';
const REACT_APP_CLIENT_ID = 'tonoXWFW9VLF9FHkzNxiUULKtibDkTuf';
const REACT_APP_AUDIENCE = 'https://bknight.dev.groceryspend.io';

ReactDOM.render(
  <React.StrictMode>
    <Auth0Provider
      domain={REACT_APP_DOMAIN}
      clientId={REACT_APP_CLIENT_ID}
      audience={REACT_APP_AUDIENCE}
      redirectUri={window.location.origin}
      scope="read:current_user update:current_user_metadata"
      onRedirectCallback={onRedirectCallback}
    >
      <App />
    </Auth0Provider>
  </React.StrictMode>,
  document.getElementById('root')
);