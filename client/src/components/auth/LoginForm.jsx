import React, { useState } from 'react'
import { useNavigate } from 'react-router-dom';
import { Mail, Lock, Eye, EyeOff } from 'lucide-react';
import { FaFacebookF, FaGoogle, FaLinkedinIn } from 'react-icons/fa';
import SocialButton from '../ui/SocialButton.jsx';
import InputIcon from '../ui/InputIcon.jsx';
import { useAuth } from '../../Context/AuthContext.jsx';
import { loginUser } from "../Services/apiService.js";



export default function LoginForm() {

    const [email, setEmail] = useState("");
    const [password, setPassword] = useState("");
    const [showPassword, setShowPassword] = useState(false);


    const { login } = useAuth();
    const navigate = useNavigate();

    const togglePasswordVisibility = () => setShowPassword(!showPassword);


    const handleLogin = async (e) => {
        e.preventDefault();

        try {
            const user = await loginUser(email, password);
            login(user);
            navigate('/');
        } catch (error) {
            console.error('Login failed:', error);
            navigate('/login');
        }
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
                <InputIcon
                    icon={<Mail size={18} />}
                    type="email"
                    placeholder="Email"
                    value={email}
                    onChange={(e) => setEmail(e.target.value)}
                />
                <InputIcon
                    icon={<Lock size={18} />}
                    type={showPassword ? "text" : "password"}
                    placeholder="Password"
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                />
                <button
                    type="button"
                    onClick={togglePasswordVisibility}
                    className="absolute right-3 top-1/2 -translate-y-1/2 text-gray-400 hover:text-purple-600 transition-colors cursor-pointer outline-none"
                >
                    {showPassword ? <EyeOff size={18} /> : <Eye size={18} />}
                </button>
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
