import React, { createContext, useState, useContext, useEffect, useCallback } from 'react';
import { useAuth } from './AuthContext';

const WebSocketContext = createContext();

export const useWebSocket = () => useContext(WebSocketContext);

export const WebSocketProvider = ({ children }) => {
    const [socket, setSocket] = useState(null);
    const [connected, setConnected] = useState(false);
    const [buses, setBuses] = useState([]);
    const { user } = useAuth();

    const connectWebSocket = useCallback(() => {
        if (user && !socket) {
            const ws = new WebSocket(import.meta.env.VITE_WS_URL || 'ws://localhost:8080/ws');

            ws.onopen = () => {
                console.log('WebSocket connected');
                setConnected(true);
            };

            ws.onmessage = (event) => {
                const message = JSON.parse(event.data);

                switch (message.type) {
                    case 'location_update':
                        // Update buses with new location data
                        setBuses(prevBuses => {
                            const updatedBuses = [...prevBuses];
                            const index = updatedBuses.findIndex(b => b.id === message.content.busId);

                            if (index !== -1) {
                                updatedBuses[index] = {
                                    ...updatedBuses[index],
                                    latitude: message.content.latitude,
                                    longitude: message.content.longitude,
                                    status: message.content.status,
                                    lastUpdated: new Date(message.content.timestamp * 1000)
                                };
                            }

                            return updatedBuses;
                        });
                        break;

                    case 'status_update':
                        // Update bus status
                        setBuses(prevBuses => {
                            const updatedBuses = [...prevBuses];
                            const index = updatedBuses.findIndex(b => b.id === message.content.busId);

                            if (index !== -1) {
                                updatedBuses[index] = {
                                    ...updatedBuses[index],
                                    status: message.content.status,
                                    lastUpdated: new Date(message.content.timestamp * 1000)
                                };
                            }

                            return updatedBuses;
                        });
                        break;

                    default:
                        console.log('Unknown message type:', message.type);
                }
            };

            ws.onclose = () => {
                console.log('WebSocket disconnected');
                setConnected(false);
                setSocket(null);

                // Attempt to reconnect after a delay
                setTimeout(() => {
                    connectWebSocket();
                }, 3000);
            };

            ws.onerror = (error) => {
                console.error('WebSocket error:', error);
                ws.close();
            };

            setSocket(ws);
        }
    }, [user, socket]);

    useEffect(() => {
        if (user) {
            connectWebSocket();
        }

        return () => {
            if (socket) {
                socket.close();
                setSocket(null);
                setConnected(false);
            }
        };
    }, [user, connectWebSocket]);

    const sendMessage = (type, content) => {
        if (socket && connected) {
            socket.send(JSON.stringify({ type, content }));
        }
    };

    const updateBusLocation = (busId, latitude, longitude) => {
        if (user?.role === 'driver') {
            sendMessage('location_update', { busId, latitude, longitude });
        }
    };

    const updateBusStatus = (busId, status) => {
        if (user?.role === 'driver') {
            sendMessage('status_update', { busId, status });
        }
    };

    const value = {
        connected,
        buses,
        setBuses,
        updateBusLocation,
        updateBusStatus
    };

    return <WebSocketContext.Provider value={value}>{children}</WebSocketContext.Provider>;
};