"use client";

import { useEffect, useState } from "react";
import { api, errorMessage } from "../lib/api";
import { UserProfile as Profile } from "../lib/types";

export default function UserProfile() {
    const [profile, setProfile] = useState<Profile | null>(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState("");

    useEffect(() => {
        api.get<Profile>("/api/profile")
            .then((res) => setProfile(res.data))
            .catch((err) => setError(errorMessage(err, "Failed to load profile")))
            .finally(() => setLoading(false));
    }, []);

    if (loading) {
        return (
            <div className="d-flex justify-content-center py-5" role="status">
                <div className="spinner-border text-primary" aria-hidden="true" />
                <span className="visually-hidden">Loading…</span>
            </div>
        );
    }
    if (error) return <div className="alert alert-danger">{error}</div>;
    if (!profile) return null;

    const initials = profile.username.slice(0, 2).toUpperCase();

    return (
        <div className="tm-card p-4 p-md-5 mx-auto" style={{ maxWidth: 560 }}>
            <div className="d-flex align-items-center gap-3 mb-4">
                <div className="tm-empty-icon m-0 fw-bold" style={{ fontSize: "1.1rem" }} aria-hidden="true">
                    {initials}
                </div>
                <div>
                    <h1 className="h4 fw-bold mb-0">{profile.username}</h1>
                    <span className={`badge ${profile.role === "admin" ? "text-bg-warning" : "text-bg-secondary"}`}>
                        {profile.role}
                    </span>
                </div>
            </div>
            <dl className="row mb-0">
                <dt className="col-sm-4 text-muted-2 fw-normal">User ID</dt>
                <dd className="col-sm-8">{profile.id}</dd>
                <dt className="col-sm-4 text-muted-2 fw-normal">Member since</dt>
                <dd className="col-sm-8 mb-0">
                    {new Date(profile.createdAt).toLocaleDateString(undefined, { day: "numeric", month: "long", year: "numeric" })}
                </dd>
            </dl>
        </div>
    );
}
