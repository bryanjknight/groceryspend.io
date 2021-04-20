// use-api.js
export interface UseApiOptions {
  accessToken: string;
  audience?: string;
  scope?: string;
  headers?: Record<string, string>;
  method: string;
}

export interface ResponseState {
  error?: Error;
  loading: boolean;
  data?: Record<string, unknown>;
}

export const useApi = async (
  url: string,
  options: UseApiOptions,
  data: BodyInit | null | undefined
): Promise<ResponseState> => {
  try {
    const { audience, scope, accessToken, ...fetchOptions } = options;
    const res = await fetch(url, {
      ...fetchOptions,
      headers: {
        ...fetchOptions.headers,
        // Add the Authorization header to the existing headers
        Authorization: `Bearer ${accessToken}`,
      },
      method: options.method,
      body: data,
    });
    return {
      data: await res.json(),
      loading: false,
    };
  } catch (error) {
    return {
      error: error,
      loading: false,
    };
  }
};
