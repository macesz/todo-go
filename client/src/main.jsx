import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import { createBrowserRouter, RouterProvider } from 'react-router-dom'

import './index.css'
import MainLayout from './components/layout/MainLayout.jsx'
import ProtectedRoute from './components/util/ProtectedRoute.jsx';
import HomePage from './pages/HomePage.jsx'
import { Navigate } from 'react-router-dom';

import TodoCard from './pages/TodoCard.jsx'
import AuthPage from './pages/AuthPage.jsx'
import AuthProvider from './context/AuthContext.jsx';

const router = createBrowserRouter([
  // Public Routes (No Layout, No Protection)
  {
    path: "/auth",
    element: <AuthPage />
  },

  // Protected Routes (Wrapped in Layout)
  {
    path: "/",
    element: (
      <ProtectedRoute>
        <MainLayout />
      </ProtectedRoute>
    ),
    children: [
      {
        index: true, // Matches path: "/"
        element: <HomePage />, 
      },
      {
        path: "lists/:listId", 
        element: <TodoCard />,
      },
      // Catch all unknown routes and redirect to home
      {
        path: "*",
        element: <Navigate to="/" replace />
      }
    ]
  },
]);

createRoot(document.getElementById('root')).render(
  <StrictMode>
    <AuthProvider>
    <RouterProvider router={router} />
    </AuthProvider>
  </StrictMode>
)
