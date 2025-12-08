import React, { useState,  } from "react";
import { X, Plus, Check, } from "lucide-react";

export default function EditLabelsModal({ isOpen, onClose, labels, onAdd, onUpdate, onDelete }) {
    const [newLabelName, setNewLabelName] = useState("");
    const [focusedId, setFocusedId] = useState(null);

    if (!isOpen) return null;

    // --- Handlers ---

    const handleCreate = () => {
        if (newLabelName.trim()) {
            onAdd(newLabelName);
            setNewLabelName("");
        }
    };

    return (
        <div className="fixed inset-0 z-50 flex items-center justify-center p-4">
            {/* Backdrop */}
            <div
                className="absolute inset-0 bg-black/20 backdrop-blur-sm transition-opacity"
                onClick={onClose}
            ></div>

            {/* Modal Content */}
            <div className="relative bg-white w-full max-w-xs md:max-w-sm rounded-lg shadow-xl overflow-hidden animation-fade-in">

                {/* Header */}
                <div className="p-4 border-b border-gray-100">
                    <h3 className="text-lg font-medium text-gray-800">Edit labels</h3>
                </div>

                <div className="p-2 max-h-[60vh] overflow-y-auto custom-scrollbar">

                    {/* 1. Create New Label Row */}
                    <div className="flex items-center gap-2 p-2 mb-2 transition-colors">
                        {/* Left Icon: Clear (X) or Plus (+) depending on focus/content */}
                        <button
                            onClick={() => setNewLabelName("")}
                            className={`p-2 rounded-full text-gray-500 hover:bg-gray-100 transition-colors ${newLabelName ? 'opacity-100' : 'opacity-50 pointer-events-none'}`}
                        >
                            {newLabelName ? <X size={18} /> : <Plus size={18} />}
                        </button>

                        <input
                            type="text"
                            placeholder="Create new label"
                            className="flex-1 bg-transparent border-b border-transparent focus:border-gray-300 outline-none text-sm py-1 placeholder-gray-400 text-gray-700"
                            value={newLabelName}
                            onChange={(e) => setNewLabelName(e.target.value)}
                            onKeyDown={(e) => e.key === 'Enter' && handleCreate()}
                        />

                        {/* Right Icon: Checkmark (Confirm) */}
                        <button
                            onClick={handleCreate}
                            disabled={!newLabelName.trim()}
                            className={`p-2 rounded-full transition-colors ${newLabelName.trim() ? 'text-green-600 hover:bg-green-50' : 'text-gray-300'}`}
                        >
                            <Check size={18} />
                        </button>
                    </div>

                    {/* 2. Existing Labels List */}
                    <div className="space-y-1">
                        {labels.map((label) => (
                            <EditLabelRow
                                key={label.id}
                                label={label}
                                onUpdate={onUpdate}
                                onDelete={onDelete}
                                isFocused={focusedId === label.id}
                                setFocused={setFocusedId}
                            />
                        ))}
                    </div>
                </div>

                {/* Footer */}
                <div className="p-3 border-t border-gray-100 flex justify-end bg-gray-50">
                    <button
                        onClick={onClose}
                        className="px-4 py-1.5 bg-white border border-gray-200 text-gray-600 rounded hover:bg-gray-100 text-sm font-medium transition-colors"
                    >
                        Done
                    </button>
                </div>
            </div>
        </div>
    );
}
