import { getBrowserInstance } from "./browser";
import pkce from "pkce-challenge";
import urlParse from "url-parse";

export const CODE_VERIFIER_KEY = "codeVerifier";

const browser = getBrowserInstance();

const redirectUri = browser.identity.getRedirectURL();

// TODO: extract these from process.env
const domain = "https://groceryspend-dev.us.auth0.com";
const clientId = "tonoXWFW9VLF9FHkzNxiUULKtibDkTuf";
const audience = "https://bknight.dev.groceryspend.io";
const permissions = ["openid", "profile"];

const scope = permissions.join("%20");

const handleAuthorizationCode = (
  codeVerifier: string,
  state: string,
  cb: (accessToken: string) => void
) => (responseUrl: string | undefined) => {
  if (browser.runtime.lastError) {
    throw new Error(browser.runtime.lastError.message);
  }

  if (!responseUrl) {
    throw new Error("Response URL was not defined");
  }

  // step 3 - after the user is redirected back, verify the state
  const { code, state: responseState } = urlParse(responseUrl, true).query;
  if (responseState !== state) {
    throw new Error("Cross-site request forgery attack detected.");
  }

  if (!code) {
    throw new Error("Code is not set");
  }

  const options = new URLSearchParams({
    grant_type: "authorization_code",
    client_id: clientId,
    code_verifier: codeVerifier,
    code: code,
    redirect_uri: redirectUri,
  });

  // step 4 - exchange auth code and code verifier for an access token
  fetch(`${domain}/oauth/token`, {
    method: "POST",
    headers: {
      "Content-type": "application/x-www-form-urlencoded",
    },
    body: options.toString(),
  })
    .then((response) => response.json())
    .then((data) => {
      if (data.error) {
        throw new Error(data.error_description.split(/\r?\n/)[0]);
      } else {
        return data;
      }
    })
    .then((data) => {
      cb(data.access_token);
    });
};

// flow inspired by https://github.com/ukhan/add-to-ms-todo and
// https://www.oauth.com/playground/authorization-code-with-pkce.html
export const login = async (setToken: (token: string) => void) => {
  // step 1 - create a secret code and a code challenge
  const { code_verifier: tmpState } = pkce(43);
  const { code_verifier: codeVerifier, code_challenge: codeChallenge } = pkce(
    50
  );

  // take the first 12 characters
  const state = tmpState.substr(0, 12);

  // step 2 - build the authorization url and have the user go through the login
  const queryParamsObj = {
    client_id: clientId,
    response_type: "code",
    redirect_uri: redirectUri,
    response_mode: "query",
    scope,
    state,
    code_challenge: codeChallenge,
    code_challenge_method: "S256",
    prompt: "login",
    audience: audience,
  };
  const queryParams = new URLSearchParams(queryParamsObj);
  const authURL = `${domain}/authorize?${queryParams.toString()}`;
  console.log(authURL);

  // launchWebAuthFlow is the magic that tells the browser "hey, this weird redirect
  // is actually an extension, so send it to the callback function instead"
  // see https://stackoverflow.com/a/35773982/704525
  browser.identity.launchWebAuthFlow(
    {
      url: authURL,
      interactive: true,
    },
    // step 3 - after the user is redirected back, verify the state
    handleAuthorizationCode(codeVerifier, state, setToken)
  );
};

export const verifyToken = (token: string | null): boolean => {
  if (!token) return false;

  // TODO: implement token verification
  return true;
};

export interface UserInfo {
  name: string;
}

export const getUserProfile = async (token: string): Promise<UserInfo> => {
  return new Promise((resolve) => {
    return {
      name: "TBD",
    };
  });
};

export const logout = (token: string) => {
  // build out logout URL

  // fetch

  // run callback (or make this a promise)
  throw new Error("Do something");
};
