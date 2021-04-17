// use-api.js
import { useEffect, useState } from "react";
import { GetTokenSilentlyOptions, useAuth0 } from "@auth0/auth0-react";
import { FetchOptions } from "@auth0/auth0-spa-js";

interface AuthState {
  error: Error | null;
  loading: boolean;
  data: Record<string, unknown> | null;
  // refresh: () => void;
}

export interface UseApiOptions extends GetTokenSilentlyOptions, FetchOptions {}

export const useApi = (url: string, options: UseApiOptions = {}): AuthState => {
  const { getAccessTokenSilently } = useAuth0();
  const initialState: AuthState = {
    error: null,
    loading: true,
    data: null,
    // eslint-disable-next-line @typescript-eslint/no-empty-function
    // refresh: () => {},
  };

  const [state, setState] = useState(initialState);
  // const [refreshIndex, setRefreshIndex] = useState(0);

  useEffect(() => {
    (async () => {
      try {
        const { audience, scope, ...fetchOptions } = options;
        const accessToken = await getAccessTokenSilently({ audience, scope });
        const res = await fetch(url, {
          ...fetchOptions,
          headers: {
            ...fetchOptions.headers,
            // Add the Authorization header to the existing headers
            Authorization: `Bearer ${accessToken}`,
          },
        });
        setState({
          ...state,
          data: await res.json(),
          error: null,
          loading: false,
        });
      } catch (error) {
        setState({
          ...state,
          error,
          loading: false,
        });
      }
    })();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  return {
    ...state,
    // refresh: () => setRefreshIndex(refreshIndex + 1),
  };
};
