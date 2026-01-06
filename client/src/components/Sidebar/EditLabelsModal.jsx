import React, { useState, } from "react";
import { X, Plus, Check, } from "lucide-react";
import EditLabelRow from "./EditLabelRow";

export default function EditLabelsModal({ isOpen, onClose, labels, onUpdate, onDelete }) {
    const [focusedId, setFocusedId] = useState(null);

    if (!isOpen) return null;


    return (
        <div className="fixed inset-0 z-50 flex items-center justify-center p-4">
            {/* Backdrop: The dark blurry background */}
            <div
                className="absolute inset-0 bg-black/20 backdrop-blur-sm transition-opacity"
                onClick={onClose}
            ></div>

            {/* Modal Box */}
            <div className="relative bg-white w-full max-w-xs md:max-w-sm rounded-lg shadow-xl overflow-hidden">

                {/* Header */}
                <div className="p-4 border-b border-gray-100">
                    <h3 className="text-lg font-medium text-gray-800">Edit labels</h3>
                </div>

                {/* List of Labels */}
                <div className="p-2 max-h-[60vh] overflow-y-auto custom-scrollbar">
                    <div className="space-y-1">
                        {/* We loop through the labels we found in our lists */}
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

                {/* Footer: Just a button to close the modal */}
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
