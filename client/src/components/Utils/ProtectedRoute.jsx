import { useAuth } from '../../Context/AuthContext.jsx';
import AuthFailureModal from './AuthFailureModal.jsx';


const ProtectedRoute = ({ children }) => {
  const { user } = useAuth();

  return (
    <>
      {user ? children : <AuthFailureModal />}
    </>
  );
};

export default ProtectedRoute;
