import React from 'react';

export default function SocialButton({ icon }) {
    return (
        <a href="#" className="w-10 h-10 border border-gray-300 rounded-full flex justify-center items-center text-gray-600 hover:bg-gray-100 hover:text-purple-600 transition-all duration-300">
            {icon}
        </a>)
};
