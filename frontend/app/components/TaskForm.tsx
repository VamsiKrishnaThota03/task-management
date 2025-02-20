"use client";

import React, { useState } from 'react';
import axios from 'axios';

export const TaskForm = () => {
    const [title, setTitle] = useState('');
    const [description, setDescription] = useState('');
    const [assignedTo, setAssignedTo] = useState('');
    const [dueDate, setDueDate] = useState('');
    const [priority, setPriority] = useState('MEDIUM');
    const [suggestions, setSuggestions] = useState('');
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState('');

    const getSuggestions = async () => {
        try {
            if (!description.trim()) {
                setError('Please enter a task description first');
                return;
            }

            setLoading(true);
            setError('');
            const token = localStorage.getItem('token');
            
            if (!token) {
                setError('Please login first');
                return;
            }

            const response = await axios.post(
                'http://localhost:8080/tasks/suggestions',
                { description },
                { 
                    headers: { 
                        'Authorization': `Bearer ${token}`,
                        'Content-Type': 'application/json'
                    }
                }
            );

            if (response.data.suggestions) {
                setSuggestions(response.data.suggestions);
            } else {
                setError('No suggestions received');
            }
        } catch (err: any) {
            const errorMessage = err.response?.data?.error || 'Failed to get suggestions';
            setError(errorMessage);
        } finally {
            setLoading(false);
        }
    };

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        try {
            setError('');
            const token = localStorage.getItem('token');
            await axios.post(
                'http://localhost:8080/tasks',
                {
                    title,
                    description,
                    assignedTo: parseInt(assignedTo),
                    dueDate: new Date(dueDate).toISOString(),
                    priority
                },
                { 
                    headers: { 
                        'Authorization': `Bearer ${token}`,
                        'Content-Type': 'application/json'
                    }
                }
            );
            
            // Clear form after successful submission
            setTitle('');
            setDescription('');
            setAssignedTo('');
            setDueDate('');
            setPriority('MEDIUM');
            setSuggestions('');
            
        } catch (err: any) {
            setError(err.response?.data?.error || 'Failed to create task');
        }
    };

    return (
        <div className="container py-5">
            <div className="row justify-content-center">
                <div className="col-md-8">
                    <div className="card shadow">
                        <div className="card-header bg-primary text-white">
                            <h3 className="mb-0 text-center">Create New Task</h3>
                        </div>
                        <div className="card-body">
                            {error && (
                                <div className="alert alert-danger" role="alert">
                                    {error}
                                </div>
                            )}

                            <form onSubmit={handleSubmit}>
                                <div className="mb-3">
                                    <label htmlFor="title" className="form-label">Title</label>
                                    <input
                                        type="text"
                                        className="form-control"
                                        id="title"
                                        value={title}
                                        onChange={(e) => setTitle(e.target.value)}
                                        required
                                    />
                                </div>

                                <div className="mb-3">
                                    <label htmlFor="description" className="form-label">Description</label>
                                    <textarea
                                        className="form-control"
                                        id="description"
                                        rows={4}
                                        value={description}
                                        onChange={(e) => setDescription(e.target.value)}
                                        required
                                    />
                                    <div className="mt-2">
                                        <button
                                            type="button"
                                            className="btn btn-outline-primary"
                                            onClick={getSuggestions}
                                            disabled={loading}
                                        >
                                            {loading ? (
                                                <>
                                                    <span className="spinner-border spinner-border-sm me-2" role="status" aria-hidden="true"></span>
                                                    Getting suggestions...
                                                </>
                                            ) : (
                                                'Get AI Suggestions'
                                            )}
                                        </button>
                                    </div>
                                    {suggestions && (
                                        <div className="alert alert-info mt-3">
                                            <h6 className="alert-heading">AI Suggestions:</h6>
                                            <div className="mt-2" style={{ whiteSpace: 'pre-line' }}>
                                                {suggestions}
                                            </div>
                                        </div>
                                    )}
                                </div>

                                <div className="row mb-3">
                                    <div className="col-md-6">
                                        <label htmlFor="priority" className="form-label">Priority</label>
                                        <select
                                            className="form-select"
                                            id="priority"
                                            value={priority}
                                            onChange={(e) => setPriority(e.target.value)}
                                        >
                                            <option value="LOW">Low</option>
                                            <option value="MEDIUM">Medium</option>
                                            <option value="HIGH">High</option>
                                        </select>
                                    </div>
                                    <div className="col-md-6">
                                        <label htmlFor="assignedTo" className="form-label">Assign To (User ID)</label>
                                        <input
                                            type="number"
                                            className="form-control"
                                            id="assignedTo"
                                            value={assignedTo}
                                            onChange={(e) => setAssignedTo(e.target.value)}
                                            required
                                        />
                                    </div>
                                </div>

                                <div className="mb-3">
                                    <label htmlFor="dueDate" className="form-label">Due Date</label>
                                    <input
                                        type="date"
                                        className="form-control"
                                        id="dueDate"
                                        value={dueDate}
                                        onChange={(e) => setDueDate(e.target.value)}
                                        required
                                    />
                                </div>

                                <div className="d-grid gap-2">
                                    <button type="submit" className="btn btn-primary btn-lg">
                                        Create Task
                                    </button>
                                </div>
                            </form>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    );
};
