import React, { useState } from 'react'
import LoginForm from '../components/auth/LoginForm';
import SignUpForm from '../components/auth/SignUpForm';
import OverlayPanel from '../components/auth/OverlayPanel';

export default function AuthPage() {

    const [isSignUp, setIsSignUp] = useState(false);


    return (
        <div className="min-h-screen flex items-center justify-center bg-gray-50 py-12 px-4 sm:px-6 lg:px-8 relative overflow-hidden">

            {/* Background Decor (Optional blobs for aesthetic) */}
            <div className="absolute top-0 left-0 w-64 h-64 bg-violet-300 rounded-full mix-blend-multiply filter blur-xl opacity-30 animate-blob"></div>
            <div className="absolute top-0 right-0 w-64 h-64 bg-violet-300 rounded-full mix-blend-multiply filter blur-xl opacity-30 animate-blob animation-delay-2000"></div>

            {/* Main Card Container */}
            <div className="relative bg-white rounded-[20px] shadow-[0_14px_28px_rgba(0,0,0,0.1),0_10px_10px_rgba(0,0,0,0.08)] overflow-hidden w-full max-w-[850px] min-h-[550px] flex">

                {/* Sign Up Form (Left, hidden initially) */}
                <div
                    className={`absolute top-0 left-0 h-full w-1/2 flex flex-col justify-center items-center p-10 transition-all duration-700 ease-in-out bg-white z-10
          ${isSignUp ? 'translate-x-[100%] opacity-100' : 'opacity-0 z-0'}
          `}
                >
                    <SignUpForm />
                </div>

                {/* Sign In Form (Right, visible initially) */}
                <div
                    className={`absolute top-0 left-0 h-full w-1/2 flex flex-col justify-center items-center p-10 transition-all duration-700 ease-in-out bg-white z-20
          ${isSignUp ? 'translate-x-[100%] opacity-0 z-0' : 'translate-x-0 opacity-100'}
          `}
                >
                    <LoginForm />
                </div>

                {/* Overlay Container (The Purple Sliding Box) */}
                <div
                    className={`absolute top-0 left-1/2 w-1/2 h-full overflow-hidden transition-transform duration-700 ease-in-out z-[100]
          ${isSignUp ? '-translate-x-full' : 'translate-x-0'}
          `}
                >
                    {/* The Gradient Background */}
                    <div
                        className={`bg-gradient-to-br from-violet-300 to-purple-500 text-white relative -left-full h-full w-[200%] transform transition-transform duration-700 ease-in-out
            ${isSignUp ? 'translate-x-1/2' : 'translate-x-0'}
            `}
                    >
                        {/* Left Panel (Visible when showing Sign Up form) */}
                        <OverlayPanel
                            isLeft={true}
                            isSignUp={isSignUp}
                            title="Welcome Back!"
                            description="To keep connected with us please login with your personal info"
                            buttonText="Sign In"
                            onClick={() => setIsSignUp(false)}
                        />

                        {/* Right Panel (Visible when showing Sign In form) */}
                        <OverlayPanel
                            isLeft={false}
                            isSignUp={isSignUp}
                            title="Hello, Friend!"
                            description="Enter your personal details and start journey with us"
                            buttonText="Sign Up"
                            onClick={() => setIsSignUp(true)}
                        />
                    </div>
                </div>

            </div>
        </div>
    );
}
