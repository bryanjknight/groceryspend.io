# Browser Extension

The browswer extension contains the following components:

- A background script, which is used to make 3rd party calls to auth0 as well as our GrocerySpend.io server
- A content script which is used to extract the DOM of the active tab
- A popup script, which is used to receive user input (e.g. press a button) and execute actions

## Notes on OAuth2

OAuth2 and browser extensions do not have a well implemented solution. This is because `laundWebAuthFlow` assumes only one exchange of information between the auth server and the client. The new guidance is to use PKCE, which makes two requests (one for the code challenge, one to get the access token). Due to this limiation, the OAuth2 workflow had to be implemented manually as opposed to a framework. This is further shown with the deprecation of the [auth0-chrome example repo](https://github.com/auth0-community/auth0-chrome) as well as [feedback from Auth0 of no longer supporting the repo](https://community.auth0.com/t/chrome-extension-advice/38887/9)

## Testing

Things to consider when testing:

- Changes to the `content.ts` file require both the plugin \*_AND_- the active tab to be refreshed

### Chrome

Options include:

- Using the inspect popup widget (insert picture)
- go to `chrome-extension://gpmoghmaibomfddfbofkionknjjeoaef/popup.html` to render the plugin popup
- see https://developer.chrome.com/docs/extensions/mv3/tut_debugging/

### Firefox

### Safari (To be implemented)

### Edge (To be implemented)

# Resources

- [Developing a browser extension with create react app](https://mmazzarolo.medium.com/developing-a-browser-extension-with-create-react-app-b0dcd3b32b3f)
- [Chrome Extension boilerplate](https://github.com/sivertschou/react-typescript-chrome-extension-boilerplate)
