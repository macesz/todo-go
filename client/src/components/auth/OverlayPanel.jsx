import React from 'react'

export default function OverlayPanel({ isLeft, isSignUp, title, description, buttonText, onClick }) {
    return (
        <div className={`absolute top-0 ${isLeft ? 'left-0' : 'right-0'} flex flex-col justify-center items-center h-full w-1/2 px-10 text-center transition-transform duration-700 ease-in-out
      ${isLeft ? (isSignUp ? 'translate-x-0' : '-translate-x-[20%]') : (isSignUp ? 'translate-x-[20%]' : 'translate-x-0')}
    `}>
            <h1 className="font-bold text-3xl mb-4">{title}</h1>
            <p className="text-sm font-light mb-8 leading-relaxed text-purple-100">
                {description}
            </p>
            <button
                onClick={onClick}
                className="bg-transparent border-2 border-white text-white font-bold py-3 px-12 rounded-full tracking-wider uppercase text-xs transition-colors hover:bg-white hover:text-purple-400"
            >
                {buttonText}
            </button>
        </div>
    )
}
