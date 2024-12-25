interface Config {
  apiUrl: string;
}

function getApiUrl(): string {
  if (import.meta.env.DEV) {
    return import.meta.env.PUBLIC_API_URL || 'http://localhost:8090';
  }
  
  // In production, use the environment variable or fall back to /api
  const apiHost = import.meta.env.PUBLIC_API_HOST || 'api.localhavencms.com';
  const apiPort = import.meta.env.PUBLIC_API_PORT;
  
  return apiPort ? `http://${apiHost}:${apiPort}` : `http://${apiHost}`;
}

export const config: Config = {
  apiUrl: getApiUrl()
};
