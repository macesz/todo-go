import React, { useState } from 'react'
import { Mail, Lock, User, Eye, EyeOff } from 'lucide-react';
import { FaFacebookF, FaGoogle, FaLinkedinIn } from 'react-icons/fa';
import InputIcon from '../ui/InputIcon';
import SocialButton from '../ui/SocialButton';
import { registerUser } from "../Services/apiService.js";
import { useNavigate } from 'react-router-dom';


export default function SignUpForm() {
    const navigate = useNavigate();

    const [name, setName] = useState("");
    const [email, setEmail] = useState("");
    const [password, setPassword] = useState("");

    const [showPassword, setShowPassword] = useState(false);

    const togglePasswordVisibility = () => setShowPassword(!showPassword);

    const handleSubmit = async (e) => {
        e.preventDefault();

        const userData = {
            name,
            email,
            password
        }
        await registerUser(userData);

        navigate('/login');
    };

    return (
        <form className="w-full flex flex-col items-center text-center" onSubmit={handleSubmit}>
            <h1 className="font-bold text-3xl text-purple-900 mb-6">Create Account</h1>

            <div className="flex gap-4 mb-6">
                <SocialButton icon={<FaFacebookF size={20} />} />
                <SocialButton icon={<FaGoogle size={20} />} />
                <SocialButton icon={<FaLinkedinIn size={20} />} />

            </div>

            <span className="text-sm text-gray-400 mb-6">or use your email for registration</span>

            <div className="w-full space-y-4 mb-8">
                <InputIcon
                    icon={<User size={18} />}
                    type="text"
                    placeholder="Name"
                    value={name}
                    onChange={(e) => setName(e.target.value)}
                />
                <InputIcon
                    icon={<Mail size={18} />}
                    type="email"
                    placeholder="Email"
                    value={email}
                    onChange={(e) => setEmail(e.target.value)}
                />
                <InputIcon icon={<Lock size={18} />}
                    type="password"
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

            <button type="submit" className="bg-purple-600 text-white font-bold py-3 px-12 rounded-full tracking-wider uppercase text-xs transition-transform active:scale-95 shadow-lg hover:bg-purple-700 hover:shadow-purple-500/30">
                Sign Up
            </button>
        </form>

    )
}
