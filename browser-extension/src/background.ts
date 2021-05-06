import { getBrowserInstance } from "./lib/browser";
import { login, verifyToken } from "./lib/auth";
import { createAES256LocalStorage } from "./lib/storage";
import { useApi } from "./lib/api";

const browser = getBrowserInstance();

// handle events when extension is installed and/or updated
browser.runtime.onInstalled.addListener((details) => {
  if (details.reason === "update") {
    // do something if the extension was updated
  }
});

export interface BackgroundMessageRequest {
  action: string;
  data?: unknown;
}

// TODO: extract from process.env
const ENC_STORAGE_CIPHER = "test-cipher-replace-me-1234";

const AUTH_TOKEN_KEY = "accessToken";

const encStorage = createAES256LocalStorage(browser, ENC_STORAGE_CIPHER);

const setAuthToken = (token: string): Promise<boolean> =>
  new Promise((resolve) => {
    encStorage.setEncryptedValue(AUTH_TOKEN_KEY, token);
    resolve(true);
  });

const getAuthToken = (): Promise<string | null> =>
  new Promise((resolve) =>
    resolve(encStorage.getEncryptedValue(AUTH_TOKEN_KEY))
  );

// handle events when messages are sent to the background script
browser.runtime.onMessage.addListener(
  (message: BackgroundMessageRequest, sender, sendResponse) => {
    if (message.action === "LOGIN") {
      new Promise<string>((resolve) => login(resolve))
        .then(setAuthToken)
        .then((cookie) => {
          if (!cookie) sendResponse(false);
          else sendResponse(true);
        })
        .catch((err) => {
          console.log(err);
          sendResponse(false);
        });
      // tell the listener that this is an async operation
      return true;
    }

    if (message.action === "VERIFY") {
      getAuthToken()
        .then((token) => verifyToken(token))
        .then((isValid) => sendResponse(isValid))
        .catch((err) => {
          console.error(err);
          sendResponse(false);
        });
      // tell the listener that this is an async operation
      return true;
    }

    if (message.action === "LOGOUT") {
      new Promise<void>((resolve) => {
        browser.storage.local.remove([AUTH_TOKEN_KEY], resolve);
      })
        .then(() => {
          sendResponse();
        })
        .catch((err) => {
          console.error(err);
          sendResponse(false);
        });
      // tell the listener that this is an async operation
      return true;
    }

    if (message.action === "SEND") {
      getAuthToken()
        .then(async (token) => {
          if (!token) {
            throw new Error("No bearer token available");
          }

          return useApi(
            `${process.env.API_URL}/receipts/receipt`,
            {
              accessToken: token,
              method: "POST",
              headers: {
                "Content-type": "application/json",
              },
            },
            JSON.stringify(message.data || {})
          );
        })
        .then((resp) => {
          if (resp.error) {
            sendResponse(resp.error.message);
          } else {
            sendResponse("Success");
          }
        })
        .catch((err) => {
          sendResponse(err);
        });
      // tell the listener that this is an async operation
      return true;
    }
  }
);
