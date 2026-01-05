export const INITIAL_LABELS = [
    { id: 1, name: 'Personal', color: 'bg-red-400' },
    { id: 2, name: 'Work', color: 'bg-blue-400' },
    { id: 3, name: 'Groceries', color: 'bg-green-400' },
    { id: 4, name: 'Finance', color: 'bg-yellow-400' },
    { id: 5, name: 'Health', color: 'bg-pink-400' },
    { id: 6, name: 'Travel', color: 'bg-purple-400' },
    { id: 7, name: 'Books', color: 'bg-indigo-400' },
    { id: 8, name: 'Movies', color: 'bg-orange-400' },
    { id: 9, name: 'Someday', color: 'bg-gray-400' },
];

export const INITIAL_TASKS_LISTS = [
    {
        id: 1,
        title: 'Groceries',
        color: '',
        labels: ['Groceries'],
        items: [
            { id: 101, userId: 1, listId: 1, title: 'Buy Apples', priority: 1, completed: false },
            { id: 102, userId: 1, listId: 1, title: 'Buy Jelly', priority: 2, completed: true },
            { id: 103, userId: 1, listId: 1, title: 'Buy Oranges', priority: 3, completed: false },
            { id: 104, userId: 1, listId: 1, title: 'Buy Bread', priority: 4, completed: false },
        ]
    },
    {
        id: 2,
        title: 'Website Project',
        color: 'blue',
        labels: ['Work'],
        items: [
            { id: 201, userId: 1, listId: 2, title: 'Design Login Page', priority: 1, completed: false },
            { id: 202, userId: 1, listId: 2, title: 'Setup React Router', priority: 2, completed: false },
            { id: 203, userId: 1, listId: 2, title: 'Integrate Tailwind', priority: 3, completed: false },
        ]
    },
    {
        id: 3,
        title: 'Weekend Plans',
        color: 'green',
        labels: ['Personal'],
        items: [
            { id: 301, userId: 1, listId: 3, title: 'Hiking trip', priority: 1, completed: false },
            { id: 302, userId: 1, listId: 3, title: 'Call Mom', priority: 2, completed: false },
        ]
    }
];