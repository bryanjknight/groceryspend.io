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