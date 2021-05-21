// const custom = require("../webpack.common.js");

module.exports = {
  stories: ["../src/**/*.stories.mdx", "../src/**/*.stories.@(js|jsx|ts|tsx)"],
  addons: ["@storybook/addon-links", "@storybook/addon-essentials"],
  core: {
    builder: "webpack5",
  },
  //
  // Theoritically we could add our webpack config
  // however, not necessary (yet) and it breaks things
  //
  // webpackFinal: (config) => {
  //   return {
  //     ...config,
  //     module: {
  //       ...config.module,
  //       rules: [...config.module.rules, ...custom.module.rules],
  //     },
  //   };
  // },
};
