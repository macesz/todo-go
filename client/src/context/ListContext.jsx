import { createContext, useContext, useState, useMemo, useCallback } from 'react';
import { useAuth } from './AuthContext';
import { createTodoList, createTodoInList, updateTodoList, deleteTodoList } from '../Services/apiServices';
import { useFetchLists } from '../Hooks/useFetchLists';

const ListContext = createContext();

export const ListProvider = ({ children }) => {
    const { user } = useAuth();
    const { lists, setLists, error } = useFetchLists();
    const [loading, setLoading] = useState(true);
    const [selectedLabel, setSelectedLabel] = useState(null);

    console.log("Lists", lists);

    const updateListItemsLocally = useCallback((listId, newItems) => {
        setLists(prevLists => prevLists.map(list =>
            list.id === listId ? { ...list, items: newItems } : list
        ));
    }, [setLists]);


    // Generate unique labels from lists
    const uniqueLabels = useMemo(() => {
        const labelsSet = new Set();

        let color = 'bg-accent'; // default color

        lists.forEach(list => {
            if (list.labels) {
                list.labels.forEach(label => labelsSet.add(label));
            }
        });

        return Array.from(labelsSet).map(name => ({
            id: name, // Use name as ID
            name: name,
            color: color
        })).sort((a, b) => a.name.localeCompare(b.name));
    }, [lists]);


    // Filter lists based on selected label
    const filteredList = useMemo(() => {
        let filtered = lists;
        if (selectedLabel) {
            filtered = lists.filter(list => list.labels && list.labels.includes(selectedLabel));
        }

        return filtered;
        // return selectedLabel
        //     ? lists.filter(list => list.labels && list.labels.includes(selectedLabel))
        //     : lists;
    }, [lists, selectedLabel]);

    // Filter click handler
    const filterByLabel = (labelName) => {
        setSelectedLabel(labelName);
    }

    // Clear filter
    const clearFilter = () => {
        setSelectedLabel(null);
    }




    // Create, Update, Delete Handlers
    const handleCreateList = async (listData) => {
        const { items, ...listDetails } = listData;

        try {
            setLoading(true);

            const createdList = await createTodoList(user, listDetails);

            let createdItems = [];
            if (items && items.length > 0) {
                createdItems = await Promise.all(items.map(item =>
                    createTodoInList(user, createdList.id, {
                        title: item.title,
                        done: item.done
                    })
                ));
            }
            const fullListData = { ...createdList, todos: createdItems };


            setLists(prevLists => [fullListData, ...prevLists]);

        } catch (err) {
            console.error('Failed to create list:', err);
        } finally {
            setLoading(false);
        }
    }


    const handleDeleteList = useCallback(async (listId) => {
        if (!window.confirm("Are you sure you want to delete this list?")) return;
        try {
            await deleteTodoList(user, listId);
            setLists(prevLists => prevLists.filter(list => list.id !== listId));
        } catch (err) {
            console.error('Failed to delete list:', err);
        }
    }, [user, setLists]);

    const handleUpdateList = useCallback(async (updatedList) => {
        try {
            const updated = await updateTodoList(user, updatedList.id, updatedList);
            setLists(prevLists => prevLists.map(list => list.id === updated.id ? updated : list));
        } catch (err) {
            console.error('Failed to update list:', err);
        }
    }, [user, setLists]);

    // Rename Labels
    const renameLabelGlobally = async (oldName, newName) => {
        if (oldName === newName) return;

        const updatedLists = lists.map(list => ({
            ...list,
            labels: list.labels ? list.labels.map(l => l === oldName ? newName : l) : []
        }));
        setLists(updatedLists);

        handleUpdateList(updatedLists)
    };

    // Delete Labels
    const deleteLabelGlobally = async (labelName) => {
        const updatedLists = lists.map(list => ({
            ...list,
            labels: list.labels ? list.labels.filter(l => l !== labelName) : []
        }));
        setLists(updatedLists);

        handleUpdateList(updatedLists)
    };

    return (
        <ListContext.Provider
            value={{
                lists: filteredList,
                updateListItemsLocally,
                loading,
                error,
                uniqueLabels,
                selectedLabel,
                filterByLabel,
                clearFilter,
                handleCreateList,
                handleDeleteList,
                handleUpdateList,
                renameLabelGlobally,
                deleteLabelGlobally,
            }}
        >
            {children}
        </ListContext.Provider>
    );

}

// eslint-disable-next-line react-refresh/only-export-components
export const useLists = () => useContext(ListContext);