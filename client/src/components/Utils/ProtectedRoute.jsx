import { useAuth } from '../../Context/AuthContext.jsx';
import AuthFailureModal from './AuthFailureModal.jsx';
import { Navigate, useLocation } from 'react-router-dom';


const ProtectedRoute = ({ children }) => {
  const { user } = useAuth();
  const location = useLocation();

  if (!user) {
    return Navigate({ to: '/auth', state: { from: location }, replace: true });

  }

  return (
    children
  );
};

export default ProtectedRoute;
