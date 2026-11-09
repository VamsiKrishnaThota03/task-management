"use client";

import { useEffect, useState } from "react";
import { WS_URL } from "../lib/api";

export interface TaskEvent {
    type: string;
    taskId: number;
    receivedAt: number;
}

/**
 * Subscribes to task change events. Reconnects with exponential backoff.
 * The JWT is passed as a query parameter because browsers cannot set headers
 * on WebSocket requests.
 */
export const useWebSocket = (token: string | null | undefined) => {
    const [event, setEvent] = useState<TaskEvent | null>(null);
    const [connected, setConnected] = useState(false);

    useEffect(() => {
        if (!token || typeof window === "undefined") return;

        let ws: WebSocket | null = null;
        let retry: ReturnType<typeof setTimeout> | undefined;
        let attempts = 0;
        let closedByUs = false;

        const connect = () => {
            ws = new WebSocket(`${WS_URL}?token=${encodeURIComponent(token)}`);
            ws.onopen = () => {
                attempts = 0;
                setConnected(true);
            };
            ws.onmessage = (msg) => {
                try {
                    const data = JSON.parse(msg.data);
                    setEvent({ type: data.type, taskId: data.taskId, receivedAt: Date.now() });
                } catch {
                    /* ignore non-JSON messages */
                }
            };
            ws.onclose = () => {
                setConnected(false);
                if (closedByUs) return;
                const delay = Math.min(30000, 1000 * 2 ** attempts);
                attempts += 1;
                retry = setTimeout(connect, delay);
            };
        };

        connect();

        return () => {
            closedByUs = true;
            if (retry) clearTimeout(retry);
            ws?.close();
        };
    }, [token]);

    return { event, connected };
};
