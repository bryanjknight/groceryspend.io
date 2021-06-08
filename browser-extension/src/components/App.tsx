import { Button } from "./Button";
import {
  getBrowserInstance,
  EXTRACT_DOM_ACTION,
  ExtractDomResponse,
  fireMessage,
} from "../lib/browser";
import { UserInfo } from "../lib/auth";
import { useEffect, useState } from "react";
import { BackgroundMessageRequest } from "../background";
import { ParseReceiptRequest } from "../models";

const browser = getBrowserInstance();

export interface AppProps {
  webhookUrl: string;
}

export const App = (props: AppProps): JSX.Element => {
  const [loaded, setLoaded] = useState(false);
  const [isAuthenticated, setIsAuthenticated] = useState(false);
  const [userInfo] = useState({} as UserInfo);

  useEffect(() => {
    fireMessage<BackgroundMessageRequest, boolean>({ action: "VERIFY" }).then(
      (resp) => {
        setIsAuthenticated(resp);
        setLoaded(true);
      }
    );
  }, []);

  const handleSendButtonClick = () => {
    const handleContentScriptResponse = (resp: ExtractDomResponse) => {
      fireMessage<BackgroundMessageRequest, string>({
        action: "SEND",
        data: new ParseReceiptRequest({
          url: resp.url,
          timestamp: new Date(),
          data: resp.dom,
          // we only support html requests from the browser extension
          parseType: 1,
        }),
      }).then((resp) => {
        alert(resp);
      });
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

  const handleLogInClick = async () => {
    setLoaded(false); // show that the popup is reloading

    fireMessage<BackgroundMessageRequest, boolean>({ action: "LOGIN" }).then(
      (resp) => {
        alert(resp);
        setIsAuthenticated(resp);
        setLoaded(true);
      }
    );
  };

  const handleLogOutClick = async () => {
    setLoaded(false); // show that the popup is reloading

    fireMessage<BackgroundMessageRequest, boolean>({ action: "LOGOUT" }).then(
      () => {
        setIsAuthenticated(false);
        setLoaded(true);
      }
    );
  };

  if (!loaded) {
    return <div>Loading...</div>;
  }

  if (isAuthenticated) {
    return (
      <div className="App">
        <p>
          Hello {userInfo.name}{" "}
          <button onClick={() => handleLogOutClick()}>Log out</button>
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
