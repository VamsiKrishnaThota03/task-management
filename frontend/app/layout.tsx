"use client";

import React, { useEffect } from 'react';
import Navbar from "./components/Navbar";
import Script from 'next/script';

export default function Layout({ children }: { children: React.ReactNode }) {
    useEffect(() => {
        // Check saved theme preference
        const savedTheme = localStorage.getItem('theme');
        if (savedTheme) {
            document.documentElement.setAttribute('data-bs-theme', savedTheme);
        }
    }, []);

    return (
        <html lang="en">
            <head>
                <link 
                    href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.0/dist/css/bootstrap.min.css" 
                    rel="stylesheet"
                />
            </head>
            <body>
                <Navbar />
                <main className="container py-4">
                    {children}
                </main>
                <Script 
                    src="https://cdn.jsdelivr.net/npm/bootstrap@5.3.0/dist/js/bootstrap.bundle.min.js"
                    strategy="afterInteractive"
                />
            </body>
        </html>
    );
}
