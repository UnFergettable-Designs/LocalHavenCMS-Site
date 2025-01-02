interface Config {
  apiUrl: string;
}

function getApiUrl(): string {
  const apiBase = import.meta.env.PUBLIC_API_URL || '';
  
  if (apiBase) {
    return apiBase;
  }
  
  if (import.meta.env.DEV) {
    return 'http://localhost:8090';
  }

  return '/api';
}

export const config: Config = {
  apiUrl: getApiUrl()
};
