export type MessageHandler = (msg: WSMessage) => void;

export interface WSMessage {
  id: string;
  type: string;
  data?: unknown;
}

export class DMCNWebSocket {
  private ws: WebSocket | null = null;
  private url: string;
  private token: string;
  private handlers = new Map<string, MessageHandler>();
  private reconnectTimer: ReturnType<typeof setTimeout> | null = null;
  private closed = false;
  private authenticated = false;

  constructor(token: string) {
    const proto = location.protocol === 'https:' ? 'wss:' : 'ws:';
    this.url = `${proto}//${location.host}/ws`;
    this.token = token;
  }

  connect() {
    if (this.closed) return;
    this.authenticated = false;
    this.ws = new WebSocket(this.url);

    this.ws.onopen = () => {
      // Send session token as the first message rather than in the URL
      // query string, so it doesn't leak into logs or browser history.
      this.ws!.send(JSON.stringify({
        id: 'auth',
        type: 'authenticate',
        data: { token: this.token },
      }));
    };

    this.ws.onmessage = (event) => {
      try {
        const msg: WSMessage = JSON.parse(event.data);
        if (msg.type === 'authenticated') {
          this.authenticated = true;
          return;
        }
        const handler = this.handlers.get(msg.type);
        if (handler) handler(msg);
      } catch (e) {
        console.error('WebSocket message parse error:', e);
      }
    };

    this.ws.onclose = () => {
      this.authenticated = false;
      if (!this.closed) {
        this.reconnectTimer = setTimeout(() => this.connect(), 3000);
      }
    };

    this.ws.onerror = () => {
      this.ws?.close();
    };
  }

  on(type: string, handler: MessageHandler) {
    this.handlers.set(type, handler);
  }

  send(msg: WSMessage) {
    if (this.ws?.readyState === WebSocket.OPEN && this.authenticated) {
      this.ws.send(JSON.stringify(msg));
    }
  }

  close() {
    this.closed = true;
    if (this.reconnectTimer) clearTimeout(this.reconnectTimer);
    this.ws?.close();
  }
}
