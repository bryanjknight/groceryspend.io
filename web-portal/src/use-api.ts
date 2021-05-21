/* eslint-disable @typescript-eslint/no-explicit-any */
import { useEffect, useState } from 'react';
import { useAuth0 } from '@auth0/auth0-react';

// useApi wraps around an api call and returns the type desired
export const useApi = <T>(
  apiCall: (bearerToken: string) => Promise<T>,
  options: any = {}
): { error?: Error | null; loading: boolean; data?: T | null } => {
  const { getAccessTokenSilently } = useAuth0();
  const [state, setState] = useState({
    error: null as Error | undefined | null,
    loading: true,
    data: null as T | null,
  });

  useEffect(() => {
    (async () => {
      try {
        const { audience, scope } = options;
        const accessToken = await getAccessTokenSilently({ audience, scope, timeoutInSeconds: 60*60 });
        const data = await apiCall(accessToken);
        setState({
          ...state,
          data: data,
          error: null,
          loading: false,
        });
      } catch (error) {
        // TODO: hook for o11y here (e.g. sentry)
        console.error(`Error reported: ${error}`);
        setState({
          ...state,
          error,
          loading: false,
        });
      }
    })();
  }, []); // eslint-disable-line react-hooks/exhaustive-deps

  return state;
};