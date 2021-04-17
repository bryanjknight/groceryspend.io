import "./App.css";
import { Button } from "./components/Button";
import {
  getBrowserInstance,
  EXTRACT_DOM_ACTION,
  ExtractDomResponse,
} from "./browser";
import { useAuth0 } from "@auth0/auth0-react";
import { useApi, UseApiOptions } from "./api";

const browser = getBrowserInstance();
export interface AppProps {
  webhookUrl: string;
}

export const App = (props: AppProps) => {
  const { error, isAuthenticated, isLoading, user, logout } = useAuth0();

  if (isLoading) {
    return <div>Loading...</div>;
  }
  if (error) {
    return <div>Oops... {error.message}</div>;
  }

  const handleSendButtonClick = () => {
    const handleContentScriptResponse = (resp: ExtractDomResponse) => {
      alert(resp);
    };

    // send a message to the content script to extract the dom
    browser.tabs.query({ active: true, currentWindow: true }, (tabs) => {
      if (tabs.length !== 1) {
        console.error(`Expected 1 tab, got ${tabs.length}`);
        return;
      }
      const tab = tabs[0];
      if (!tab.id) {
        console.error(`Received null tab id`);
        return;
      }
      browser.tabs.sendMessage(
        tab.id,
        { action: EXTRACT_DOM_ACTION },
        handleContentScriptResponse
      );
    });
  };

  const handleLogInClick = () => {
    browser.runtime.sendMessage({ action: "AUTH" }, (response) => {
      console.log(response);
    });
  };

  if (isAuthenticated) {
    return (
      <div className="App">
        <p>
          Hello {user.name}{" "}
          <button onClick={() => logout({ returnTo: window.location.origin })}>
            Log out
          </button>
        </p>
        <header className="App-header">
          <Button
            onClick={handleSendButtonClick}
            text={"Send to GrocerySpend.io"}
          />
        </header>
      </div>
    );
  }
  return <button onClick={() => handleLogInClick()}>Log In</button>;
};
