API design
===

There are three types of models we are trying to tie together:
* The web model, which represents what the user sees on various screens
* The business model, which represents how we model the data in a usable manner
* The database model, which represents how we store the data effectively

Our API is consumed by the browser features (the extension, the web portal); thus, we should priortize the web and business models aligning more closely than the database model.

* The database model should only care about how the data is store and have no knowledge of the use of the data
* The repo should offer different ways to represent the data, specifically in a lightweight (only the bare minimum data) and a heavyweight (all the details)
* The business and web models should be relatively the same

The API is RESTful, so we should think about:

POST /api/v1/requests - submit a receipt for processing
GET /api/v1/requests - get all parse requests, support pagination
GET /api/v1/requests/<uuid> - get details on a specific request
GET /api/v1/receipts - get all receipts, perhaps support pagination
GET /api/v1/receipts/<uuid> - get details on a parsed receipt
PATCH /api/v1/receipts/<uuid>/item/<uuid> - Update a category

The next part is understanding what the schema of the responses are so that the consumers can expect fields without issue. The plan is to use JsonSchema or Pact.io as our contract testing framework. The flow would be:
1. Developer creates, updates, or changes API endpoints on server
1. Server build process recreates the jsonschema file, triggers rebuild of dependent projects (i.e. browser-extension and web-portal)
1. If build of dependencies break, fail the server build