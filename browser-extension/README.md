# Browser Extesnion

The browswer extension contains the following components:

- A background script, which is used to access the DOM of the active tab
- A popup script, which is used to receive user input (e.g. press a button) and execute actions

## Testing

Things to consider when testing:

- Changes to the `content.ts` file require both the plugin **AND** the active tab to be refreshed

### Chrome

### Firefox

### Safari (To be implemented)

### Edge (To be implemented)

# Resources

- [Developing a browser extension with create react app](https://mmazzarolo.medium.com/developing-a-browser-extension-with-create-react-app-b0dcd3b32b3f)
- [Chrome Extension boilerplate](https://github.com/sivertschou/react-typescript-chrome-extension-boilerplate)
