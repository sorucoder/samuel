import React from 'react';
import ReactDOM from 'react-dom/client';
import { createBrowserRouter, RouterProvider } from 'react-router-dom';

import 'bootstrap/dist/css/bootstrap.min.css';
import './styles/theme.css';

import LoginPage, { loginLoader, loginAction } from './pages/LoginPage.js';
import LogoutPage, { logoutLoader } from './pages/LogoutPage.js';
import PasswordChangePage, { passwordChangeAction } from './pages/PasswordChangePage.js';
import DashboardPage, { dashboardLoader } from './pages/DashboardPage.js';
import AuditViewPage, { auditViewLoader } from './pages/audit/AuditViewPage.js';
import ErrorPage from './pages/ErrorPage.js';

const router = createBrowserRouter([
    {
        path: '/',
        loader: () => {
            throw {
                status: 404,
                message: 'Not Implemented Yet',
                redirect: {
                    path: '/login',
                    name: 'Login Page'
                }
            };
        },
        errorElement: <ErrorPage />
    },
    {
        path: '/login',
        element: <LoginPage />,
        loader: loginLoader,
        action: loginAction,
        errorElement: <ErrorPage />
    },
    {
        path: '/logout',
        element: <LogoutPage />,
        loader: logoutLoader,
        errorElement: <ErrorPage />
    },
    {
        path: '/password_change/:token?',
        element: <PasswordChangePage />,
        action: passwordChangeAction,
        errorElement: <ErrorPage />
    },
    {
        path: '/dashboard',
        element: <DashboardPage />,
        loader: dashboardLoader,
        errorElement: <ErrorPage />
    },
    {
        path: '/audit/view',
        element: <AuditViewPage />,
        loader: auditViewLoader,
        errorElement: <ErrorPage />
    }
]);

document.documentElement.setAttribute(
    'data-bs-theme',
    window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light'
);

const root = ReactDOM.createRoot(document.getElementById('root'));
root.render(
    <React.StrictMode>
        <RouterProvider router={router} />
    </React.StrictMode>
);
