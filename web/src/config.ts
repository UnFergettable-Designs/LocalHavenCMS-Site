interface Config {
  apiUrl: string;
}

function getApiUrl(): string {
  if (import.meta.env.DEV) {
    return import.meta.env.PUBLIC_API_URL || 'http://localhost:8090';
  }
  
  // In production, use relative URL to respect the current domain
  return '/api';
}

export const config: Config = {
  apiUrl: getApiUrl()
};
