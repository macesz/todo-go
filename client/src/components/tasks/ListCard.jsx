import React, { useState } from 'react';
import { Plus } from 'lucide-react';
import TaskItem from './TaskItem';
import { COLOR_PALETTE } from '../../data/ColorPalette';
import {
    DndContext,
    KeyboardSensor,
    PointerSensor,
    useSensor,
    useSensors,
    closestCorners,
} from "@dnd-kit/core";
import {
    SortableContext,
    verticalListSortingStrategy,
    arrayMove,
    sortableKeyboardCoordinates
} from '@dnd-kit/sortable';

const ListCard = ({ list }) => {
    const [todos, setTodos] = useState(list.items);
    const [inputValue, setInputValue] = useState("");

    const theme = COLOR_PALETTE[list.color] || COLOR_PALETTE.default;

    const handleToggle = (id) => {
        setTodos(todos.map(todo =>
            todo.id === id ? { ...todo, completed: !todo.completed } : todo
        ));
    };

    const handleDelete = (id) => {
        setTodos(todos.filter(todo => todo.id !== id));
    };

    const handleAdd = (e) => {
        if (e.key === 'Enter' && inputValue.trim()) {
            setTodos([...todos, { id: Date.now(), title: inputValue, completed: false }]);
            setInputValue("");
        }
    };


    // DnD Kit Setup activationConstraint: tells dnd-kit: "Wait until the user moves 8px before assuming it's a drag."
    // This allows clicks on checkboxes/buttons to pass through instantly.
    const sensors = useSensors(
        useSensor(PointerSensor, { activationConstraint: { distance: 8 } }),
        useSensor(KeyboardSensor, {
            coordinateGetter: sortableKeyboardCoordinates,
        })
    );

    const getTaskPos = (id) => {
        return todos.findIndex((t) => t.id === id);
    }

    const handleDragEnd = (event) => {
        const { active, over } = event;

        if (active.id === over.id) return;

        setTodos((todos) => {
            const originalPos = getTaskPos(active.id);
            const newPos = getTaskPos(over.id);

            return arrayMove(todos, originalPos, newPos);
        });
    };

    // Sort: Active first, Completed last
    const activeTodos = todos.filter(todo => !todo.completed);
    const completedTodos = todos.filter(todo => todo.completed);

    return (
        <DndContext
            sensors={sensors}
            collisionDetection={closestCorners}
            onDragEnd={handleDragEnd}
        >
            <div className={`break-inside-avoid mb-6 rounded-2xl shadow-[0_4px_20px_-4px_rgba(0,0,0,0.1)] bg-light hover:shadow-xl transition-shadow duration-300 flex flex-col overflow-hidden border border-gray-100`}>
                <div className="card-body p-0">

                    <ul className="list w-full">

                        {/* Header Color Strip (Optional, based on label color) */}
                        <div className={`h-1.5 w-full ${theme.bar}`}></div>



                        {/* --- HEADER: Title & Labels --- */}
                        <li className="p-4 pb-1 flex flex-col gap-2">
                            <div className="flex justify-between items-start">
                                <span className="text-sm uppercase font-bold opacity-70 tracking-wide text-gray-800">
                                    <h3 className="font-bold text-lg text-gray-800">{list.title}</h3>
                                    <p className="text-xs text-gray-400 font-medium uppercase tracking-wider mt-1">
                                        {list.labels.join(' • ')}
                                    </p>

                                </span>
                            </div>
                        </li>

                        {/* Active Tasks (Draggable) */}
                        <ul className='space-y-1 mb-4'>
                            <SortableContext items={activeTodos} strategy={verticalListSortingStrategy}>
                                {activeTodos.map(todo => (
                                    <TaskItem
                                        key={todo.id}
                                        todoItem={todo}
                                        checkboxColor={theme.checkbox}
                                        hoverColor={theme.hover}
                                        onToggle={handleToggle}
                                        onDelete={handleDelete}
                                    />
                                ))}
                            </SortableContext>
                        </ul>

                        {activeTodos.length > 0 && completedTodos.length > 0 && (
                            <div className="divider my-0 px-4 text-xs text-gray-400 font-semibold uppercase">
                                Completed
                            </div>)}

                        <ul className='space-y-1 mb-4'>
                            {completedTodos.map(item => (
                                <TaskItem
                                    key={item.id}
                                    todoItem={item}
                                    onToggle={handleToggle}
                                    onDelete={handleDelete}
                                    checkboxColor={theme.checkbox}
                                    hoverColor={theme.hover}
                                />
                            ))}
                        </ul>


                        {/* --- FOOTER: Input Field --- */}
                        <li className={`flex items-center ${theme.hover} rounded-lg shadow-sm gap-3 mt-auto m-4 pt-3 px-6 pb-4`}>
                            <Plus size={18} className="text-accent" />

                            <div className="w-full">
                                <input
                                    type="text"
                                    className={`bg-transparent outline-none w-full text-sm ${theme.placeholder}`}
                                    placeholder="Add list item..."
                                    value={inputValue}
                                    onChange={(e) => setInputValue(e.target.value)}
                                    onKeyDown={handleAdd}
                                />
                            </div>
                        </li>
                    </ul>
                </div>
            </div>
        </DndContext>
    );
};

export default ListCard;