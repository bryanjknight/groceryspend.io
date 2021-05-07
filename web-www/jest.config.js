module.exports = {
  preset: 'ts-jest',
  testEnvironment: 'jsdom',
  coverageDirectory: "./.coverage",
  collectCoverageFrom: [
    "**/*.{ts,tsx}",
    "!**/*.test.{ts,tsx}",
    "!**/node_modules/**",
    "!**/vendor/**"
  ],
  // let jest know how to handle css and less files
  moduleNameMapper: {
    "\\.css$": "identity-obj-proxy",
    "\\.less$": "identity-obj-proxy",
  },
};