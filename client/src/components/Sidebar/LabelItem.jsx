
import { COLOR_PALETTE } from '../../data/ColorPalette';

export default function LabelItem({ label, onClick }) {

    if (!label) return null;

    const theme = COLOR_PALETTE[label.color] || COLOR_PALETTE.default;
    
    return (
        <li
            onClick={onClick}
            className="flex items-center gap-3 p-2 rounded-lg cursor-pointer text-gray-600 hover:bg-gray-100 transition-colors"
        >
            <div className={`w-3 h-3 rounded-full ${theme.bar}`} />
            <span className="text-sm font-medium">{label.name}</span>
        </li>)
}
