export type TaskStatus = "PENDING" | "IN_PROGRESS" | "COMPLETED";
export type Priority = "LOW" | "MEDIUM" | "HIGH";

export interface Task {
    id: number;
    title: string;
    description: string;
    status: TaskStatus;
    priority: Priority;
    dueDate: string | null;
    assignedTo: number;
    assignedToName: string;
    assignedBy: number;
    assignedByName: string;
    createdAt: string;
    updatedAt: string;
}

export interface Assignee {
    id: number;
    username: string;
}

export interface UserProfile {
    id: number;
    username: string;
    role: string;
    createdAt: string;
}

export const STATUS_LABELS: Record<TaskStatus, string> = {
    PENDING: "Pending",
    IN_PROGRESS: "In progress",
    COMPLETED: "Completed",
};

export const PRIORITY_LABELS: Record<Priority, string> = {
    LOW: "Low",
    MEDIUM: "Medium",
    HIGH: "High",
};
