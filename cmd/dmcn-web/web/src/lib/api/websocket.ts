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

  constructor(token: string) {
    const proto = location.protocol === 'https:' ? 'wss:' : 'ws:';
    this.url = `${proto}//${location.host}/ws?token=${encodeURIComponent(token)}`;
    this.token = token;
  }

  connect() {
    if (this.closed) return;
    this.ws = new WebSocket(this.url);

    this.ws.onopen = () => {
      console.log('WebSocket connected');
    };

    this.ws.onmessage = (event) => {
      try {
        const msg: WSMessage = JSON.parse(event.data);
        const handler = this.handlers.get(msg.type);
        if (handler) handler(msg);
      } catch (e) {
        console.error('WebSocket message parse error:', e);
      }
    };

    this.ws.onclose = () => {
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
    if (this.ws?.readyState === WebSocket.OPEN) {
      this.ws.send(JSON.stringify(msg));
    }
  }

  close() {
    this.closed = true;
    if (this.reconnectTimer) clearTimeout(this.reconnectTimer);
    this.ws?.close();
  }
}
