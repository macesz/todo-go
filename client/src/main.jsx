import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import { createBrowserRouter, RouterProvider } from 'react-router-dom'

import './index.css'
import ProtectedRoute from './components/Utils/ProtectedRoute.jsx'
import { Navigate } from 'react-router-dom';

import AuthPage from './Pages/AuthPage.jsx'
import AuthProvider from './Context/AuthContext.jsx';
import MainLayout from './Layouts/MainLayout.jsx';
import TodoCard from './Pages/TodoCard.jsx';
import HomePage from './pages/HomePage.jsx'
import { ListProvider } from './Context/ListContext.jsx'

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
      <ListProvider >
        <RouterProvider router={router} />
      </ListProvider>
    </AuthProvider>
  </StrictMode>
)
