
import axios from "axios";

const api = axios.create({
    baseURL: '/api',
    timeout: 10000,
    headers: {
        "Content-Type": "application/json",
    },
});

api.interceptors.request.use(
    (config) => {
        if (config.user && config.user.jwt) {
            config.headers['Authorization'] = `Bearer ${config.user.jwt}`;
            delete config.user;
        }
        return config;
    },
    (error) => {
        return Promise.reject(error);
    }
);

// List all todo lists for a user
export const fetchTodoLists = async (user) => {
    try {
        const response = await api.get("/lists", {
            user
        });

        return response.data;
    } catch (error) {
        console.error('Error getting todo lists:', error);
    }
};

// Create a new todo list
export const createTodoList = async (user, listData) => {
    try {
        const response = await api.post("/lists", listData, {
            user
        });

        return response.data;
    } catch (error) {
        console.error('Error creating todo list:', error);
    }
};

// Get List by ID
export const fetchTodoListById = async (user, listId) => {
    try {
        const response = await api.get(`/lists/${listId}`, {
            user
        });

        return response.data;
    } catch (error) {
        console.error('Error getting todo list by ID:', error);
    }
};

// Update a todo list
export const updateTodoList = async (user, listId, updatedData) => {
    try {
        const response = await api.put(`/lists/${listId}`, updatedData, {
            user
        });

        return response.data;
    } catch (error) {
        console.error('Error updating todo list:', error);
    }
};

// Delete a todo list
export const deleteTodoList = async (user, listId) => {
    try {
        const response = await api.delete(`/lists/${listId}`, {
            user
        });

        return response.data;
    } catch (error) {
        console.error('Error deleting todo list:', error);
    }
};

// Get all todos for a specific list
export const fetchTodosByListId = async (user, listId) => {
    try {
        const response = await api.get(`/lists/${listId}/todos`, {
            user
        });

        return response.data;
    } catch (error) {
        console.error('Error getting todos by list ID:', error);
    }
};

// Get todo list by label
export const fetchTodoListByLabel = async (user, label) => {
    try {
        const response = await api.get(`/lists/${label}`, {
            user
        });

        return response.data;
    } catch (error) {
        console.error('Error getting todo list by label:', error);
    }
};

// Get todos on a specific list
export const fetchTodosInList = async (user, listId) => {
    try {
        const response = await api.get(`/lists/${listId}/todos`, {
            user
        });

        return response.data;
    } catch (error) {
        console.error('Error getting todos in list:', error);
    }
}

// Create a new todo in a specific list
export const createTodoInList = async (user, listId, todoData) => {
    try {
        const response = await api.post(`/lists/${listId}/todos`, todoData, {
            user
        });

        return response.data;
    } catch (error) {
        console.error('Error creating todo in list:', error);
    }
};

// Update a todo item
export const updateTodoItem = async (user, listId, todoId, updatedData) => {
    try {
        const response = await api.put(`/lists/${listId}/todos/${todoId}`, updatedData, {
            user
        });

        return response.data;
    } catch (error) {
        console.error('Error updating todo item:', error);
    }
};

// Delete a todo item
export const deleteTodoItem = async (user, listId, todoId) => {
    try {
        const response = await api.delete(`/lists/${listId}/todos/${todoId}`, {
            user
        });

        return response.data;
    } catch (error) {
        console.error('Error deleting todo item:', error);
    }
};

// User authentication (login)
export const loginUser = async (credentials) => {
    try {
        const response = await api.post("/auth/login", credentials);
        return response.data;
    } catch (error) {
        console.error('Error logging in user:', error);
    }
};

// User registration
export const registerUser = async (userData) => {
    try {
        const response = await api.post("/auth/register", userData);
        return response.data;
    } catch (error) {
        console.error('Error registering user:', error);
    }
};  


export default api;
