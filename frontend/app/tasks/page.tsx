"use client";

import Link from "next/link";
import RequireAuth from "../components/RequireAuth";
import { TaskList } from "../components/TaskList";

export default function TasksPage() {
    return (
        <RequireAuth>
            {(session) => (
                <>
                    <div className="tm-page-header d-flex flex-wrap justify-content-between align-items-center gap-3 mb-4">
                        <div>
                            <h1>Tasks</h1>
                            <p className="text-muted-2 mb-0">
                                {session.role === "admin" ? "All tasks in the workspace." : "Tasks assigned to you or created by you."}
                            </p>
                        </div>
                        <Link href="/tasks/new" className="btn btn-primary">
                            + New task
                        </Link>
                    </div>
                    <TaskList session={session} />
                </>
            )}
        </RequireAuth>
    );
}
