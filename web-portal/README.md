# Web

Interacting with server through a web interface.

## Tools to install
* `npm install -g yarn`
* React Dev Extension for Chrome

## Setup
* `yarn` to install all packages
* `yarn watch` to update webpack and run on `localhost:3000`

## Known hacks
### Storybook 6 and Webpack 5 don't play well nice together (yet)
In order to get it to work, I had to do the following:
* `npx sb@next init --builder webpack5` to fresh install storybook
* Update `.storybook/main.js` to leverage the existing webpack5 config