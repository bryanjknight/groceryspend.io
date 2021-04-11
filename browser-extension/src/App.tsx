import "./App.css";
import { Button } from "./components/Button";
import axios from "axios";
import {
  getBrowserInstance,
  EXTRACT_DOM_ACTION,
  ExtractDomResponse,
} from "./browser";
import { useAuth0 } from "@auth0/auth0-react";

const browser = getBrowserInstance();
export interface AppProps {
  webhookUrl: string;
}

export const App = (props: AppProps) => {
  const {
    isLoading,
    isAuthenticated,
    error,
    user,
    loginWithPopup,
    logout,
  } = useAuth0();

  if (isLoading) {
    return <div>Loading...</div>;
  }
  if (error) {
    return <div>Oops... {error.message}</div>;
  }

  const getAxios = (webhookUrl: string) => {
    return axios.create({
      baseURL: webhookUrl,
      // TODO: http proxy?
    });
  };

  const handleContentScriptResponse = (
    extractDomRequest: ExtractDomResponse
  ) => {
    const payload = {
      url: extractDomRequest.url,
      timestamp: new Date().toISOString(),
      data: `<html><body>${extractDomRequest.dom}</body></html>`,
    };
    const httpClient = getAxios(props.webhookUrl);

    httpClient
      .post("/", payload)
      .then((resp) => {
        alert("Success");
      })
      .catch((err) => {
        alert(`Failure: ${err}`);
      });
  };

  const handleSendButtonClick = () => {
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

  return <button onClick={() => loginWithPopup()}>Log in</button>;
};
