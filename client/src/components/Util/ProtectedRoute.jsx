import React from 'react';
import { useAuth } from '../../context/AuthContext.jsx';
import AuthFailureModal from '../auth/AuthFailureModal.jsx';    

const ProtectedRoute = ({ children }) => {
  const { user } = useAuth();

  return (
    <>
      {user ? children : <AuthFailureModal />}
    </>
  );
};

export default ProtectedRoute;
