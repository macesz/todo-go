import React, { useState } from 'react'
import { INITIAL_TASKS_LISTS } from '../data/MockData.js';
import ListCard from '../components/Todos/ListCard.jsx';
import { useFetchLists } from '../Hooks/useFetchLists.jsx';
import Loading from '../components/Loading/Loading.jsx';
import ErrorComponent from '../components/Utils/ErrorComponent.jsx';
import { useAuth } from '../Context/AuthContext';
import { createTodoList, deleteTodoList, updateTodoList, createTodoItem } from '../Services/apiServices.js';
import CreateList from '../components/Todos/CreateList.jsx';



export default function HomePage() {

    const { user } = useAuth();

    const { lists, setLists, error } = useFetchLists();

    const [loading, setLoading] = useState(false);



    const handleCreateList = async (listData) => {
        const { items, ...listDetails } = listData;

        try {
            setLoading(true);

            const createdList = await createTodoList(user, listDetails);

            let createdItems = [];
            if (items && items.length > 0) {
                createdItems = await Promise.all(items.map(item =>
                    createTodoItem(user, createdList.id, {
                        title: item.title,
                        completed: item.completed
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


    const handleDeleteList = async (listId) => {
        if (!window.confirm("Are you sure you want to delete this list?")) return;

        try {
            await deleteTodoList(user, listId);
            setLists(prevLists => prevLists.filter(list => list.id !== listId));
        } catch (err) {
            console.error('Failed to delete list:', err);
        }
    };

    const handleUpdateList = async (updatedList) => {
        try {
            const updated = await updateTodoList(user, updatedList.id, updatedList);
            setLists(prevLists => prevLists.map(list => list.id === updated.id ? updated : list));
        } catch (err) {
            console.error('Failed to update list:', err);
        }
    }


    if (loading) return <Loading />;
    if (error) return <ErrorComponent message={error} />;


    return (
        // <div className="min-h-screen bg-base-100 p-8 rounded-lg shadow-[0_4px_20px_-4px_rgba(0,0,0,0.1)]">
        <div className='container mx-auto p-4'>
            <CreateList
                onSave={handleCreateList}
            />
            {/* Masonry Grid Layout */}

            <div className="columns-1 md:columns-2 lg:columns-3 gap-6 space-y-6 mx-auto max-w-6xl">
                {lists.map(list => (
                    <ListCard
                        key={list.id}
                        list={list}
                        onDelete={() => handleDeleteList(list.id)}
                        onUpdate={handleUpdateList}
                    />
                ))}
            </div>
        </div>
    );
}
