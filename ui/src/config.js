export const BASE_URL = process.env.REACT_APP_BASE || 'http://localhost:8080';

const wsProtocol = window.location.protocol === 'https:' ? 'wss' : 'ws';

export const BASE_URL_WS = process.env.REACT_APP_BASE ?
  `${wsProtocol}://${window.location.host}/api` : 'ws://localhost:8080';