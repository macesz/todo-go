import {useRef, useState } from "react";
import { Trash2, Pencil, Check } from "lucide-react";


export default function EditLabelRow({ label, onUpdate, onDelete, isFocused, setFocused }) {
    const [tempName, setTempName] = useState(label.name);
    const inputRef = useRef(null);

    const handleSave = () => {
        if (tempName.trim() && tempName !== label.name) {
            onUpdate(label.id, tempName.trim());
        } else {
            setTempName(label.name)
        }
        setFocused(null);
    };

    const handleFocus = () => {
        setFocused(label.id);
        inputRef.current?.focus();
    };

    return (
        <div className={`group flex items-center gap-2 p-1 px-2 rounded-full hover:bg-gray-50 transition-colors ${isFocused ? 'bg-gray-50' : ''}`}>

            {/* Left: Delete Button */}
            <button
                onClick={() => onDelete(label.id)}
                className="p-2 text-gray-500 hover:bg-red-100 hover:text-red-500 rounded-full transition-colors"
                title="Delete label"
            >
                <Trash2 size={16} />
            </button>

            {/* Middle: Input */}
            <input
                ref={inputRef}
                type="text"
                value={tempName}
                onChange={(e) => setTempName(e.target.value)}
                onFocus={() => setFocused(label.id)}
                onBlur={handleSave}
                onKeyDown={(e) => e.key === 'Enter' && e.target.blur()}
                className="flex-1 bg-transparent border-b border-transparent focus:border-gray-300 outline-none text-sm text-gray-700 py-1 font-medium"
            />

            {/* Right: Pencil (Edit) / Check (Save) */}
            <button
                onClick={isFocused ? () => inputRef.current.blur() : handleFocus}
                className="p-2 text-gray-400 hover:text-gray-700 rounded-full transition-colors"
            >
                {isFocused ? <Check size={16} /> : <Pencil size={16} />}
            </button>
        </div>
    );
}