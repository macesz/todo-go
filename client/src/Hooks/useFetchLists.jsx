import React, { useEffect, useState } from 'react'
import { useAuth } from '../Context/AuthContext';
import { fetchTodoLists } from '../Services/apiServices';

export const useFetchLists = () => {



    const [lists, setLists] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);

    const { user } = useAuth();

    useEffect(() => {
        const fetchLists = async () => {

            if (!user) {
                setLoading(false);
                return;
            }

            try {
                setLoading(true);
                const result = await fetchTodoLists(user);

                if (result && Array.isArray(result) && result.length > 0) {
                    setLists(result)
                } else {
                    console.log("Backend returned empty.");
                    setLists([]);
                }

                // setLists(result);
            } catch (err) {
                setError(err.message);
                console.error("Fetch error:", err);
                setError(err.message);
                // Fallback on error too
                setLists([]);
            } finally {
                setLoading(false);
            }
        };

        if (user) fetchLists();
    }, [user]);

    return { lists, setLists, loading, error };

}
