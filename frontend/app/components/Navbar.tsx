"use client";

import React, { useEffect, useState } from 'react';
import Link from 'next/link';
import { useRouter } from 'next/navigation';
import { logoutUser } from '../api/auth';

export default function Navbar() {
    const router = useRouter();
    const [isLoggedIn, setIsLoggedIn] = useState(false);
    const [isDarkMode, setIsDarkMode] = useState(false);
    const [role, setRole] = useState<string | null>(null);

    useEffect(() => {
        // Check auth status on mount and when localStorage changes
        const checkAuth = () => {
            const token = localStorage.getItem('token');
            setIsLoggedIn(!!token);
        };

        // Initial check
        checkAuth();

        // Listen for storage events (when token is added/removed)
        window.addEventListener('storage', checkAuth);

        // Custom event for auth state changes
        window.addEventListener('authStateChange', checkAuth);

        // Check theme preference from localStorage
        const savedTheme = localStorage.getItem('theme');
        if (savedTheme === 'dark') {
            setIsDarkMode(true);
            document.documentElement.setAttribute('data-bs-theme', 'dark');
        }

        setRole(localStorage.getItem('role'));

        return () => {
            window.removeEventListener('storage', checkAuth);
            window.removeEventListener('authStateChange', checkAuth);
        };
    }, []);

    const handleLogout = () => {
        logoutUser();
        router.push('/login');
    };

    const handleCreateTask = () => {
        router.push('/tasks/new');  // Updated path to match the new task page
    };

    const handleViewTasks = () => {
        router.push('/tasks');
    };

    const toggleTheme = () => {
        const newTheme = !isDarkMode;
        setIsDarkMode(newTheme);
        document.documentElement.setAttribute('data-bs-theme', newTheme ? 'dark' : 'light');
        localStorage.setItem('theme', newTheme ? 'dark' : 'light');
    };

    return (
        <nav className={`navbar navbar-expand-lg ${isDarkMode ? 'navbar-dark bg-dark' : 'navbar-light bg-light'}`}>
            <div className="container-fluid">
                <Link className="navbar-brand" href="/">
                    Task Management System
                </Link>
                <button 
                    className="navbar-toggler" 
                    type="button" 
                    data-bs-toggle="collapse" 
                    data-bs-target="#navbarSupportedContent"
                >
                    <span className="navbar-toggler-icon"></span>
                </button>

                <div className="collapse navbar-collapse" id="navbarSupportedContent">
                    <ul className="navbar-nav me-auto mb-2 mb-lg-0">
                        <li className="nav-item">
                            <Link className="nav-link active" href="/">
                                Home
                            </Link>
                        </li>
                        {role && (
                            <>
                                <li className="nav-item dropdown">
                                    <a 
                                        className="nav-link dropdown-toggle" 
                                        href="#" 
                                        role="button" 
                                        data-bs-toggle="dropdown"
                                    >
                                        Manage
                                    </a>
                                    <ul className="dropdown-menu">
                                        <li>
                                            <button className="dropdown-item" onClick={handleCreateTask}>
                                                Create Task
                                            </button>
                                        </li>
                                        <li>
                                            <button className="dropdown-item" onClick={handleViewTasks}>
                                                View Tasks
                                            </button>
                                        </li>
                                    </ul>
                                </li>
                            </>
                        )}
                    </ul>

                    <div className="d-flex align-items-center">
                        <div className="form-check form-switch me-3">
                            <input
                                className="form-check-input"
                                type="checkbox"
                                role="switch"
                                id="themeSwitch"
                                checked={isDarkMode}
                                onChange={toggleTheme}
                            />
                            <label className={`form-check-label ${isDarkMode ? 'text-light' : 'text-dark'}`} htmlFor="themeSwitch">
                                {isDarkMode ? 'Dark' : 'Light'} Mode
                            </label>
                        </div>

                        {!isLoggedIn ? (
                            <>
                                <Link href="/login" className={`btn ${isDarkMode ? 'btn-outline-light' : 'btn-outline-dark'} me-2`}>
                                    Login
                                </Link>
                                <Link href="/register" className={`btn ${isDarkMode ? 'btn-light' : 'btn-dark'}`}>
                                    Register
                                </Link>
                            </>
                        ) : (
                            <button 
                                onClick={handleLogout} 
                                className={`btn ${isDarkMode ? 'btn-outline-light' : 'btn-outline-dark'}`}
                            >
                                Logout
                            </button>
                        )}
                    </div>
                </div>
            </div>
        </nav>
    );
}
