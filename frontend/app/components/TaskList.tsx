"use client";

import { useCallback, useEffect, useMemo, useState } from "react";
import Link from "next/link";
import { useWebSocket } from "../hooks/useWebSocket";
import { api, errorMessage, Session } from "../lib/api";
import { STATUS_LABELS, Task, TaskStatus } from "../lib/types";
import { PriorityBadge, StatusBadge } from "./StatusBadge";

type Filter = "ALL" | TaskStatus;
const FILTERS: Filter[] = ["ALL", "PENDING", "IN_PROGRESS", "COMPLETED"];

function formatDate(value: string | null) {
    if (!value) return "No due date";
    return new Date(value).toLocaleDateString(undefined, { day: "numeric", month: "short", year: "numeric", timeZone: "UTC" });
}

function isOverdue(task: Task) {
    if (!task.dueDate || task.status === "COMPLETED") return false;
    const due = new Date(task.dueDate);
    due.setUTCHours(23, 59, 59, 999);
    return due.getTime() < Date.now();
}

export const TaskList = ({ session }: { session: Session }) => {
    const [tasks, setTasks] = useState<Task[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState("");
    const [notice, setNotice] = useState("");
    const [filter, setFilter] = useState<Filter>("ALL");
    const [busyId, setBusyId] = useState<number | null>(null);
    const { event, connected } = useWebSocket(session.token);
    const isAdmin = session.role === "admin";

    const fetchTasks = useCallback(async () => {
        try {
            const res = await api.get<Task[]>("/api/tasks");
            setTasks(res.data);
            setError("");
        } catch (err) {
            setError(errorMessage(err, "Failed to load tasks"));
        } finally {
            setLoading(false);
        }
    }, []);

    useEffect(() => {
        fetchTasks();
    }, [fetchTasks]);

    useEffect(() => {
        if (event) fetchTasks();
    }, [event, fetchTasks]);

    const counts = useMemo(() => {
        const c: Record<Filter, number> = { ALL: tasks.length, PENDING: 0, IN_PROGRESS: 0, COMPLETED: 0 };
        tasks.forEach((t) => (c[t.status] += 1));
        return c;
    }, [tasks]);

    const visible = filter === "ALL" ? tasks : tasks.filter((t) => t.status === filter);

    const updateStatus = async (task: Task, status: TaskStatus) => {
        setBusyId(task.id);
        setNotice("");
        try {
            const res = await api.put<Task>(`/api/tasks/${task.id}/status`, { status });
            setTasks((prev) => prev.map((t) => (t.id === task.id ? res.data : t)));
        } catch (err) {
            setNotice(errorMessage(err, "Failed to update task"));
        } finally {
            setBusyId(null);
        }
    };

    const deleteTask = async (task: Task) => {
        if (!window.confirm(`Delete "${task.title}"? This cannot be undone.`)) return;
        setBusyId(task.id);
        setNotice("");
        try {
            await api.delete(`/api/tasks/${task.id}`);
            setTasks((prev) => prev.filter((t) => t.id !== task.id));
        } catch (err) {
            setNotice(errorMessage(err, "Failed to delete task"));
        } finally {
            setBusyId(null);
        }
    };

    if (loading) {
        return (
            <div className="row g-3" aria-busy="true" aria-label="Loading tasks">
                {[0, 1, 2].map((i) => (
                    <div className="col-md-6 col-xl-4" key={i}>
                        <div className="tm-card p-4">
                            <div className="tm-skeleton mb-3" style={{ height: 18, width: "60%" }} />
                            <div className="tm-skeleton mb-2" style={{ height: 12 }} />
                            <div className="tm-skeleton mb-4" style={{ height: 12, width: "80%" }} />
                            <div className="tm-skeleton" style={{ height: 28, width: "40%" }} />
                        </div>
                    </div>
                ))}
            </div>
        );
    }

    if (error) {
        return (
            <div className="tm-card tm-empty" role="alert">
                <div className="tm-empty-icon">!</div>
                <h2 className="h5 text-body">Couldn&apos;t load tasks</h2>
                <p className="mb-3">{error}</p>
                <button className="btn btn-primary" onClick={() => { setLoading(true); fetchTasks(); }}>
                    Try again
                </button>
            </div>
        );
    }

    return (
        <>
            <div className="d-flex flex-wrap justify-content-between align-items-center gap-2 mb-3">
                <div className="tm-filters" role="group" aria-label="Filter by status">
                    {FILTERS.map((f) => (
                        <button
                            key={f}
                            type="button"
                            className={`btn ${filter === f ? "btn-primary" : "btn-outline-secondary"}`}
                            aria-pressed={filter === f}
                            onClick={() => setFilter(f)}
                        >
                            {f === "ALL" ? "All" : STATUS_LABELS[f]}
                            <span className="ms-2 opacity-75">{counts[f]}</span>
                        </button>
                    ))}
                </div>
                <span className={`tm-live ${connected ? "on" : ""}`} title="Real-time updates">
                    {connected ? "Live updates on" : "Live updates offline"}
                </span>
            </div>

            {notice && (
                <div className="alert alert-warning alert-dismissible" role="alert">
                    {notice}
                    <button type="button" className="btn-close" aria-label="Close" onClick={() => setNotice("")} />
                </div>
            )}

            {visible.length === 0 ? (
                <div className="tm-card tm-empty">
                    <div className="tm-empty-icon">✓</div>
                    <h2 className="h5 text-body">{tasks.length === 0 ? "No tasks yet" : "Nothing here"}</h2>
                    <p className="mb-3">
                        {tasks.length === 0
                            ? "Tasks assigned to you or created by you will show up here."
                            : `No ${STATUS_LABELS[filter as TaskStatus].toLowerCase()} tasks.`}
                    </p>
                    {tasks.length === 0 && (
                        <Link href="/tasks/new" className="btn btn-primary">
                            Create your first task
                        </Link>
                    )}
                </div>
            ) : (
                <div className="row g-3">
                    {visible.map((task) => {
                        const canEdit = isAdmin || task.assignedBy === session.userId;
                        const overdue = isOverdue(task);
                        return (
                            <div className="col-md-6 col-xl-4" key={task.id}>
                                <article className="tm-card tm-task p-4">
                                    <div className="d-flex justify-content-between align-items-start gap-2 mb-2">
                                        <h3>{task.title}</h3>
                                        <PriorityBadge priority={task.priority} />
                                    </div>
                                    {task.description && <p className="tm-desc mb-3">{task.description}</p>}

                                    <dl className="tm-meta mb-3 mt-auto">
                                        <div className="d-flex justify-content-between">
                                            <dt className="fw-normal">Assignee</dt>
                                            <dd className="mb-1 text-body">{task.assignedToName || `#${task.assignedTo}`}</dd>
                                        </div>
                                        <div className="d-flex justify-content-between">
                                            <dt className="fw-normal">Created by</dt>
                                            <dd className="mb-1 text-body">{task.assignedByName || `#${task.assignedBy}`}</dd>
                                        </div>
                                        <div className="d-flex justify-content-between">
                                            <dt className="fw-normal">Due</dt>
                                            <dd className={`mb-0 ${overdue ? "tm-overdue" : "text-body"}`}>
                                                {formatDate(task.dueDate)}
                                                {overdue && " · overdue"}
                                            </dd>
                                        </div>
                                    </dl>

                                    <div className="d-flex align-items-center gap-2 flex-wrap border-top pt-3">
                                        <StatusBadge status={task.status} />
                                        <label className="visually-hidden" htmlFor={`status-${task.id}`}>
                                            Status for {task.title}
                                        </label>
                                        <select
                                            id={`status-${task.id}`}
                                            className="form-select form-select-sm w-auto ms-auto"
                                            value={task.status}
                                            disabled={busyId === task.id}
                                            onChange={(e) => updateStatus(task, e.target.value as TaskStatus)}
                                        >
                                            <option value="PENDING">Pending</option>
                                            <option value="IN_PROGRESS">In progress</option>
                                            <option value="COMPLETED">Completed</option>
                                        </select>
                                        {canEdit && (
                                            <Link href={`/tasks/edit/${task.id}`} className="btn btn-sm btn-outline-secondary">
                                                Edit
                                            </Link>
                                        )}
                                        {isAdmin && (
                                            <button
                                                className="btn btn-sm btn-outline-danger"
                                                onClick={() => deleteTask(task)}
                                                disabled={busyId === task.id}
                                            >
                                                Delete
                                            </button>
                                        )}
                                    </div>
                                </article>
                            </div>
                        );
                    })}
                </div>
            )}
        </>
    );
};
