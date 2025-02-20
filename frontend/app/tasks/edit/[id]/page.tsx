import EditTaskForm from './EditTaskForm';

export default function EditTaskPage({ params }: { params: { id: string } }) {
    return <EditTaskForm taskId={params.id} />;
} 