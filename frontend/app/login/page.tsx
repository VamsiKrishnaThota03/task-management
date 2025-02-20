"use client";

import React, { useState } from "react";
import { useRouter } from "next/navigation";
import { loginUser } from "../api/auth";
import axios from "axios";

export default function LoginPage() {
    const [username, setUsername] = useState("");
    const [password, setPassword] = useState("");
    const [error, setError] = useState("");
    const router = useRouter();

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        try {
            const response = await axios.post('http://localhost:8080/login', { username, password });
            localStorage.setItem('token', response.data.token);
            // Dispatch custom event to notify auth state change
            window.dispatchEvent(new Event('authStateChange'));
            router.push('/tasks');
        } catch (err: any) {
            setError(err.response?.data?.error || 'Login failed');
        }
    };

    return (
        <div className="flex flex-col items-center mt-10">
            <h2 className="text-2xl font-bold">Login</h2>
            {error && <p className="text-red-500">{error}</p>}
            <form onSubmit={handleSubmit} className="flex flex-col gap-2">
                <input type="text" placeholder="Username" value={username} onChange={(e) => setUsername(e.target.value)} className="border p-2" />
                <input type="password" placeholder="Password" value={password} onChange={(e) => setPassword(e.target.value)} className="border p-2" />
                <button type="submit" className="bg-blue-500 text-white px-4 py-2">Login</button>
            </form>
        </div>
    );
}
