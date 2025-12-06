
export default function LabelItem({ label, onClick }) {

    if (!label) return null;


    console.log("Rendering LabelItem: color:", label.color);


    const labelColorClass = label.color || 'bg-purple-400';
    return (
        <li
            onClick={onClick}
            className="flex items-center gap-3 p-2 rounded-lg cursor-pointer text-gray-600 hover:bg-gray-100 transition-colors"
        >
            <div className={`w-3 h-3 rounded-full ${labelColorClass}`} />
            <span className="text-sm font-medium">{label.name}</span>
        </li>)
}
