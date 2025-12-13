import React, { useEffect, useState } from 'react'
import { fetchTodosInList } from '../Services/apiServices';
import { useAuth } from '../Context/AuthContext';


export const useFetchTodoItems = (listId) => {

    const [todoItems, setTodoItems] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);
    const { user } = useAuth();


    useEffect(() => {
        const fetchTodoItems = async () => {
            try {
                setLoading(true);
                const result = await fetchTodosInList(user, listId);

                setTodoItems(result);
            } catch (err) {
                setError(err.message);
            } finally {
                setLoading(false);
            }
        };

        if (listId) fetchTodoItems();
    }, [user, listId]);

    return { todoItems, setTodoItems, loading, error };
}