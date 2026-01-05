import { INITIAL_TASKS_LISTS } from '../data/MockData.js';
import ListCard from '../components/Todos/ListCard.jsx';
import Loading from '../components/Loading/Loading.jsx';
import ErrorComponent from '../components/Utils/ErrorComponent.jsx';
import CreateList from '../components/Todos/CreateList.jsx';
import { useLists } from '../Context/ListContext.jsx';



export default function HomePage() {

const { lists, handleCreateList, error} = useLists(); 

    // if (loading) return <Loading />;
    if (error) return <ErrorComponent message={error} />;


    return (
        <div className='container mx-auto p-4'>
            <CreateList
                onSave={handleCreateList}
            />
            <div className="columns-1 md:columns-2 lg:columns-3 gap-6 space-y-6 mx-auto max-w-6xl">
                {lists.map(list => (
                    <ListCard
                        key={list.id}
                        list={list}
                    />
                ))}
            </div>
        </div>
    );
}
