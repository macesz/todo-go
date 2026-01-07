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
    const [searchQuery, setSearchQuery] = useState("");
    const [view, setView] = useState('home');


    const updateListItemsLocally = useCallback((listId, newItems) => {
        setLists(prevLists => prevLists.map(list =>
            list.id === listId ? { ...list, items: newItems } : list
        ));
    }, [setLists]);


    // Generate unique labels from lists
    const uniqueLabels = useMemo(() => {
        const labelMap = new Map();

        lists.forEach(list => {
            if (list.labels && Array.isArray(list.labels)) {
                list.labels.forEach(labelName => {
                    // 2. Only add if we haven't seen this label yet 
                    // (or you can overwrite if you want the "latest" color)
                    if (!labelMap.has(labelName)) {
                        labelMap.set(labelName, list.color || 'default');
                    }
                });
            }
        });

        // 3. Convert the Map into your desired array format
        return Array.from(labelMap.entries()).map(([name, color]) => ({
            id: name,
            name: name,
            color: color
        })).sort((a, b) => a.name.localeCompare(b.name));
    }, [lists]);


    // Filter lists

    const applySearchFilter = (allLists, query) => {
        const cleanQuery = query.toLowerCase().trim();
        if (!cleanQuery) return;

        return allLists.filter(list => {
            const inTitle = list.title.toLowerCase().includes(cleanQuery);
            const inItems = list.items?.some(item =>
                item.title.toLowerCase().includes(cleanQuery)
            );
            return inTitle || inItems
        });
    };

    const filteredList = useMemo(() => {

        let visibleLists = lists;

        if (view === 'bin') {
            visibleLists = visibleLists.filter(list => list.deleted === true)
        } else {
            visibleLists = visibleLists.filter(list => list.deleted === false);
        }

        if (selectedLabel) {
            visibleLists = visibleLists.filter(list => list.labels?.includes(selectedLabel))
        }

        if (searchQuery.trim() != "") {
            return applySearchFilter(visibleLists, searchQuery)
        }

        return visibleLists

    }, [lists, selectedLabel, searchQuery, view]);

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
            const fullListData = { ...createdList, items: createdItems };


            setLists(prevLists => [fullListData, ...prevLists]);

        } catch (err) {
            console.error('Failed to create list:', err);
        } finally {
            setLoading(false);
        }
    }

    const moveToBin = useCallback(async (listId) => {
        const listToUpdate = lists.find(list => list.id === listId);
        const updatedList = { ...listToUpdate, deleted: true };
        await updateTodoList(user, listId, updatedList);
        setLists(prevLists => prevLists.map(list =>
            list.id === listId ? updatedList : list
        ));
    }, [user, lists, setLists]);

    const restoreFromBin = useCallback(async (listId) => {
        const listToMove = lists.find(l => l.id === listId);
        const updatedList = { ...listToMove, deleted: false };

        await updateTodoList(user, listId, updatedList);
        setLists(prevLists => prevLists.map(list =>
            list.id === listId ? updatedList : list
        ));
    }, [user, lists, setLists]);


    const handleDeleteList = useCallback(async (listId) => {
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

        const affectedLists = lists.filter(l => l.labels?.includes(oldName));
        for (const list of affectedLists) {
            const updatedLabels = list.labels.map(l => l === oldName ? newName : l);
            await handleUpdateList({ ...list, labels: updatedLabels });
        }
    };

    // Delete Labels
    const deleteLabelGlobally = async (labelName) => {
        const updatedLists = lists.map(list => ({
            ...list,
            labels: list.labels ? list.labels.filter(l => l !== labelName) : []
        }));
        setLists(updatedLists);

        const affectedLists = lists.filter(l => l.labels?.includes(labelName));
        for (const list of affectedLists) {
            const updatedLabels = list.labels.filter(l => l !== labelName);
            await handleUpdateList({ ...list, labels: updatedLabels });
        }
    };

    return (
        <ListContext.Provider
            value={{
                lists: filteredList,
                updateListItemsLocally,
                loading,
                error,
                uniqueLabels,
                searchQuery,
                setSearchQuery,
                view,
                setView,
                selectedLabel,
                filterByLabel: (name) => {
                    setView('home');
                    setSelectedLabel(name);
                },
                clearFilter: () => {
                    setSelectedLabel(null);
                    setSearchQuery("");
                },
                handleCreateList,
                moveToBin,
                restoreFromBin,
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