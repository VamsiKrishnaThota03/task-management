import React from 'react';
import { TaskForm } from './components/TaskForm';
import { TaskList } from './components/TaskList';


export default function Home() {
    return (
        <div>
            <TaskForm />
            <TaskList />
        </div>
    );
}