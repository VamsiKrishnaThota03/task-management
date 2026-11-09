"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import { usePathname, useRouter } from "next/navigation";
import { clearSession } from "../lib/api";
import { useSession } from "../hooks/useSession";

export default function Navbar() {
    const router = useRouter();
    const pathname = usePathname();
    const { session, ready } = useSession();
    const [open, setOpen] = useState(false);
    const [theme, setTheme] = useState<"light" | "dark">("light");

    useEffect(() => {
        const current = document.documentElement.getAttribute("data-bs-theme");
        setTheme(current === "dark" ? "dark" : "light");
    }, []);

    // Close the mobile menu after navigating.
    useEffect(() => setOpen(false), [pathname]);

    const toggleTheme = () => {
        const next = theme === "dark" ? "light" : "dark";
        setTheme(next);
        document.documentElement.setAttribute("data-bs-theme", next);
        localStorage.setItem("theme", next);
    };

    const logout = () => {
        clearSession();
        router.push("/login");
    };

    const navLink = (href: string, label: string) => (
        <li className="nav-item">
            <Link
                href={href}
                className={`nav-link ${pathname === href ? "active" : ""}`}
                aria-current={pathname === href ? "page" : undefined}
            >
                {label}
            </Link>
        </li>
    );

    return (
        <nav className="navbar navbar-expand-md tm-navbar sticky-top" aria-label="Main navigation">
            <div className="container">
                <Link className="tm-brand me-4" href="/">
                    {/* eslint-disable-next-line @next/next/no-img-element */}
                    <img src="/icon.svg" alt="" />
                    Task Manager
                </Link>
                <button
                    className="navbar-toggler"
                    type="button"
                    aria-controls="main-nav"
                    aria-expanded={open}
                    aria-label="Toggle navigation"
                    onClick={() => setOpen((o) => !o)}
                >
                    <span className="navbar-toggler-icon" />
                </button>

                <div className={`collapse navbar-collapse ${open ? "show" : ""}`} id="main-nav">
                    <ul className="navbar-nav me-auto gap-1 mt-2 mt-md-0">
                        {ready && session && (
                            <>
                                {navLink("/tasks", "Tasks")}
                                {navLink("/tasks/new", "New task")}
                                {navLink("/profile", "Profile")}
                            </>
                        )}
                    </ul>

                    <div className="d-flex align-items-center gap-2 py-2 py-md-0">
                        <button
                            type="button"
                            className="btn btn-sm btn-outline-secondary"
                            onClick={toggleTheme}
                            aria-label={`Switch to ${theme === "dark" ? "light" : "dark"} mode`}
                            title={`Switch to ${theme === "dark" ? "light" : "dark"} mode`}
                        >
                            {theme === "dark" ? "☀ Light" : "☾ Dark"}
                        </button>

                        {ready &&
                            (session ? (
                                <>
                                    <span className="small text-muted-2 d-none d-lg-inline">
                                        {session.username}
                                        {session.role === "admin" && (
                                            <span className="badge text-bg-warning ms-2">admin</span>
                                        )}
                                    </span>
                                    <button onClick={logout} className="btn btn-sm btn-outline-primary">
                                        Log out
                                    </button>
                                </>
                            ) : (
                                <>
                                    <Link href="/login" className="btn btn-sm btn-outline-primary">
                                        Log in
                                    </Link>
                                    <Link href="/register" className="btn btn-sm btn-primary">
                                        Sign up
                                    </Link>
                                </>
                            ))}
                    </div>
                </div>
            </div>
        </nav>
    );
}
