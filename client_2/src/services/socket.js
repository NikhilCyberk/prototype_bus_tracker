export class SocketService {
    constructor(url) {
        this.url = url;
        this.socket = null;
        this.listeners = new Map();
    }

    connect() {
        if (this.socket) return;

        this.socket = new WebSocket(this.url);

        this.socket.onopen = () => {
            console.log('WebSocket connected');
            this.notifyListeners('connect', {});
        };

        this.socket.onmessage = (event) => {
            try {
                const message = JSON.parse(event.data);
                this.notifyListeners(message.type, message.content);
            } catch (error) {
                console.error('Error parsing WebSocket message:', error);
            }
        };

        this.socket.onclose = () => {
            console.log('WebSocket disconnected');
            this.notifyListeners('disconnect', {});
            this.socket = null;
            // Attempt to reconnect after 5 seconds
            setTimeout(() => this.connect(), 5000);
        };

        this.socket.onerror = (error) => {
            console.error('WebSocket error:', error);
            this.notifyListeners('error', error);
        };
    }

    disconnect() {
        if (this.socket) {
            this.socket.close();
            this.socket = null;
        }
    }

    addListener(type, callback) {
        if (!this.listeners.has(type)) {
            this.listeners.set(type, new Set());
        }
        this.listeners.get(type).add(callback);
        return () => this.removeListener(type, callback);
    }

    removeListener(type, callback) {
        if (this.listeners.has(type)) {
            this.listeners.get(type).delete(callback);
        }
    }

    notifyListeners(type, data) {
        if (this.listeners.has(type)) {
            this.listeners.get(type).forEach(callback => callback(data));
        }
    }

    send(type, data) {
        if (this.socket && this.socket.readyState === WebSocket.OPEN) {
            this.socket.send(JSON.stringify({ type, data }));
        }
    }
}

// Singleton instance
export const socketService = new SocketService(import.meta.env.VITE_WS_URL || 'ws://localhost:8080/ws');