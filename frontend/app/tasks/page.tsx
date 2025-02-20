"use client";

import React from 'react';
import { useRouter } from 'next/navigation';
import { TaskList } from '../components/TaskList';

export default function TasksPage() {
    const router = useRouter();

    return (
        <div className="container">
            <div className="d-flex justify-content-between align-items-center mb-4">
                <h1>Tasks</h1>
                <button 
                    className="btn btn-primary"
                    onClick={() => router.push('/tasks/new')}
                >
                    Create New Task
                </button>
            </div>
            <TaskList />
        </div>
    );
}