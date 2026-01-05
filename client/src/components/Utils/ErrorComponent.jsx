const ErrorComponent = (message) => {
    return (
        <div>
            <h2>An unexpected error occurred</h2>
            <p>{message}</p>
        </div>
    )
};

export default ErrorComponent;