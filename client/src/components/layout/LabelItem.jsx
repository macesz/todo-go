import React from 'react'

export default function LabelItem({ label, onClick }) {
    return (
        <li
            onClick={onClick}
            className="flex items-center gap-3 p-2 rounded-lg cursor-pointer text-gray-600 hover:bg-gray-100 transition-colors"
        >
            <div className={`w-3 h-3 rounded-full ${label.color}`} />
            <span className="text-sm font-medium">{label.name}</span>
        </li>)
}
