import { getBrowserInstance, EXTRACT_DOM_ACTION } from "./browser";

const browser = getBrowserInstance();

browser.runtime.onMessage.addListener((message, sender, sendResponse) => {
  if (message.action === EXTRACT_DOM_ACTION) {
    sendResponse({
      dom: document.body.innerHTML,
      url: window.location.href,
    });
  }
});
