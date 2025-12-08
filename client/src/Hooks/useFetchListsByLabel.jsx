import { useEffect, useState } from 'react'
import { useAuth } from '../Context/AuthContext'; 
import { fetchTodoListByLabel } from '../Services/apiServices';

export const useFetchListsByLabel = (label) => {
    const [lists, setLists] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);

    const { user } = useAuth();

    useEffect(() => {
        const fetchListsByLabel = async () => {
            try {
                setLoading(true);
                const result = fetchTodoListByLabel(user, label);
                setLists(result);
            } catch (err) {
                setError(err.message);
            } finally {
                setLoading(false);
            }
        };

        if (user) fetchListsByLabel();
    }, [user, label]);

    return { lists, loading, error };   
 
}
