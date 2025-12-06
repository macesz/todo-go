import React, { useState } from 'react'
import { X, Plus } from 'lucide-react';


export default function LabelPopUP({ isOpen, onClose, onUpdate, labels = [], }) {

    const [inputValue, setInputValue] = useState("");

    if (!isOpen) return null;

    // Handlers

    const handleAddLabel = (e) => {
        if (e.key === 'Enter') {
            e.preventDefault();
            const newLabel = inputValue.trim();
            if (newLabel && !labels.includes(newLabel)) {
                onUpdate([...labels, newLabel]);
                setInputValue("");
            }
        }
    };

    const handleRemoveLabel = (labelToRemove) => {
        const updateLabels = labels.filter(label => label !== labelToRemove);
        onUpdate(updateLabels);
    }

    return (
        <>
            {/* A. The Invisible Backdrop (Closes menu when clicking outside) */}
            <div
                className="fixed inset-0 z-40"
                onClick={() => onClose()}
            ></div>

            {/* B. The Actual Popup */}
            <div className="left-0 mb-2 z-50 bg-white p-3 rounded-xl shadow-xl border border-gray-100 w-72 cursor-default animate-in fade-in zoom-in-95 duration-200">

                <div className='space-y-3'>
                    <h4 className='text-xs font-semibold text-gray-500 uppercase tracking-wider'>
                        Manage Labels
                    </h4>
                    {/* A. Active Labels List (Chips) */}
                    <div className='flex flex-wrap gap-2'>
                        {labels.map((label, index) => (
                            <span
                                key={index}
                                className="inline-flex items-center px-2 py-1 rounded-md text-xs font-medium bg-gray-100 text-gray-700 border border-gray-200"
                            >
                                {label}
                                <button
                                    onClick={() => handleRemoveLabel(label)}
                                    className="ml-2 text-gray-400 hover:text-purple-500 focus:outline-none"
                                    aria-label={`Remove ${label}`}
                                >
                                    <X size={12} />
                                </button>
                            </span>
                        ))}

                        {labels.length === 0 && (
                            <p className="text-xs text-gray-400">No labels yet..</p>
                        )}
                    </div>
                    {/* B. Input Field to Add New Label */}
                    <div className="relative flex items-center">
                        <input
                            type="text"
                            autoFocus
                            className="w-full pl-2 pr-8 py-1.5 text-sm border border-gray-200 rounded-lg focus:outline-none focus:border-blue-500 focus:ring-1 focus:ring-blue-500 transition-all placeholder-gray-400"
                            placeholder="Type & press Enter..."
                            value={inputValue}
                            onChange={(e) => setInputValue(e.target.value)}
                            onKeyDown={handleAddLabel}
                        />
                        <div className="absolute right-2 text-gray-400 pointer-events-none">
                            <Plus size={14} />
                        </div>
                    </div>
                </div>
            </div>
        </>
    );
};
