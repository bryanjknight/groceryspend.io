/* eslint-disable @typescript-eslint/no-explicit-any */
import React from "react";
import { withAuthenticationRequired } from "@auth0/auth0-react";
import { Route } from "react-router-dom";

// eslint-disable-next-line @typescript-eslint/explicit-module-boundary-types
export const ProtectedRoute = ({
  component,
  ...args
}: React.PropsWithChildren<any>) => (
  <Route
    render={(props) => {
      const Component = withAuthenticationRequired(component, {
        // If using a Hash Router, you need to pass the hash fragment as `returnTo`
        // returnTo: () => window.location.hash.substr(1),
      });
      return <Component {...props} />;
    }}
    {...args}
  />
);
