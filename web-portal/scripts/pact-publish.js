/* eslint-disable @typescript-eslint/no-var-requires */
// if (process.env.CI !== "true") {
//   console.log("skipping Pact publish as not on CI...")
//   process.exit(0)
// }
require('dotenv').config()

const { Publisher } = require("@pact-foundation/pact")
const { versionFromGitTag } = require("@pact-foundation/absolute-version")
const path = require("path")
const childProcess = require("child_process")

const exec = command =>
  childProcess
    .execSync(command)
    .toString()
    .trim()

const branch =
  process.env.TRAVIS_BRANCH || exec("git rev-parse --abbrev-ref HEAD")
// Your version numbers need to be unique for every different version of your consumer
// see https://docs.pact.io/getting_started/versioning_in_the_pact_broker/ for details.
// If you use git tags, then you can use @pact-foundation/absolute-version as we do here.
const consumerVersion = versionFromGitTag()

const opts = {
  pactFilesOrDirs: [path.resolve(process.cwd(), "pact", "pacts")],
  pactBroker: process.env.PACT_BROKER_BASE_URL,
  pactBrokerToken: process.env.PACT_BROKER_TOKEN,
  tags: [branch],
  consumerVersion,
}

new Publisher(opts)
  .publishPacts()
  .then(() => {
    console.log(
      `Pact contract for consumer version ${consumerVersion} published!`
    )
    console.log("")
    console.log("Head over to https://groceryspend.pactflow.io")
    console.log("to see your published contracts.")
  })
  .catch(e => {
    console.log("Pact contract publishing failed: ", e)
  })