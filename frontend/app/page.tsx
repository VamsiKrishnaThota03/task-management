"use client";

import Link from "next/link";
import { useSession } from "./hooks/useSession";

const FEATURES = [
    { icon: "⇄", title: "Assign and track", text: "Create tasks with due dates and priorities, assign them to teammates and move them from pending to done." },
    { icon: "✨", title: "AI suggestions", text: "Ask Google Gemini to tighten a task description and break it into concrete subtasks." },
    { icon: "⚡", title: "Real-time updates", text: "Task lists refresh over a WebSocket whenever someone creates, edits or completes a task." },
];

export default function Home() {
    const { session, ready } = useSession();

    return (
        <>
            <section className="tm-hero">
                <span className="badge rounded-pill text-bg-light border mb-3 px-3 py-2 fw-medium">Go · Gin · Next.js · Gemini</span>
                <h1>
                    Keep your team&apos;s work
                    <br className="d-none d-md-block" /> moving forward.
                </h1>
                <p className="lead">
                    A simple task manager with roles, live updates and AI help for writing better tasks.
                </p>
                <div className="d-flex gap-2 justify-content-center flex-wrap" style={{ minHeight: 48 }}>
                    {ready &&
                        (session ? (
                            <>
                                <Link href="/tasks" className="btn btn-primary btn-lg px-4">
                                    Go to my tasks
                                </Link>
                                <Link href="/tasks/new" className="btn btn-outline-primary btn-lg px-4">
                                    New task
                                </Link>
                            </>
                        ) : (
                            <>
                                <Link href="/register" className="btn btn-primary btn-lg px-4">
                                    Get started
                                </Link>
                                <Link href="/login" className="btn btn-outline-primary btn-lg px-4">
                                    Log in
                                </Link>
                            </>
                        ))}
                </div>
            </section>

            <section className="row g-3 pb-4">
                {FEATURES.map((f) => (
                    <div className="col-md-4" key={f.title}>
                        <div className="tm-card tm-feature">
                            <div className="tm-empty-icon" aria-hidden="true">{f.icon}</div>
                            <h2>{f.title}</h2>
                            <p>{f.text}</p>
                        </div>
                    </div>
                ))}
            </section>
        </>
    );
}
