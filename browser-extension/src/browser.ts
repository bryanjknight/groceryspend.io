export const EXTRACT_DOM_ACTION = "groceryspend.get_dom";

export interface ExtractDomResponse {
  dom: string;
  url: string;
}

/* eslint-disable @typescript-eslint/no-explicit-any */
export function getBrowserInstance(): typeof chrome {
  // Get extension api Chrome or Firefox
  const browserInstance = window.chrome || (window as any)["browser"];
  return browserInstance;
}

export function isFirefox(): boolean {
  // TODO: implement this
  return false;
}

export const createTab = async (
  tabOptions: chrome.tabs.CreateProperties
): Promise<chrome.tabs.Tab> =>
  new Promise((resolve) =>
    getBrowserInstance().tabs.create(tabOptions, (tab) => resolve(tab))
  );

export const updateTab = (
  tabId: number,
  updateProps: chrome.tabs.UpdateProperties
): void => getBrowserInstance().tabs.update(tabId, updateProps);

export const closeTab = async (tabId: number): Promise<void> =>
  new Promise((resolve) => getBrowserInstance().tabs.remove(tabId, resolve));

export const notification = (message: string | undefined): Promise<string> =>
  new Promise((resolve) =>
    getBrowserInstance().notifications.create(
      {
        message,
      },
      (notificationId) => resolve(notificationId)
    )
  );
