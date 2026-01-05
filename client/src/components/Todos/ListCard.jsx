import React, { useState } from 'react';
import { Plus, PaletteIcon, TagsIcon, Trash2 } from 'lucide-react';
import TaskItem from './TodoItem';
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
import ColorPopUp from '../Ui/ColorPopUp';
import Modal from '../Utils/Modal';
import LabelPopUP from '../Ui/LabelPopUP';
import { useFetchTodoItems } from '../../Hooks/useFetchTodoItems';
import { useAuth } from '../../Context/AuthContext'; // Import Auth
import { updateTodoItem, createTodoInList, deleteTodoItem } from '../../Services/apiServices'; // Import API services
import Loading from '../Loading/Loading.jsx';
import ErrorComponent from '../Utils/ErrorComponent.jsx';
import { useLists } from '../../Context/ListContext.jsx';

export default function ListCard({ list }) {
    const { user } = useAuth();

    const { todoItems, setTodoItems, loading, error } = useFetchTodoItems(list.id);
    const { handleDeleteList, handleUpdateList } = useLists();

    const [inputValue, setInputValue] = useState("");
    const [activeMenu, setActiveMenu] = useState(null);

    const [themeColor, setThemeColor] = useState(list.color);


    const theme = COLOR_PALETTE[list.color] || COLOR_PALETTE.default;


    const handleToggleTodo = async (id) => {
        //Find the todo to toggle
        const todoToToggle = todoItems.find(todo => todo.id === id);
        if (!todoToToggle) return;

        // Optimistic UI Update
        const previousTodos = [...todoItems];
        setTodoItems(todoItems.map(todo =>
            todo.id === id ? { ...todo, completed: !todo.completed } : todo
        ));

        // API Call
        try {
            await updateTodoItem(user, list.id, id, { completed: !todoToToggle.completed });
        } catch (err) {
            console.error('Failed to toggle todo:', err);
            setTodoItems(previousTodos); // Revert on failure
        }
    };


    const handleDeleteTodo = async (id) => {
        const previousTodos = [...todoItems];

        setTodoItems(todoItems.filter(todo => todo.id !== id));

        // API Call
        try {
            deleteTodoItem(user, list.id, id);
        } catch (err) {
            console.error('Failed to delete todo:', err);
            setTodoItems(previousTodos); // Revert on failure
        }
    };

    const handleEditTodo = async (id, newTitle) => {
        const previousTodos = [...todoItems];

        setTodoItems(todoItems.map(todo =>
            todo.id === id ? { ...todo, title: newTitle } : todo
        ));

        // API Call
        try {
            updateTodoItem(user, list.id, id, { title: newTitle });
        } catch (err) {
            console.error('Failed to edit todo:', err);
            setTodoItems(previousTodos); // Revert on failure
        }
    };

    const handleAddTodo = async (e) => {
        if (e.key === 'Enter' && inputValue.trim()) {
            setTodoItems([...todoItems, { id: Date.now(), title: inputValue, completed: false }]);
            setInputValue("");
        }
        try {
            const newTodo = await createTodoInList(user, list.id, { title: inputValue.trim(), completed: false });
            setTodoItems(prevTodos => [...prevTodos, newTodo]);
            setInputValue("");
        } catch (err) {
            console.error('Failed to add todo:', err);
        }
    };

    const handleColorChange = (colorKey) => {
        const updatedList = { ...list, color: colorKey };
        handleUpdateList(updatedList);
        setActiveMenu(null);
    }

    const handleLabelsChange = (newLabels) => {
        const updatedList = { ...list, labels: newLabels };
        handleUpdateList(updatedList);
        setActiveMenu(null);
    };

    const confirmDelete = (listId) => {
        handleDeleteList(listId);
        setActiveMenu(null);
    }


    // DnD Kit Setup activationConstraint: tells dnd-kit: "Wait until the user moves 8px before assuming it's a drag."
    // This allows clicks on checkboxes/buttons to pass through instantly.
    const sensors = useSensors(
        useSensor(PointerSensor, { activationConstraint: { distance: 8 } }),
        useSensor(KeyboardSensor, {
            coordinateGetter: sortableKeyboardCoordinates,
        })
    );

    const getTaskPos = (id) => {
        return todoItems.findIndex((t) => t.id === id);
    }

    const handleDragEnd = (event) => {
        const { active, over } = event;

        if (active.id === over.id) return;

        setTodoItems((todos) => {
            const originalPos = getTaskPos(active.id);
            const newPos = getTaskPos(over.id);

            return arrayMove(todos, originalPos, newPos);
        });
    };

    // Sort: Active first, Completed last

    const safeTodos = Array.isArray(todoItems) ? todoItems : [];


    const activeTodos = safeTodos.filter(todo => !todo.completed);
    const completedTodos = todoItems.filter(todo => todo.completed);

    if (loading) return <Loading />;
    if (error) return <ErrorComponent message={error} />;

    return (
        <DndContext sensors={sensors} collisionDetection={closestCorners} onDragEnd={handleDragEnd}>

            {/* CARD CONTAINER */}
            <div className={`break-inside-avoid mb-6 rounded-2xl shadow-sm bg-white hover:shadow-xl transition-all duration-300 flex flex-col border border-gray-100 group overflow-visible`}>

                <div className="flex flex-col h-full">

                    {/* 1. COLOR BAR */}
                    <div className={`h-1.5 w-full rounded-t-2xl ${theme.bar}`}></div>

                    {/* 2. HEADER & LIST CONTENT */}
                    <div className="p-0">
                        {/* Title */}
                        <div className="p-4 pb-1">
                            <h3 className="font-bold text-lg text-gray-800">{list.title}</h3>
                            <p className="text-xs text-gray-400 font-medium uppercase tracking-wider mt-1">
                                {list.labels.join(' • ')}
                            </p>
                        </div>

                        {/* Active Tasks */}
                        <ul className='space-y-1 mb-2 w-full'>
                            <SortableContext items={activeTodos} strategy={verticalListSortingStrategy}>
                                {activeTodos.map(todo => (
                                    <TaskItem
                                        key={todo.id}
                                        todoItem={todo}
                                        checkboxColor={theme.checkbox}
                                        hoverColor={theme.hover}
                                        onToggle={handleToggleTodo}
                                        onDelete={handleDeleteTodo}
                                        onEdit={handleEditTodo}
                                    />
                                ))}
                            </SortableContext>
                        </ul>

                        {/* Completed Divider */}
                        {activeTodos.length > 0 && completedTodos.length > 0 && (
                            <div className="divider my-0 px-4 text-xs text-gray-300 font-semibold uppercase">Completed</div>
                        )}

                        {/* Completed Tasks */}
                        <ul className='space-y-1 mb-2 w-full'>
                            {completedTodos.map(item => (
                                <TaskItem
                                    key={item.id}
                                    todoItem={item}
                                    onToggle={handleToggleTodo}
                                    onDelete={handleDeleteTodo}
                                    onEdit={handleEditTodo}
                                    checkboxColor={theme.checkbox}
                                    hoverColor={theme.hover} />
                            ))}
                        </ul>
                    </div>

                    {/* 3. INPUT FIELD  */}
                    <div className="px-4 pb-2 pt-2 mt-auto">
                        <div className={`flex items-center gap-3 p-2 rounded-lg ${theme.hover} transition-colors`}>
                            <Plus size={18} className="text-gray-400 flex-shrink-0" />
                            <input
                                type="text"
                                className={`bg-transparent outline-none w-full text-sm placeholder-gray-400 text-gray-700`}
                                placeholder="Add list item..."
                                value={inputValue}
                                onChange={(e) => setInputValue(e.target.value)}
                                onKeyDown={handleAddTodo}
                            />
                        </div>
                    </div>

                    {/* 4. BOTTOM ACTION BAR  */}
                    <div className={`
                        relative px-4 overflow-visible
                        w-full
                        transition-all duration-300 ease-in-out
                        max-h-0 opacity-0 
                        /* Show when hovering card OR when a menu is open */
                        ${activeMenu ? 'max-h-12 opacity-100 py-3 border-t border-gray-50' : 'group-hover:max-h-12 group-hover:opacity-100 group-hover:py-3 group-hover:border-t group-hover:border-gray-50'}
                    `}>

                        <div className="flex items-center w-full">
                            <div className="flex gap-2">

                                {/* COLOR PALETTE */}
                                <div className="relative">
                                    <button
                                        onClick={() => setActiveMenu(activeMenu === 'palette' ? null : 'palette')}
                                        className="p-1.5 rounded hover:bg-gray-100 text-gray-400 hover:text-gray-700 transition-colors tooltip tooltip-bottom"
                                        data-tip="Change Color"
                                    >
                                        <PaletteIcon size={16} />
                                    </button>

                                    <div className="absolute bottom-full left-0 mb-2">
                                        <ColorPopUp
                                            isOpen={activeMenu === 'palette'}
                                            onClose={() => setActiveMenu(null)}
                                            selectedColor={themeColor}
                                            onSelect={(key) => {
                                                setThemeColor(key); // TODO Update local state, temp DELETE later
                                                handleColorChange(key);
                                            }}
                                            palette={COLOR_PALETTE}
                                        />
                                    </div>
                                </div>

                                {/* LABELS */}
                                <div className="relative">
                                    <button
                                        onClick={() => setActiveMenu(activeMenu === 'labels' ? null : 'labels')}
                                        className="p-1.5 rounded hover:bg-gray-100 text-gray-400 hover:text-gray-700 transition-colors tooltip tooltip-bottom"
                                        data-tip="Edit Labels"
                                    >
                                        <TagsIcon size={16} />
                                    </button>

                                    <div className="absolute bottom-full left-0 mb-2">
                                        <LabelPopUP
                                            isOpen={activeMenu === 'labels'}
                                            onClose={() => setActiveMenu(null)}
                                            selectedLabels={list.labels}
                                            onUpdate={handleLabelsChange}
                                        />
                                    </div>
                                </div>
                            </div>

                            {/* Right: Delete Button */}
                            <button
                                onClick={() => setActiveMenu(activeMenu === 'delete' ? null : 'delete')}
                                className="p-1.5 ml-auto rounded hover:bg-red-50 text-gray-400 hover:text-red-500 transition-colors tooltip tooltip-bottom"
                                data-tip="Delete List"
                            >
                                <Trash2 size={16} />
                            </button>
                        </div>
                    </div>

                    {/* 5. DELETE MODAL */}
                    <Modal
                        openModal={activeMenu === 'delete'}
                        closeModal={() => setActiveMenu(null)}
                    >
                        <div className="flex flex-col items-center text-center">
                            <h3 className="text-lg font-bold text-gray-800 mb-2">Delete this list?</h3>
                            <p className="text-gray-500 text-sm mb-6">
                                Are you sure you want to delete <strong>"{list.title}"</strong>?
                            </p>
                            <div className="flex gap-3 w-full">
                                <button onClick={() => setActiveMenu(null)} className="flex-1 px-4 py-2 bg-gray-100 hover:bg-gray-200 text-gray-700 rounded-lg font-medium">
                                    Cancel
                                </button>
                                <button onClick={() => confirmDelete(list.id)} className="flex-1 px-4 py-2 bg-red-500 hover:bg-red-600 text-white rounded-lg font-medium shadow-sm">
                                    Delete
                                </button>
                            </div>
                        </div>
                    </Modal>

                </div>
            </div>
        </DndContext>
    );
};
