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

    const handleClose = () => {
        closeModal();
    };

    return (
        <dialog ref={ref} onCancel={handleClose}
            className="absolute inset-1/2 transform -translate-x-1/2 -translate-y-1/2 bg-transparent p-6 rounded-lg w-120 overflow-y-hidden">
            <div className="bg-violet-50 p-6 rounded-lg shadow-lg w-96">
                {children(handleClose)}
            </div>
        </dialog>
    );
}

export default Modal;