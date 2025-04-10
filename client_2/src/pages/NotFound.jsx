import React from 'react';
import { Button, Container } from 'react-bootstrap';
import { Link } from 'react-router-dom';

const NotFound = () => {
    return (
        <Container className="text-center py-5">
            <h1 className="display-1">404</h1>
            <h2 className="mb-4">Page Not Found</h2>
            <p className="lead mb-4">
                The page you're looking for doesn't exist or has been moved.
            </p>
            <Button as={Link} to="/" variant="primary" size="lg">
                Go to Homepage
            </Button>
        </Container>
    );
};

export default NotFound;