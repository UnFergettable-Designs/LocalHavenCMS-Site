interface Config {
  apiUrl: string;
}

function getApiUrl() {
  const url = import.meta.env.VITE_API_URL;
  if (!url) {
    console.warn('VITE_API_URL not set, using default');
    return 'http://localhost:8090';
  }
  return url;
}

export const config: Config = {
  apiUrl: getApiUrl()
};
