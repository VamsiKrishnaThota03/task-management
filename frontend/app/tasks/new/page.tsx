"use client";

import RequireAuth from "../../components/RequireAuth";
import { TaskForm } from "../../components/TaskForm";

export default function NewTaskPage() {
    return (
        <RequireAuth>
            {() => (
                <div className="mx-auto" style={{ maxWidth: 760 }}>
                    <div className="tm-page-header mb-4">
                        <h1>New task</h1>
                        <p className="text-muted-2 mb-0">Describe the work, pick an assignee and set a due date.</p>
                    </div>
                    <TaskForm />
                </div>
            )}
        </RequireAuth>
    );
}
