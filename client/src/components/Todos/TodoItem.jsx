import { Trash2 } from 'lucide-react';
import { useSortable } from '@dnd-kit/sortable';
import { CSS } from '@dnd-kit/utilities';
import { useEffect, useRef, useState } from 'react';

export default function TaskItem({ todoItem, onToggle, onDelete, onEdit, checkboxColor, hoverColor }) {

    const [isEditing, setIsEditing] = useState(false);
    const [editValue, setEditValue] = useState(todoItem.title);
    const inputRef = useRef(null);

    const { attributes, listeners, setNodeRef, transform, transition, isDragging } = useSortable({ id: todoItem.id });

    const style = {
        transition,
        transform: CSS.Transform.toString(transform),
        opacity: isDragging ? 0.5 : 1,
        zIndex: isDragging ? 50 : 'auto',
    };

    useEffect(() => {
        if (isEditing && inputRef.current) {
            inputRef.current.focus();
        }
    }, [isEditing]);

    const handleSave = () => {
        if (editValue.trim()) {
            console.log("Saving edit: ", editValue.trim());

            onEdit(todoItem.id, editValue.trim());
        } else {
            setEditValue(todoItem.title); // Revert to original if empty
        }
        setIsEditing(false);
    }

    const handleKeyDown = (e) => {

        if (e.key === ' ' || e.key === 'Enter') {
            e.stopPropagation();
        }

        if (e.key === 'Enter') {
            e.preventDefault();
            handleSave();
        } else if (e.key === 'Escape') {
            setEditValue(todoItem.title);
            setIsEditing(false);
        }
    }

    return (
        <li
            ref={setNodeRef}
            style={style}
            {...attributes}
            {...listeners}
            className={`group flex items-center gap-3 p-2 ${hoverColor} rounded-lg shadow-sm mt-auto m-4 pt-3 transition-colors cursor-pointer min-h-[40px]`}>

            {/* 1. Checkbox (Fixed width to prevent squishing) */}
            <div className="flex-shrink-0">
                <input
                    id={`checkbox-${todoItem.id}`}
                    type="checkbox"
                    className={`checkbox checkbox-xs rounded-sm border-2 transition-all ${checkboxColor}`}
                    checked={todoItem.completed}
                    onChange={() => onToggle(todoItem.id)}
                />

            </div>

            {/* 2. Text (Takes remaining space) */}
            <div className="flex-grow min-w-0">
                {isEditing ? (
                    <input
                        ref={inputRef}
                        type="text"
                        value={editValue}
                        onChange={(e) => setEditValue(e.target.value)}
                        onBlur={handleSave}
                        onKeyDown={handleKeyDown}
                        // IMPORTANT: Stop propagation so dnd-kit doesn't think we are dragging
                        onPointerDown={(e) => e.stopPropagation()}
                        className="w-full bg-transparent outline-none text-sm text-gray-800 font-medium p-0 border-b border-blue-400"
                    />
                ) : (
                    <span
                        onClick={(e) => {
                            e.stopPropagation(); // Prevent drag start on click
                            setIsEditing(true);
                        }}
                        className={`block text-sm truncate cursor-text ${todoItem.completed
                            ? 'line-through text-gray-400'
                            : 'text-gray-700 font-medium'
                            }`}
                    >
                        {todoItem.title}
                    </span>
                )}
            </div>

            {/* Column 3: Action Buttons (Hidden until hover) */}
            <button
                onClick={() => onDelete(todoItem.id)}
                className="btn btn-square btn-ghost btn-sm opacity-0 group-hover:opacity-100 transition-opacity"
            >
                <Trash2 className="size-4 text-gray-950" />
            </button>
        </li>
    );
};

