import axios, { AxiosError } from "axios";

export const API_URL = (process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080").replace(/\/$/, "");

/** WebSocket URL derived from NEXT_PUBLIC_WS_URL or the API URL. */
export const WS_URL =
    process.env.NEXT_PUBLIC_WS_URL || `${API_URL.replace(/^http/, "ws")}/ws`;

const TOKEN_KEY = "token";
const ROLE_KEY = "role";
const USER_KEY = "username";
const USER_ID_KEY = "userId";
export const AUTH_EVENT = "authStateChange";

export interface Session {
    token: string;
    role: string;
    username: string;
    userId: number;
}

export function getSession(): Session | null {
    if (typeof window === "undefined") return null;
    const token = localStorage.getItem(TOKEN_KEY);
    if (!token) return null;
    return {
        token,
        role: localStorage.getItem(ROLE_KEY) || "user",
        username: localStorage.getItem(USER_KEY) || "",
        userId: Number(localStorage.getItem(USER_ID_KEY) || 0),
    };
}

export function saveSession(session: Session) {
    localStorage.setItem(TOKEN_KEY, session.token);
    localStorage.setItem(ROLE_KEY, session.role);
    localStorage.setItem(USER_KEY, session.username);
    localStorage.setItem(USER_ID_KEY, String(session.userId));
    window.dispatchEvent(new Event(AUTH_EVENT));
}

export function clearSession() {
    localStorage.removeItem(TOKEN_KEY);
    localStorage.removeItem(ROLE_KEY);
    localStorage.removeItem(USER_KEY);
    localStorage.removeItem(USER_ID_KEY);
    window.dispatchEvent(new Event(AUTH_EVENT));
}

export const api = axios.create({
    baseURL: API_URL,
    headers: { "Content-Type": "application/json" },
    timeout: 30000,
});

api.interceptors.request.use((config) => {
    const session = getSession();
    if (session) {
        config.headers.Authorization = `Bearer ${session.token}`;
    }
    return config;
});

api.interceptors.response.use(
    (response) => response,
    (error: AxiosError) => {
        // Expired or invalid token: drop the session so the UI asks for a new login.
        if (error.response?.status === 401 && getSession()) {
            clearSession();
        }
        return Promise.reject(error);
    }
);

/** Human-readable message for an API error. */
export function errorMessage(err: unknown, fallback: string): string {
    if (axios.isAxiosError(err)) {
        const data = err.response?.data as { error?: string } | undefined;
        if (data?.error) return data.error;
        if (!err.response) return `Cannot reach the server at ${API_URL}. Is the backend running?`;
    }
    return fallback;
}
