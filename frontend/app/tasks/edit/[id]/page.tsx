"use client";

import { useParams } from "next/navigation";
import RequireAuth from "../../../components/RequireAuth";
import EditTaskForm from "./EditTaskForm";

export default function EditTaskPage() {
    const params = useParams<{ id: string }>();
    return <RequireAuth>{() => <EditTaskForm taskId={params.id} />}</RequireAuth>;
}
