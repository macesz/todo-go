import React from 'react'
import { useState } from 'react';
import { INITIAL_TASKS_LISTS } from '../data/MockData.js';
import ListCard from '../components/Todos/ListCard.jsx';



export default function HomePage() {

    const [lists, setLists] = useState(INITIAL_TASKS_LISTS);

    //TODO: Fetch lists from backend API and setLists
    //TODO implement delete, update list functions and pass as props to ListCard



    return (
        <div className="min-h-screen bg-base-100 p-8 rounded-lg shadow-[0_4px_20px_-4px_rgba(0,0,0,0.1)]">
            {/* Masonry Grid Layout */}
            <div className="columns-1 md:columns-2 lg:columns-3 gap-6 space-y-6 mx-auto max-w-6xl">
                {lists.map(list => (

                    <ListCard
                        key={list.id}
                        list={list}
                        onSetLists={setLists}
                    />
                ))}
            </div>
        </div>
    )
}
