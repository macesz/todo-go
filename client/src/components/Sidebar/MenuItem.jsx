import React from 'react'

export default function MenuItem({ icon, label, active, onClick }) {
    return (
        <li
            onClick={onClick}
            className={`flex items-center gap-3 p-2 rounded-lg cursor-pointer transition-colors 
      ${active ? 'bg-white shadow-sm text-gray-900' : 'text-gray-600 hover:bg-gray-100 hover:text-gray-900'}
    `}
        >
            {React.cloneElement(icon, { size: 18, className: active ? 'text-purple-600' : 'text-gray-500' })}
            <span className="text-sm font-medium">{label}</span>
        </li>)
}
