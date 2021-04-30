import React from "react";
import { useAuth0 } from "@auth0/auth0-react";
import { createBrowserHistory } from "history";
import { Route, Router, Switch } from "react-router-dom";
import { ProtectedRoute } from "./ProtectedRoute";
import { Nav } from "./Nav";
import { Error } from "./Error";
import { Loading } from "./Loading";
import { Receipts } from "./Receipts";
import { ReceiptDetails } from "./ReceiptDetails";
import { Analytics } from "./Analytics"

// Use `createHashHistory` to use hash routing
export const history = createBrowserHistory();

function App() {
  const { isLoading, error } = useAuth0();

  if (isLoading) {
    return <Loading />;
  }

  return (
    <Router history={history}>
      <Nav />
      {error && <Error message={error.message} />}
      <Switch>
        <Route path="/" exact />
        <ProtectedRoute exact path="/receipts" component={Receipts} />
        <ProtectedRoute exact path="/receipts/:ID" component={ReceiptDetails} />
        <ProtectedRoute exact path="/analytics" component={Analytics} />
      </Switch>
    </Router>
  );
}

export default App;
