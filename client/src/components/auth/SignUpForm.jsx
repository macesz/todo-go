import React from 'react'
import { Mail, Lock, User } from 'lucide-react';
import { FaFacebookF, FaGoogle, FaLinkedinIn } from 'react-icons/fa';
import InputIcon from '../ui/InputIcon';
import SocialButton from '../ui/SocialButton';


export default function SignUpForm() {

    const handleSubmit = async (e) => {
        e.preventDefault();
        // TODO: Fetch to Go backend (e.g., POST /user)
        // const response = await fetch('http://localhost:8080/user', { method: 'POST', body: ... });
        // Handle response
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
                <InputIcon icon={<User size={18} />} type="text" placeholder="Name" />
                <InputIcon icon={<Mail size={18} />} type="email" placeholder="Email" />
                <InputIcon icon={<Lock size={18} />} type="password" placeholder="Password" />
            </div>

            <button type="submit" className="bg-purple-600 text-white font-bold py-3 px-12 rounded-full tracking-wider uppercase text-xs transition-transform active:scale-95 shadow-lg hover:bg-purple-700 hover:shadow-purple-500/30">
                Sign Up
            </button>
        </form>

    )
}
