import { useState, useCallback, useEffect } from 'react';
import { updateTodoItem, createTodoInList, deleteTodoItem } from '../Services/apiServices';

export const useTodoItems = (initialItems, listId, user, onSyncGlobal) => {
    const [todoItems, setTodoItems] = useState(initialItems || []);

    // Sync to global list context whenever local items change
    useEffect(() => {
        onSyncGlobal(listId, todoItems)
    }, [todoItems, listId, onSyncGlobal])

    // 1. ADD TODO
   const addTodo = useCallback(async (title) => {
        try {
            const newTodo = await createTodoInList(user, listId, { title, done: false });
            if (newTodo) {
                setTodoItems(prev => [...prev, newTodo]);
            }
        } catch (err) {
            console.error('Failed to add todo:', err);
        }
    }, [user, listId]);

    // 2. TOGGLE DONE
    const toggleTodo = useCallback(async (id) => {
        setTodoItems(prev => {
            const todo = prev.find(t => t.id === id);
            if (!todo) return prev;

            // Perform API call in background
            updateTodoItem(user, listId, id, { done: !todo.done })
                .catch(() => setTodoItems(initial => [...initial])); // Simple rollback logic

            return prev.map(t => t.id === id ? { ...t, done: !t.done } : t);
        });
    }, [user, listId]);

    // 3. DELETE TODO
    const deleteTodo = useCallback(async (id) => {
        setTodoItems(prev => prev.filter(t => t.id !== id));
        try {
            await deleteTodoItem(user, listId, id);
        } catch (err) {
            console.error('Delete failed:', err);
        }
    }, [user, listId]);

    // 4. EDIT TODO
   const editTodo = useCallback(async (id, newTitle) => {
        setTodoItems(prev => prev.map(t => t.id === id ? { ...t, title: newTitle } : t));
        try {
            await updateTodoItem(user, listId, id, { title: newTitle });
        } catch (err) {
            console.error('Edit failed:', err);
        }
    }, [user, listId]);

    return {
        todoItems,
        setTodoItems,
        addTodo,
        toggleTodo,
        deleteTodo,
        editTodo
    };
};