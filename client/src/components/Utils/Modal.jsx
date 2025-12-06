import { useEffect, useRef } from "react";

const Modal = ({ openModal, closeModal, children, }) => {
    const ref = useRef();

    useEffect(() => {
        if (openModal && ref.current) {
            ref.current?.showModal();
        } else if (ref.current) {
            ref.current?.close();
        }
    }, [openModal]);


    return (
        <dialog
            ref={ref}
            onCancel={closeModal}
            onClick={(e) => {
                // Close if clicking the backdrop (the dialog element itself)
                if (e.target === ref.current) closeModal();
            }}
            className="m-auto backdrop:bg-black/30 bg-transparent p-0 rounded-lg shadow-2xl outline-none"
        >
            <div className="bg-white p-6 rounded-xl w-96 border border-gray-100">
                {children}
            </div>
        </dialog>
    );
}

export default Modal;