import React from 'react'
import { useNavigate } from 'react-router-dom';
import { Mail, Lock } from 'lucide-react';
import { FaFacebookF, FaGoogle, FaLinkedinIn } from 'react-icons/fa';
import SocialButton from '../ui/SocialButton.jsx';
import InputIcon from '../ui/InputIcon.jsx';
import { useAuth } from '../../context/AuthContext.jsx';



export default function LoginForm() {

    const { login } = useAuth();
    const navigate = useNavigate();

    const handleLogin = async (e) => {
        e.preventDefault();


        // TODO: Fetch to Go backend (e.g., POST /login)
        // const response = await fetch('http://localhost:8080/login', { method: 'POST', body: ... });
        // const data = await response.json();
        // if (data.token) login(data.token, data.user);

        // For now, i simulate a successful login:
        login('fake-token', { name: 'User' }); // Call the login function from context

        navigate('/'); // Redirect to home page after login
    }

    return (

        <form className="w-full flex flex-col items-center text-center" onSubmit={handleLogin}>
            <h1 className="font-bold text-3xl text-purple-900 mb-6">Sign in to TodoApp</h1>

            <div className="flex gap-4 mb-6">
                <SocialButton icon={<FaFacebookF size={20} />} />
                <SocialButton icon={<FaGoogle size={20} />} />
                <SocialButton icon={<FaLinkedinIn size={20} />} />
            </div>

            <span className="text-sm text-gray-400 mb-6">or use your email account</span>

            <div className="w-full space-y-4 mb-4">
                <InputIcon icon={<Mail size={18} />} type="email" placeholder="Email" />
                <InputIcon icon={<Lock size={18} />} type="password" placeholder="Password" />
            </div>

            <a href="#" className="text-xs text-gray-500 mb-8 hover:text-purple-700 border-b border-transparent hover:border-purple-700 transition-colors">
                Forgot your password?
            </a>

            <button type="submit" className="bg-purple-600 text-white font-bold py-3 px-12 rounded-full tracking-wider uppercase text-xs transition-transform active:scale-95 shadow-lg hover:bg-purple-700 hover:shadow-purple-500/30">
                Sign In
            </button>
        </form>
    );
}
