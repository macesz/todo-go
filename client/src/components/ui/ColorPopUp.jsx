import { Check } from 'lucide-react';

export default function ColorPopUp({ isOpen, onClose, selectedColor, onSelect, palette }) {

    if (!isOpen) return null;



    const handleCheckColor = (key) => {
        onSelect(key)

        setTimeout(() => {
            onClose();
        }, 800)
    };

    return (
        <>
            {/* A. The Invisible Backdrop (Closes menu when clicking outside) */}
            <div
                className="fixed inset-0 z-30"
                onClick={() => onClose(null)}
            ></div>

            {/* B. The Actual Popup */}
            <div className="absolute left-0 mb-2 z-60 bg-white p-3 rounded-xl shadow-xl border border-gray-100 w-64 cursor-default animation-fade-in">

                {/* Grid of Colors */}
                <div className="grid grid-cols-5 gap-2">
                    {Object.entries(palette).map(([key, themeData]) => {

                        const isSelected = selectedColor === key;
                        console.log("isSelected: ", isSelected);


                        return (
                            <button
                                key={key}
                                onClick={(e) => {
                                    e.preventDefault();
                                    e.stopPropagation();
                                    handleCheckColor(key);
                                }}
                                className={`
                                relative w-8 h-8 rounded-full transition-transform hover:scale-110 focus:outline-none
                                ${themeData.bg} /* Uses the background color from your palette */
                                ${isSelected ? 'ring-2 ring-offset-2 ring-gray-400' : 'border border-transparent'}
                            `}
                                title={key} // Shows color name on hover
                            >
                                {/* Show Checkmark if selected */}
                                {isSelected && (

                                    <div className="absolute inset-0 flex items-center justify-center text-black drop drop-shadow-md">
                                        <Check size={18} strokeWidth={3} />
                                    </div>
                                )}
                            </button>
                        );
                    })}
                </div>
            </div>
        </>
    )
}
