import React, { useEffect, useState } from 'react'
import { useAuth } from '../Context/AuthContext';
import { fetchTodoLists } from '../Services/apiServices';

export const useFetchLists = () => {

    const [lists, setLists] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);

    const {user} = useAuth();

    useEffect(() => {
        const fetchLists = async () => {
            try {
                setLoading(true);
                const result = fetchTodoLists(user);
                setLists(result);
            } catch (err) {
                setError(err.message);
            } finally {
                setLoading(false);
            }
        };

        if(user) fetchLists();
    }, [user]);

    return { lists, loading, error }; 

}
