# React + Vite

This template provides a minimal setup to get React working in Vite with HMR and some ESLint rules.

Currently, two official plugins are available:

- [@vitejs/plugin-react](https://github.com/vitejs/vite-plugin-react/blob/main/packages/plugin-react) uses [Babel](https://babeljs.io/) (or [oxc](https://oxc.rs) when used in [rolldown-vite](https://vite.dev/guide/rolldown)) for Fast Refresh
- [@vitejs/plugin-react-swc](https://github.com/vitejs/vite-plugin-react/blob/main/packages/plugin-react-swc) uses [SWC](https://swc.rs/) for Fast Refresh

## React Compiler

The React Compiler is not enabled on this template because of its impact on dev & build performances. To add it, see [this documentation](https://react.dev/learn/react-compiler/installation).

## Expanding the ESLint configuration

If you are developing a production application, we recommend using TypeScript with type-aware lint rules enabled. Check out the [TS template](https://github.com/vitejs/vite/tree/main/packages/create-vite/template-react-ts) for information on how to integrate TypeScript and [`typescript-eslint`](https://typescript-eslint.io) in your project.


client/
├── node_modules/
├── public/
│   └── vite.svg
├── src/
│   ├── assets/                 # Static images, logos
│   │   └── logo.png
│   │
│   ├── components/             # Reusable UI components
│   │   ├── auth/               # Authentication specific components
│   │   │   └── AuthForm.jsx    # The Login/Register toggle form
│   │   │
│   │   ├── layout/             # Structural components
│   │   │   ├── Sidebar.jsx     # The "Sticky Wall" side menu
│   │   │   ├── Header.jsx      # Search bar, view toggles, user avatar
│   │   │   └── MainLayout.jsx  # Wrapper to hold Sidebar + Content
│   │   │
│   │   ├── tasks/              # Task-related components
│   │   │   ├── TaskGrid.jsx    # Logic for switching between Grid/List view
│   │   │   ├── TaskCard.jsx    # The individual note card
│   │   │   └── TaskItem.jsx    # Individual checklist items (with the hover-delete)
│   │   │
│   │   └── ui/                 # Generic UI elements (Buttons, Inputs)
│   │       ├── Button.jsx
│   │       └── IconButton.jsx
│   │
│   ├── data/                   # Mock data (temporary until backend is ready)
│   │   └── mockData.js         # Your INITIAL_TASKS and INITIAL_LABELS
│   │
│   ├── pages/                  # Main full-screen pages
│   │   ├── LoginPage.jsx       # The Purple login screen
│   │   └── DashboardPage.jsx   # The main app view
│   │
│   ├── hooks/                  # Custom logic hooks
│   │   └── useLocalStorage.js  # (Optional) To save data in browser
│   │
│   ├── App.jsx                 # Main entry point, handles Routing
│   ├── index.css               # Tailwind imports (@tailwind base; etc.)
│   └── main.jsx                # React DOM rendering
│
├── .gitignore
├── index.html
├── package.json
├── postcss.config.js
├── tailwind.config.js          # Theme colors (Purple) configuration
└── vite.config.js