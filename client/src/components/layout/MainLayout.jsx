import { Outlet } from "react-router-dom";
import Sidebar from "./Sidebar.jsx";
import { useState } from "react";
import { Menu } from "lucide-react";


const MainLayout = () => {

    const [isSidebarOpen, setIsSidebarOpen] = useState(false);

    return (
        <div className="flex h-screen bg-white overflow-hidden font-sans">

            {/* SIDEBAR */}
            {/* Pass state and close function to the sidebar */}
            <Sidebar
                isOpen={isSidebarOpen}
                onClose={() => setIsSidebarOpen(false)}
            />

            {/* MAIN CONTENT WRAPPER */}
            <div className="flex-1 flex flex-col min-w-0 overflow-hidden">

                {/* MOBILE HEADER (Visible only on small screens) */}
                {/* This allows opening the menu on mobile */}
                <header className="md:hidden bg-white border-b border-gray-200 p-4 flex items-center gap-4 sticky top-0 z-10">
                    <button
                        onClick={() => setIsSidebarOpen(true)}
                        className="p-2 hover:bg-gray-100 rounded-lg text-gray-600 transition-colors"
                        aria-label="Open Menu"
                    >
                        <Menu size={24} />
                    </button>
                    <h1 className="font-bold text-xl text-gray-800">Menu</h1>
                </header>

                {/* PAGE CONTENT AREA */}
                {/* 'Outlet' renders the child routes (like HomePage, BinPage, etc.) */}
                <main className="flex-1 overflow-y-auto p-4 md:p-8 bg-white scroll-smooth">
                    <Outlet />
                </main>

            </div>
        </div>
    );
};

export default MainLayout;  
