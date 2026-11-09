"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import { api, errorMessage } from "../../../lib/api";
import { Task } from "../../../lib/types";
import { TaskForm } from "../../../components/TaskForm";

export default function EditTaskForm({ taskId }: { taskId: string }) {
    const [task, setTask] = useState<Task | null>(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState("");

    useEffect(() => {
        api.get<Task>(`/api/tasks/${encodeURIComponent(taskId)}`)
            .then((res) => setTask(res.data))
            .catch((err) => setError(errorMessage(err, "Failed to load task")))
            .finally(() => setLoading(false));
    }, [taskId]);

    return (
        <div className="mx-auto" style={{ maxWidth: 760 }}>
            <div className="tm-page-header mb-4">
                <h1>Edit task</h1>
            </div>
            {loading && (
                <div className="d-flex justify-content-center py-5" role="status">
                    <div className="spinner-border text-primary" aria-hidden="true" />
                    <span className="visually-hidden">Loading…</span>
                </div>
            )}
            {error && (
                <div className="alert alert-danger d-flex justify-content-between align-items-center" role="alert">
                    {error}
                    <Link href="/tasks" className="btn btn-sm btn-outline-danger">
                        Back to tasks
                    </Link>
                </div>
            )}
            {task && <TaskForm task={task} />}
        </div>
    );
}
