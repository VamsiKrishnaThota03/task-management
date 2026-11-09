"use client";

import { useEffect } from "react";
import { useRouter } from "next/navigation";
import { Session } from "../lib/api";
import { useSession } from "../hooks/useSession";

/** Renders children only for signed-in users; otherwise redirects to /login. */
export default function RequireAuth({ children }: { children: (session: Session) => React.ReactNode }) {
    const { session, ready } = useSession();
    const router = useRouter();

    useEffect(() => {
        if (ready && !session) router.replace("/login");
    }, [ready, session, router]);

    if (!ready || !session) {
        return (
            <div className="d-flex justify-content-center py-5" role="status" aria-live="polite">
                <div className="spinner-border text-primary" aria-hidden="true" />
                <span className="visually-hidden">Loading…</span>
            </div>
        );
    }
    return <>{children(session)}</>;
}
