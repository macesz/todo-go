import React, { useState, useRef, useEffect, useCallback } from 'react';
import { 
    CheckSquare, 
    Plus, 
    X, 
    Palette, 
    Tag, 
    MoreVertical, 
    Undo, 
    Redo 
} from 'lucide-react';
import { COLOR_PALETTE } from '../../data/ColorPalette'; 
import ColorPopUp from '../Ui/ColorPopUp';
import LabelPopUP from '../Ui/LabelPopUP';

export default function CreateList({ onSave }) {
    // --- State ---
    const [isExpanded, setIsExpanded] = useState(false);
    
    // Form Data
    const [title, setTitle] = useState("");
    const [listItems, setListItems] = useState([]); 
    const [newTodoTitle, setNewTodoTitle] = useState("");
    
    // Meta Data
    const [selectedColor, setSelectedColor] = useState("default");
    const [selectedLabels, setSelectedLabels] = useState([]);

    // UI State
    const [activeMenu, setActiveMenu] = useState(null); // 'palette', 'labels', etc.
    const containerRef = useRef(null);


    // Reset form to initial state
    const resetForm = () => {
        setTitle("");
        setListItems([]);
        setNewTodoTitle("");
        setSelectedColor("default");
        setSelectedLabels([]);
        setIsExpanded(false);
        setActiveMenu(null);
    };

    const handleAddItem = () => {
        if (newTodoTitle.trim()) {
            const newItem = {
                id: Date.now(),
                title: newTodoTitle,
                done: false
            };
            setListItems([...listItems, newItem]);
            setNewTodoTitle("");
        }
    };

    const handleDeleteItem = (id) => {
        setListItems(listItems.filter(item => item.id !== id));
    };

    const handleSave = useCallback (() => {
        // Only save if there is content
        if (title.trim() || listItems.length > 0 || newTodoTitle.trim()) {
            
            // If user typed in the input but didn't press enter, add it now
            let finalItems = [...listItems];
            if (newTodoTitle.trim()) {
                finalItems.push({
                    title: newTodoTitle,
                    done: false
                });
            }

            const newList = {
                title: title.trim() || "Untitled List", // Fallback title
                items: finalItems,
                color: selectedColor,
                labels: selectedLabels,
                done: false
            };

            onSave(newList);
        }
        resetForm();
    }, [title, listItems, newTodoTitle, selectedColor, selectedLabels, onSave]);

    // --- Event Handlers ---

    // Handle "Enter" key in the add item input
    const handleKeyDown = (e) => {
        if (e.key === 'Enter') {
            e.preventDefault();
            handleAddItem();
        }
    };

    useEffect(() => {
        const handleClickOutside = (event) => {
            if (containerRef.current && !containerRef.current.contains(event.target)) {
                if (isExpanded) {
                    handleSave();
                }
            }
        };

        document.addEventListener("mousedown", handleClickOutside);
        return () => document.removeEventListener("mousedown", handleClickOutside);
    }, [isExpanded, handleSave]);


    // --- RENDER ---

    // 1. COLLAPSED VIEW
    if (!isExpanded) {
        return (
            <div className="w-full flex justify-center mb-8 px-4">
                <div 
                    onClick={() => setIsExpanded(true)}
                    className="w-full max-w-2xl bg-white shadow-sm border border-gray-200 rounded-lg flex items-center justify-between p-3 cursor-text hover:shadow-md transition-shadow"
                >
                    <span className="text-gray-500 font-medium pl-2">Take a note...</span>
                    
                    <div className="flex gap-2">
                         <button className="p-2 text-gray-500 hover:bg-gray-100 rounded-full transition-colors" title="New List">
                            <CheckSquare size={20} />
                        </button>
                    </div>
                </div>
            </div>
        );
    }

    // 2. EXPANDED VIEW
    return (
        <div className="w-full flex justify-center mb-8 px-4 relative z-10">
            <div 
                ref={containerRef}
                className='w-full max-w-2xl rounded-lg shadow-xl border border-gray-200 flex flex-col transition-colors duration-300 bg-white'
            >
                {/* A. Title Input */}
                <div className="p-4 pb-2">
                    <input 
                        type="text" 
                        placeholder="Title" 
                        value={title}
                        onChange={(e) => setTitle(e.target.value)}
                        className="w-full bg-transparent text-lg font-medium placeholder-gray-500 text-gray-800 outline-none"
                        autoFocus
                    />
                </div>

                {/* B. Existing List Items */}
                <div className="px-4 space-y-1">
                    {listItems.map((item) => (
                        <div key={item.id} className="flex items-center gap-2 group">
                            <div className={`w-3.5 h-3.5 border-2 border-gray-400 rounded-sm`}></div>
                            <span className="text-sm text-gray-700 flex-1 py-1 border-b border-transparent focus:border-gray-200">
                                {item.title}
                            </span>
                            <button 
                                onClick={() => handleDeleteItem(item.id)}
                                className="opacity-0 group-hover:opacity-100 p-1 text-gray-400 hover:text-gray-700"
                            >
                                <X size={14} />
                            </button>
                        </div>
                    ))}
                </div>

                {/* C. Add New Item Input */}
                <div className="px-4 py-2 flex items-center gap-2 border-b border-transparent focus-within:border-gray-200">
                    <Plus size={16} className="text-gray-400" />
                    <input 
                        type="text"
                        placeholder="List item"
                        value={newTodoTitle}
                        onChange={(e) => setNewTodoTitle(e.target.value)}
                        onKeyDown={handleKeyDown}
                        className="flex-1 bg-transparent text-sm placeholder-gray-500 text-gray-800 outline-none py-1"
                    />
                </div>

                {/* D. Bottom Toolbar */}
                <div className="p-2 flex items-center justify-between mt-2">
                    
                    {/* Left: Action Icons */}
                    <div className="flex items-center gap-1">
                        
                        {/* Palette */}
                        <div className="relative">
                            <button 
                                onClick={() => setActiveMenu(activeMenu === 'palette' ? null : 'palette')}
                                className="p-2 text-gray-600 hover:bg-black/5 rounded-full transition-colors" 
                                title="Background options"
                            >
                                <Palette size={18} />
                            </button>
                            <div className="absolute top-full left-0 mt-2 z-20">
                                <ColorPopUp 
                                    isOpen={activeMenu === 'palette'}
                                    onClose={() => setActiveMenu(null)}
                                    selectedColor={selectedColor}
                                    onSelect={setSelectedColor}
                                    palette={COLOR_PALETTE}
                                />
                            </div>
                        </div>

                        {/* Labels */}
                        <div className="relative">
                             <button 
                                onClick={() => setActiveMenu(activeMenu === 'labels' ? null : 'labels')}
                                className="p-2 text-gray-600 hover:bg-black/5 rounded-full transition-colors"
                                title="Add label"
                             >
                                <Tag size={18} />
                            </button>
                            <div className="absolute top-full left-0 mt-2 z-20">
                                <LabelPopUP 
                                    isOpen={activeMenu === 'labels'}
                                    onClose={() => setActiveMenu(null)}
                                    selectedLabels={selectedLabels}
                                    onUpdate={(newLabels) => setSelectedLabels(newLabels)}
                                />
                            </div>
                        </div>
                    </div>

                    {/* Right: Close/Save Button */}
                    <button 
                        onClick={handleSave}
                        className="px-6 py-2 text-sm font-medium text-gray-800 hover:bg-black/5 rounded-md transition-colors"
                    >
                        Save
                    </button>
                </div>
            </div>
        </div>
    );
}