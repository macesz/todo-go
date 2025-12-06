import { useState } from "react";
import { useNavigate } from 'react-router-dom';
import { INITIAL_LABELS } from '../../data/MockData.js';
import { ChevronsRight, Search, Settings, LogOut } from "lucide-react";
import { List, Trash2, Edit3, ChevronDown, ChevronUp } from "lucide-react";
import MenuItem from "./MenuItem.jsx";
import LabelItem from "./LabelItem.jsx";
import { useAuth } from "../../Context/AuthContext.jsx";



const Sidebar = ({ isOpen, onClose }) => {

  const { logout } = useAuth();
  const navigate = useNavigate();

  const [showAllLabels, setShowAllLabels] = useState(false);

  const safeLabels = INITIAL_LABELS || [];

  const visibleLabels = showAllLabels ? safeLabels : safeLabels.slice(0, 5);


  return (

    <>

      {/* OVERLAY (Mobile Only) */}
      {isOpen && (
        <div
          className="fixed inset-0 bg-black/20 backdrop-blur-sm z-20 md:hidden"
          onClick={onClose}
        />
      )}

      {/* SIDEBAR CONTAINER */}
      <aside
        className={`
          fixed md:static inset-y-0 left-0 z-30
          w-64 bg-gray-50 border-r border-gray-200 
          transform transition-transform duration-300 ease-in-out
          flex flex-col h-full
          ${isOpen ? 'translate-x-0' : '-translate-x-full md:translate-x-0'}
        `}
      >

        {/* SIDEBAR Header */}

        <div className="p-6 flex items-center justify-between">
          <h1 className="font-bold text-2xl text-gray-800">Menu</h1>
          <button onClick={onClose} className="md:hidden text-gray-500 hover:text-gray-800">
            <ChevronsRight size={24} className="rotate-180" />
          </button>
        </div>

        {/* SIDEBAR Search */}
        <div className="px-6 mb-6">
          <div className="relative">
            <Search className="absolute left-3 top-2.5 text-gray-400" size={16} />
            <input
              type="text"
              placeholder="Search"
              className="w-full bg-white border border-gray-200 rounded-lg pl-10 pr-4 py-2 text-sm focus:outline-none focus:border-purple-400 focus:ring-2 focus:ring-purple-100 transition-all"
            />
          </div>
        </div>

        {/* SIDEBAR Navigation */}
        <div className="flex-1 overflow-y-auto px-6 custom-scrollbar">

          <ul className="space-y-2 mb-8">
            {/* TODOS (Main View) */}
            <MenuItem
              icon={<List />}
              label="Todos"
              onClick={() => navigate('/')} // Goes to HomePage
              active={true} // You can make this dynamic based on current route
            />

            {/* BIN */}
            <MenuItem
              icon={<Trash2 />}
              label="Bin"
              onClick={() => navigate('/bin')}
            />
          </ul>
          {/* LABELS SECTION */}
          <div className="mb-6">
            <div className="flex items-center justify-between mb-3">
              <h3 className="text-xs font-bold text-gray-400 uppercase tracking-wider">Labels</h3>
              {/* Edit Labels (Non-functional as requested) */}
              <button className="text-gray-400 hover:text-purple-600 transition-colors" title="Edit Labels">
                <Edit3 size={14} />
              </button>
            </div>

            <ul className="space-y-2">
              {/* Render Visible Labels */}
              {visibleLabels && visibleLabels.map((label) => (
                <LabelItem
                  key={label.id}
                  label={label}
                  onClick={() => navigate(`/label/${label.id}`)}
                />
              ))}

              {/* "More" Button (Only shows if there are more than 5 labels) */}
              {INITIAL_LABELS.length > 5 && (
                <li
                  onClick={() => setShowAllLabels(!showAllLabels)}
                  className="flex items-center gap-3 p-2 rounded-lg cursor-pointer text-gray-500 hover:text-gray-800 hover:bg-gray-100 transition-colors"
                >
                  {showAllLabels ? <ChevronUp size={18} /> : <ChevronDown size={18} />}
                  <span className="text-sm font-medium">
                    {showAllLabels ? 'Show Less' : 'More...'}
                  </span>
                </li>
              )}
            </ul>
          </div>

        </div>

        {/* FOOTER */}
        <div className="p-6 border-t border-gray-200 space-y-2">
          <div className="flex items-center gap-3 p-2 text-gray-600 hover:bg-gray-100 rounded-lg cursor-pointer transition-colors">
            <Settings size={20} />
            <span className="text-sm font-medium">Settings</span>
          </div>
          <button
            onClick={logout}
            className="w-full flex items-center gap-3 p-2 text-red-500 hover:bg-red-50 rounded-lg cursor-pointer transition-colors"
          >
            <LogOut size={20} />
            <span className="text-sm font-medium">Sign out</span>
          </button>
        </div>
      </aside>

    </>

  );

};

export default Sidebar;

