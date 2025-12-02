import { Trash2 } from 'lucide-react';
import { useSortable } from '@dnd-kit/sortable';
import { CSS } from '@dnd-kit/utilities';

const TaskItem = ({ todoItem, onToggle, onDelete, checkboxColor, hoverColor }) => {

    const { attributes, listeners, setNodeRef, transform, transition } = useSortable({ id: todoItem.id });

    const style = {
        transition,
        transform: CSS.Transform.toString(transform),
    };

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
                <span className={`block text-sm truncate ${todoItem.completed
                    ? 'line-through text-gray-400'
                    : 'text-gray-700 font-medium'
                    }`}>
                    {todoItem.title}
                </span>
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

export default TaskItem;