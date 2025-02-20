"use client";

import { useEffect, useState } from "react";

const WS_URL = process.env.NEXT_PUBLIC_WS_URL || "ws://localhost:8080/ws";

export const useWebSocket = (url: string) => {
    const [message, setMessage] = useState("");
    const [wsError, setWsError] = useState<string | null>(null);

    useEffect(() => {
        if (typeof window === "undefined") {
            return;
        }

        let ws: WebSocket;
        let reconnectInterval: NodeJS.Timeout;

        const connect = () => {
            try {
                ws = new WebSocket(url);

                ws.onopen = () => {
                    console.log("WebSocket connected");
                    setWsError(null);
                    if (reconnectInterval) clearTimeout(reconnectInterval);
                };

                ws.onmessage = (event) => {
                    setMessage(event.data);
                };

                ws.onerror = (event: Event) => {
                    const error = event as ErrorEvent;
                    const errorMessage = error.message || 'Unknown WebSocket error';
                    setWsError(errorMessage);
                    // console.error("WebSocket error:", errorMessage);
                };

                ws.onclose = (event: CloseEvent) => {
                    const reason = event.reason || 'Connection closed';
                    console.log(`WebSocket disconnected: ${event.code} - ${reason}`);
                    reconnectInterval = setTimeout(connect, 3000);
                };
            } catch (error) {
                const errorMessage = error instanceof Error ? error.message : 'Failed to connect';
                setWsError(errorMessage);
                console.error("WebSocket connection failed:", errorMessage);
                reconnectInterval = setTimeout(connect, 3000);
            }
        };

        connect();

        return () => {
            if (ws) {
                ws.close();
            }
            if (reconnectInterval) {
                clearTimeout(reconnectInterval);
            }
        };
    }, [url]);

    return { message, error: wsError };
};
