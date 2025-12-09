import React from 'react'

export default function InputIcon({ icon, type, placeholder, value, onChange, }) {
    return (
        <div className="relative w-full">
            <div className="absolute inset-y-0 left-0 pl-4 flex items-center pointer-events-none text-gray-400">
                {icon}
            </div>
            <input
                className="bg-gray-100 w-full pl-10 pr-4 py-3 rounded-lg text-sm outline-none focus:ring-2 focus:ring-purple-500/50 transition-all placeholder-gray-400 text-gray-700"
                type={type}
                placeholder={placeholder}
                value={value}
                onChange={onChange}
            />
      
        </div>
    )
}
