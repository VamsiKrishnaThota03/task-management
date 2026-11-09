"use client";

import RequireAuth from "../components/RequireAuth";
import UserProfile from "../components/UserProfile";

export default function ProfilePage() {
    return <RequireAuth>{() => <UserProfile />}</RequireAuth>;
}
