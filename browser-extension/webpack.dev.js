/* eslint-disable @typescript-eslint/no-var-requires */
const { merge } = require('webpack-merge');
const common = require('./webpack.common.js');

const { DefinePlugin } = require("webpack");

// TODO: better externalization of env vars
const env = {
  REACT_APP_DOMAIN: 'https://groceryspend-dev.us.auth0.com',
  REACT_APP_CLIENT_ID: 'tonoXWFW9VLF9FHkzNxiUULKtibDkTuf',
  REACT_APP_AUDIENCE: 'https://bknight.dev.groceryspend.io',
  API_URL: 'http://localhost:8080'
}

// reduce it to a nice object, the same as before (but with the variables from the file)
const envKeys = Object.keys(env).reduce((prev, next) => {
  prev[`process.env.${next}`] = JSON.stringify(env[next]);
  return prev;
}, {});

module.exports = merge(common, {
  mode: 'development',
  devtool: 'inline-source-map',
  // devtool: "source-map",

  devServer: {
    contentBase: "./dist",
  },
  resolve: {
    alias: {
      "react-dom": "@hot-loader/react-dom",
    },
  },
  plugins: [
    new DefinePlugin(envKeys)
  ]
});