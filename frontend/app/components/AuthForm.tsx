"use client";

import { useState } from "react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { api, errorMessage, saveSession } from "../lib/api";

export default function AuthForm({ mode }: { mode: "login" | "register" }) {
    const router = useRouter();
    const isLogin = mode === "login";
    const [username, setUsername] = useState("");
    const [password, setPassword] = useState("");
    const [confirm, setConfirm] = useState("");
    const [error, setError] = useState("");
    const [submitting, setSubmitting] = useState(false);

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setError("");
        if (!isLogin && password !== confirm) {
            setError("Passwords do not match");
            return;
        }
        setSubmitting(true);
        try {
            if (!isLogin) {
                await api.post("/register", { username, password });
            }
            const res = await api.post("/login", { username, password });
            saveSession({
                token: res.data.token,
                role: res.data.role,
                username: res.data.user?.username ?? username,
                userId: res.data.user?.id ?? 0,
            });
            router.push("/tasks");
        } catch (err) {
            setError(errorMessage(err, isLogin ? "Login failed" : "Registration failed"));
            setSubmitting(false);
        }
    };

    return (
        <div className="tm-auth">
            <div className="tm-card p-4 p-md-5">
                <h1 className="h3 fw-bold mb-1">{isLogin ? "Welcome back" : "Create your account"}</h1>
                <p className="text-muted-2 mb-4">
                    {isLogin ? "Log in to see and update your tasks." : "Sign up to start assigning and tracking tasks."}
                </p>

                {error && (
                    <div className="alert alert-danger py-2" role="alert">
                        {error}
                    </div>
                )}

                <form onSubmit={handleSubmit}>
                    <div className="mb-3">
                        <label htmlFor="username" className="form-label fw-semibold">
                            Username
                        </label>
                        <input
                            id="username"
                            type="text"
                            className="form-control"
                            autoComplete="username"
                            value={username}
                            onChange={(e) => setUsername(e.target.value)}
                            minLength={isLogin ? undefined : 3}
                            maxLength={32}
                            pattern={isLogin ? undefined : "[A-Za-z0-9_.\\-]{3,32}"}
                            title={isLogin ? undefined : "3-32 characters: letters, numbers, '.', '_' or '-'"}
                            required
                            autoFocus
                        />
                    </div>
                    <div className="mb-3">
                        <label htmlFor="password" className="form-label fw-semibold">
                            Password
                        </label>
                        <input
                            id="password"
                            type="password"
                            className="form-control"
                            autoComplete={isLogin ? "current-password" : "new-password"}
                            value={password}
                            onChange={(e) => setPassword(e.target.value)}
                            minLength={isLogin ? undefined : 8}
                            maxLength={72}
                            required
                            aria-describedby={isLogin ? undefined : "password-help"}
                        />
                        {!isLogin && (
                            <div id="password-help" className="form-text">
                                At least 8 characters.
                            </div>
                        )}
                    </div>
                    {!isLogin && (
                        <div className="mb-3">
                            <label htmlFor="confirm" className="form-label fw-semibold">
                                Confirm password
                            </label>
                            <input
                                id="confirm"
                                type="password"
                                className="form-control"
                                autoComplete="new-password"
                                value={confirm}
                                onChange={(e) => setConfirm(e.target.value)}
                                required
                            />
                        </div>
                    )}
                    <button type="submit" className="btn btn-primary w-100 py-2 mt-2" disabled={submitting}>
                        {submitting && <span className="spinner-border spinner-border-sm me-2" aria-hidden="true" />}
                        {isLogin ? "Log in" : "Create account"}
                    </button>
                </form>

                <p className="text-center text-muted-2 mt-4 mb-0">
                    {isLogin ? (
                        <>
                            New here? <Link href="/register">Create an account</Link>
                        </>
                    ) : (
                        <>
                            Already have an account? <Link href="/login">Log in</Link>
                        </>
                    )}
                </p>
            </div>
        </div>
    );
}
