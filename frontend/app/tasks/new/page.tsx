"use client";

import React from 'react';
import { TaskForm } from '../../components/TaskForm';

export default function NewTaskPage() {
    return (
        <div className="container">
            <h1 className="mb-4">Create New Task</h1>
            <TaskForm />
        </div>
    );
} 