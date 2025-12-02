import { useNavigate } from "react-router-dom";

const AuthFailureModal = () => {
    const navigate = useNavigate();

    const handleClick = () => {
        navigate('/login');
    }

    return (
        <div className="fixed inset-0 bg-[#394a51] bg-opacity-50 flex items-center justify-center">
            <div className="bg-white p-6 rounded-lg shadow-lg max-w-sm text-center">
                <h2 className="text-lg font-bold mb-2">Authentication error</h2>
                <p>Unable to authenticate. Please enter your login credentials.</p>
                <button
                    onClick={handleClick}
                    className="mt-4 px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600"
                >
                    Sign in
                </button>
            </div>
        </div>
    );
};

export default AuthFailureModal;
