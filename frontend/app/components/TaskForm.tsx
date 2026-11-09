"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";
import { api, errorMessage } from "../lib/api";
import { Assignee, Priority, Task, TaskStatus } from "../lib/types";

interface TaskFormProps {
    /** When provided the form edits this task, otherwise it creates a new one. */
    task?: Task;
}

const today = () => new Date().toISOString().split("T")[0];

export const TaskForm = ({ task }: TaskFormProps) => {
    const router = useRouter();
    const editing = Boolean(task);

    const [title, setTitle] = useState(task?.title ?? "");
    const [description, setDescription] = useState(task?.description ?? "");
    const [assignedTo, setAssignedTo] = useState(task ? String(task.assignedTo) : "");
    const [dueDate, setDueDate] = useState(task?.dueDate ? task.dueDate.split("T")[0] : "");
    const [priority, setPriority] = useState<Priority>(task?.priority ?? "MEDIUM");
    const [status, setStatus] = useState<TaskStatus>(task?.status ?? "PENDING");

    const [assignees, setAssignees] = useState<Assignee[]>([]);
    const [assigneesError, setAssigneesError] = useState("");
    const [suggestions, setSuggestions] = useState("");
    const [suggesting, setSuggesting] = useState(false);
    const [saving, setSaving] = useState(false);
    const [error, setError] = useState("");

    useEffect(() => {
        api.get<Assignee[]>("/api/users/assignable")
            .then((res) => setAssignees(res.data))
            .catch((err) => setAssigneesError(errorMessage(err, "Could not load users")));
    }, []);

    const getSuggestions = async () => {
        if (!description.trim()) {
            setError("Write a short description first, then ask for suggestions.");
            return;
        }
        setSuggesting(true);
        setError("");
        try {
            const res = await api.post("/api/tasks/suggestions", { description });
            setSuggestions(res.data.suggestions || "");
            if (!res.data.suggestions) setError("No suggestions received");
        } catch (err) {
            setError(errorMessage(err, "Failed to get suggestions"));
        } finally {
            setSuggesting(false);
        }
    };

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setError("");
        if (!assignedTo) {
            setError("Choose who the task is assigned to.");
            return;
        }
        setSaving(true);
        const payload = {
            title,
            description,
            assignedTo: Number(assignedTo),
            dueDate,
            priority,
            ...(editing ? { status } : {}),
        };
        try {
            if (editing && task) {
                await api.put(`/api/tasks/${task.id}`, payload);
            } else {
                await api.post("/api/tasks", payload);
            }
            router.push("/tasks");
        } catch (err) {
            setError(errorMessage(err, editing ? "Failed to update task" : "Failed to create task"));
            setSaving(false);
        }
    };

    return (
        <form onSubmit={handleSubmit} className="tm-card p-4 p-md-5" noValidate={false}>
            {error && (
                <div className="alert alert-danger" role="alert">
                    {error}
                </div>
            )}

            <div className="mb-4">
                <label htmlFor="title" className="form-label fw-semibold">
                    Title
                </label>
                <input
                    type="text"
                    className="form-control form-control-lg"
                    id="title"
                    value={title}
                    maxLength={200}
                    placeholder="e.g. Prepare the Q3 report"
                    onChange={(e) => setTitle(e.target.value)}
                    required
                />
            </div>

            <div className="mb-4">
                <div className="d-flex justify-content-between align-items-end mb-2 gap-2 flex-wrap">
                    <label htmlFor="description" className="form-label fw-semibold mb-0">
                        Description
                    </label>
                    <button
                        type="button"
                        className="btn btn-sm btn-outline-primary"
                        onClick={getSuggestions}
                        disabled={suggesting}
                    >
                        {suggesting ? (
                            <>
                                <span className="spinner-border spinner-border-sm me-2" aria-hidden="true" />
                                Thinking…
                            </>
                        ) : (
                            "✨ Get AI suggestions"
                        )}
                    </button>
                </div>
                <textarea
                    className="form-control"
                    id="description"
                    rows={5}
                    maxLength={5000}
                    value={description}
                    placeholder="What needs to be done? Add context, links or acceptance criteria."
                    onChange={(e) => setDescription(e.target.value)}
                />
                {suggestions && (
                    <div className="tm-suggestions mt-3" aria-live="polite">
                        <div className="d-flex justify-content-between align-items-center mb-2">
                            <strong>AI suggestions</strong>
                            <button type="button" className="btn-close btn-sm" aria-label="Dismiss suggestions" onClick={() => setSuggestions("")} />
                        </div>
                        {suggestions}
                        <div className="mt-3">
                            <button
                                type="button"
                                className="btn btn-sm btn-primary"
                                onClick={() => setDescription((d) => `${d.trim()}\n\n${suggestions}`.trim())}
                            >
                                Append to description
                            </button>
                        </div>
                    </div>
                )}
            </div>

            <div className="row g-3 mb-4">
                <div className="col-md-6">
                    <label htmlFor="assignedTo" className="form-label fw-semibold">
                        Assign to
                    </label>
                    <select
                        className="form-select"
                        id="assignedTo"
                        value={assignedTo}
                        onChange={(e) => setAssignedTo(e.target.value)}
                        required
                    >
                        <option value="">Select a user…</option>
                        {assignees.map((a) => (
                            <option key={a.id} value={a.id}>
                                {a.username}
                            </option>
                        ))}
                    </select>
                    {assigneesError && <div className="form-text text-danger">{assigneesError}</div>}
                </div>
                <div className="col-md-6">
                    <label htmlFor="dueDate" className="form-label fw-semibold">
                        Due date
                    </label>
                    <input
                        type="date"
                        className="form-control"
                        id="dueDate"
                        value={dueDate}
                        min={editing ? undefined : today()}
                        onChange={(e) => setDueDate(e.target.value)}
                        required
                    />
                </div>
                <div className={editing ? "col-md-6" : "col-12"}>
                    <label htmlFor="priority" className="form-label fw-semibold">
                        Priority
                    </label>
                    <select
                        className="form-select"
                        id="priority"
                        value={priority}
                        onChange={(e) => setPriority(e.target.value as Priority)}
                    >
                        <option value="LOW">Low</option>
                        <option value="MEDIUM">Medium</option>
                        <option value="HIGH">High</option>
                    </select>
                </div>
                {editing && (
                    <div className="col-md-6">
                        <label htmlFor="status" className="form-label fw-semibold">
                            Status
                        </label>
                        <select
                            className="form-select"
                            id="status"
                            value={status}
                            onChange={(e) => setStatus(e.target.value as TaskStatus)}
                        >
                            <option value="PENDING">Pending</option>
                            <option value="IN_PROGRESS">In progress</option>
                            <option value="COMPLETED">Completed</option>
                        </select>
                    </div>
                )}
            </div>

            <div className="d-flex gap-2 justify-content-end flex-wrap">
                <Link href="/tasks" className="btn btn-outline-secondary">
                    Cancel
                </Link>
                <button type="submit" className="btn btn-primary px-4" disabled={saving}>
                    {saving && <span className="spinner-border spinner-border-sm me-2" aria-hidden="true" />}
                    {editing ? "Save changes" : "Create task"}
                </button>
            </div>
        </form>
    );
};
