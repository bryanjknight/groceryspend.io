import { getBrowserInstance } from "./browser";
import { backgroundAuth } from "./auth/login";

const browser = getBrowserInstance();

// handle events when extension is installed and/or updated
browser.runtime.onInstalled.addListener((details) => {
  if (details.reason === "update") {
    // do something if the extension was updated
  }
});

export interface BackgroundMessageRequest {
  action: string;
}

// handle events when messages are sent to the background script
browser.runtime.onMessage.addListener(
  (message: BackgroundMessageRequest, sender, sendResponse) => {
    if (message.action === "AUTH") {
      backgroundAuth(false);
    }
  }
);
