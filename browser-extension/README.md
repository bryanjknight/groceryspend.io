# Browser Extension

The browswer extension contains the following components:

- A background script, which is used to access the DOM of the active tab
- A popup script, which is used to receive user input (e.g. press a button) and execute actions

## Notes on Auth0

Auth0 and extensions don't seem to play well together. Things to consider:

- https://github.com/auth0-community/auth0-chrome is deprecated but the best option. Look at https://github.com/auth0-community/auth0-chrome/blob/master/src/PKCEClient.js#L63-L80

- Reading material:
  - https://github.com/an-object-is-a/chrome-ext-discord-oauth2/blob/master/background.js
  - https://github.com/michaeloryl/oauth2-angularjs-chrome-extension-demo

Takeaway:

- I need to register the extension in both Chrome and Firefox to get their respective browser extension IDs.
  - Chrome Extension id:
    - Web Store: gphmemfooelbfnnlnjegjofkkabhebek
  - Firefox browser add-on: ccc07083ce20415a94bd/{c5de945e-999e-4b2b-ab30-a7d4663c4058}
- Once I have the browser extension IDs, I can then register them in Auth0 as a valid callback
  - The production and development versions will have different IDs. We'll need to add the appropriate ones for the given application/tenancy
  - https://gphmemfooelbfnnlnjegjofkkabhebek.chromiumapp.org/
  - Add `key` to `manifest.json` to override key in dev mode
- From there, I can implement the `launchWebFlow` (which creates a separate tab for authentication)
- Once the webflow is complete, you can then get the access token
  - Need to figure out how to leverage the `launchWebFlow` callback to work with the request coming from Auth0

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
