import { Priority, PRIORITY_LABELS, STATUS_LABELS, TaskStatus } from "../lib/types";

export function StatusBadge({ status }: { status: TaskStatus }) {
    return <span className={`tm-badge tm-status-${status}`}>{STATUS_LABELS[status] ?? status}</span>;
}

export function PriorityBadge({ priority }: { priority: Priority }) {
    return (
        <span className={`tm-badge tm-priority-${priority}`} title="Priority">
            {PRIORITY_LABELS[priority] ?? priority}
        </span>
    );
}
