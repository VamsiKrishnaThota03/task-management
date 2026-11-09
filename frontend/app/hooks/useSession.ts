"use client";

import { useEffect, useState } from "react";
import { AUTH_EVENT, getSession, Session } from "../lib/api";

/**
 * Current session from localStorage. `ready` is false during the first render
 * (server render / hydration) so pages can avoid flashing the wrong state.
 */
export function useSession() {
    const [session, setSession] = useState<Session | null>(null);
    const [ready, setReady] = useState(false);

    useEffect(() => {
        const update = () => setSession(getSession());
        update();
        setReady(true);
        window.addEventListener("storage", update);
        window.addEventListener(AUTH_EVENT, update);
        return () => {
            window.removeEventListener("storage", update);
            window.removeEventListener(AUTH_EVENT, update);
        };
    }, []);

    return { session, ready };
}
