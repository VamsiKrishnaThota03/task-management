"use client";

import React, { useEffect, useState } from 'react';
import { useWebSocket } from '../hooks/useWebSocket';
import axios from 'axios';

interface Task {
    ID: number;
    Title: string;
    Description: string;
    Status: string;
    AssignedTo: number;
    DueDate: string;
    Priority: string;
}

export const TaskList = () => {
    const [tasks, setTasks] = useState<Task[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState('');
    const { message, error: wsError } = useWebSocket("ws://localhost:8080/ws");
    const role = typeof window !== 'undefined' ? localStorage.getItem('role') : null;

    const fetchTasks = async () => {
        try {
            const token = localStorage.getItem('token');
            const response = await axios.get('http://localhost:8080/tasks', {
                headers: { Authorization: `Bearer ${token}` }
            });
            setTasks(response.data);
            setLoading(false);
        } catch (err) {
            setError('Failed to fetch tasks');
            setLoading(false);
        }
    };

    const updateTaskStatus = async (taskId: number, status: string) => {
        try {
            const token = localStorage.getItem('token');
            await axios.put(`http://localhost:8080/tasks/${taskId}/status`, 
                { status },
                { headers: { Authorization: `Bearer ${token}` }}
            );
            fetchTasks(); // Refresh task list
        } catch (err) {
            setError('Failed to update task');
        }
    };

    useEffect(() => {
        fetchTasks();
    }, []);

    useEffect(() => {
        if (message) {
            fetchTasks(); // Refresh when receiving WebSocket updates
        }
    }, [message]);

    if (wsError) {
        console.warn('WebSocket connection issue:', wsError);
        // Optionally show a warning to the user
    }

    if (loading) return <div>Loading tasks...</div>;
    if (error) return <div className="text-red-500">{error}</div>;

    return (
        <div className="space-y-4">
            {tasks.map((task) => (
                <div key={task.ID} className="border p-4 rounded-lg shadow">
                    <h3 className="font-bold">{task.Title}</h3>
                    <p className="text-gray-600">{task.Description}</p>
                    <div className="mt-2 flex justify-between items-center">
                        <span className={`px-2 py-1 rounded ${
                            task.Status === 'COMPLETED' ? 'bg-green-200' : 
                            task.Status === 'IN_PROGRESS' ? 'bg-yellow-200' : 'bg-gray-200'
                        }`}>
                            {task.Status}
                        </span>
                        <select 
                            value={task.Status}
                            onChange={(e) => updateTaskStatus(task.ID, e.target.value)}
                            className="border rounded p-1"
                        >
                            <option value="PENDING">Pending</option>
                            <option value="IN_PROGRESS">In Progress</option>
                            <option value="COMPLETED">Completed</option>
                        </select>
                    </div>
                </div>
            ))}
        </div>
    );
};
