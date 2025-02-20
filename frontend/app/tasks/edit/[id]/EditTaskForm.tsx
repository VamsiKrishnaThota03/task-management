"use client";

import React, { useEffect, useState } from 'react';
import axios from 'axios';
import { useRouter } from 'next/navigation';

interface EditTaskFormProps {
    taskId: string;
}

export default function EditTaskForm({ taskId }: EditTaskFormProps) {
    const router = useRouter();
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState('');
    const [formData, setFormData] = useState({
        title: '',
        description: '',
        status: '',
        priority: '',
        assignedTo: '',
        dueDate: ''
    });

    useEffect(() => {
        const fetchTask = async () => {
            try {
                const token = localStorage.getItem('token');
                if (!token) {
                    router.push('/login');
                    return;
                }

                const response = await axios.get(`http://localhost:8080/tasks/${taskId}`, {
                    headers: { 'Authorization': `Bearer ${token}` }
                });

                const taskData = response.data;
                setFormData({
                    title: taskData.Title || '',
                    description: taskData.Description || '',
                    status: taskData.Status || '',
                    priority: taskData.Priority || '',
                    assignedTo: taskData.AssignedTo?.toString() || '',
                    dueDate: taskData.DueDate ? new Date(taskData.DueDate).toISOString().split('T')[0] : ''
                });
            } catch (err: any) {
                setError(err.response?.data?.error || 'Failed to fetch task');
            } finally {
                setLoading(false);
            }
        };

        fetchTask();
    }, [taskId, router]);

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        try {
            const token = localStorage.getItem('token');
            await axios.put(`http://localhost:8080/tasks/${taskId}`, 
                {
                    Title: formData.title,
                    Description: formData.description,
                    Status: formData.status,
                    Priority: formData.priority,
                    AssignedTo: parseInt(formData.assignedTo),
                    DueDate: new Date(formData.dueDate).toISOString()
                },
                {
                    headers: { 'Authorization': `Bearer ${token}` }
                }
            );
            router.push('/tasks');
        } catch (err: any) {
            setError(err.response?.data?.error || 'Failed to update task');
        }
    };

    const handleInputChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement>) => {
        const { name, value } = e.target;
        setFormData(prev => ({
            ...prev,
            [name]: value
        }));
    };

    if (loading) return <div className="text-center p-4">Loading...</div>;
    if (error) return <div className="alert alert-danger m-4">{error}</div>;

    return (
        <div className="container mt-4">
            <div className="row justify-content-center">
                <div className="col-md-8">
                    <div className="card">
                        <div className="card-header">
                            <h3 className="mb-0">Edit Task</h3>
                        </div>
                        <div className="card-body">
                            <form onSubmit={handleSubmit}>
                                <div className="mb-3">
                                    <label htmlFor="title" className="form-label">Title</label>
                                    <input
                                        type="text"
                                        className="form-control"
                                        id="title"
                                        name="title"
                                        value={formData.title}
                                        onChange={handleInputChange}
                                        required
                                    />
                                </div>

                                <div className="mb-3">
                                    <label htmlFor="description" className="form-label">Description</label>
                                    <textarea
                                        className="form-control"
                                        id="description"
                                        name="description"
                                        value={formData.description}
                                        onChange={handleInputChange}
                                        rows={4}
                                        required
                                    />
                                </div>

                                <div className="mb-3">
                                    <label htmlFor="status" className="form-label">Status</label>
                                    <select
                                        className="form-select"
                                        id="status"
                                        name="status"
                                        value={formData.status}
                                        onChange={handleInputChange}
                                        required
                                    >
                                        <option value="">Select Status</option>
                                        <option value="pending">Pending</option>
                                        <option value="in_progress">In Progress</option>
                                        <option value="completed">Completed</option>
                                    </select>
                                </div>

                                <div className="mb-3">
                                    <label htmlFor="priority" className="form-label">Priority</label>
                                    <select
                                        className="form-select"
                                        id="priority"
                                        name="priority"
                                        value={formData.priority}
                                        onChange={handleInputChange}
                                        required
                                    >
                                        <option value="">Select Priority</option>
                                        <option value="low">Low</option>
                                        <option value="medium">Medium</option>
                                        <option value="high">High</option>
                                    </select>
                                </div>

                                <div className="mb-3">
                                    <label htmlFor="assignedTo" className="form-label">Assigned To (User ID)</label>
                                    <input
                                        type="number"
                                        className="form-control"
                                        id="assignedTo"
                                        name="assignedTo"
                                        value={formData.assignedTo}
                                        onChange={handleInputChange}
                                    />
                                </div>

                                <div className="mb-3">
                                    <label htmlFor="dueDate" className="form-label">Due Date</label>
                                    <input
                                        type="date"
                                        className="form-control"
                                        id="dueDate"
                                        name="dueDate"
                                        value={formData.dueDate}
                                        onChange={handleInputChange}
                                        required
                                    />
                                </div>

                                <div className="d-flex gap-2">
                                    <button type="submit" className="btn btn-primary">
                                        Update Task
                                    </button>
                                    <button 
                                        type="button" 
                                        className="btn btn-secondary"
                                        onClick={() => router.push('/tasks')}
                                    >
                                        Cancel
                                    </button>
                                </div>
                            </form>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    );
} 