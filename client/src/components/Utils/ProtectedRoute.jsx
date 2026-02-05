import { useAuth } from '../../context/AuthContext.jsx';
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
