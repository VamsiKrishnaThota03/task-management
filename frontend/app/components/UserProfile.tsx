"use client";

import React, { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import axios from 'axios';

interface UserProfile {
    id: number;
    username: string;
    role: string;
    createdAt: string;
}

export default function UserProfile() {
    const router = useRouter();
    const [profile, setProfile] = useState<UserProfile | null>(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState('');

    useEffect(() => {
        const fetchProfile = async () => {
            try {
                const token = localStorage.getItem('token');
                if (!token) {
                    router.push('/login');
                    return;
                }

                const response = await axios.get('http://localhost:8080/api/profile', {
                    headers: { 
                        Authorization: `Bearer ${token}`,
                        'Content-Type': 'application/json'
                    }
                });

                setProfile(response.data);
            } catch (err: any) {
                setError(err.response?.data?.error || 'Failed to fetch profile');
            } finally {
                setLoading(false);
            }
        };

        fetchProfile();
    }, [router]);

    if (loading) return <div className="text-center p-4">Loading...</div>;
    if (error) return <div className="alert alert-danger m-4">{error}</div>;
    if (!profile) return null;

    return (
        <div className="container mt-4">
            <div className="row justify-content-center">
                <div className="col-md-8">
                    <div className="card">
                        <div className="card-header">
                            <h3 className="mb-0">User Profile</h3>
                        </div>
                        <div className="card-body">
                            <div className="mb-3">
                                <label className="form-label">Username</label>
                                <p className="form-control">{profile.username}</p>
                            </div>
                            <div className="mb-3">
                                <label className="form-label">Role</label>
                                <p className="form-control">{profile.role}</p>
                            </div>
                            <div className="mb-3">
                                <label className="form-label">Member Since</label>
                                <p className="form-control">
                                    {new Date(profile.createdAt).toLocaleDateString()}
                                </p>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    );
} 