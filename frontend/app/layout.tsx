import type { Metadata, Viewport } from "next";
import "bootstrap/dist/css/bootstrap.min.css";
import "./styles/globals.css";
import Navbar from "./components/Navbar";

export const metadata: Metadata = {
    title: {
        default: "Task Manager",
        template: "%s · Task Manager",
    },
    description: "Assign, track and complete tasks with real-time updates and AI-assisted task descriptions.",
};

export const viewport: Viewport = {
    width: "device-width",
    initialScale: 1,
};

// Applies the saved theme before first paint to avoid a light/dark flash.
const themeScript = `try{var t=localStorage.getItem('theme');if(!t){t=window.matchMedia('(prefers-color-scheme: dark)').matches?'dark':'light'}document.documentElement.setAttribute('data-bs-theme',t)}catch(e){}`;

export default function RootLayout({ children }: { children: React.ReactNode }) {
    return (
        <html lang="en" data-bs-theme="light" suppressHydrationWarning>
            <head>
                <script dangerouslySetInnerHTML={{ __html: themeScript }} />
            </head>
            <body>
                <Navbar />
                <main className="container py-4 py-md-5">{children}</main>
                <footer className="tm-footer py-3">
                    <div className="container d-flex flex-wrap justify-content-between gap-2">
                        <span>Task Manager · Go + Next.js</span>
                        <a href="https://github.com/VamsiKrishnaThota03/taskflow" target="_blank" rel="noreferrer">
                            Source on GitHub
                        </a>
                    </div>
                </footer>
            </body>
        </html>
    );
}
