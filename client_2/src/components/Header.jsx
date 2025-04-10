import React from 'react';
import { Navbar, Nav, Container, Button } from 'react-bootstrap';
import { Link, useNavigate } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';

const Header = () => {
    const { user, logout } = useAuth();
    const navigate = useNavigate();

    const handleLogout = () => {
        logout();
        navigate('/');
    };

    return (
        <Navbar bg="dark" variant="dark" expand="lg" sticky="top">
            <Container>
                <Navbar.Brand as={Link} to="/">
                    <i className="bi bi-bus-front me-2"></i>
                    CityBus Tracker
                </Navbar.Brand>
                <Navbar.Toggle aria-controls="main-navbar" />
                <Navbar.Collapse id="main-navbar">
                    <Nav className="me-auto">
                        <Nav.Link as={Link} to="/">Home</Nav.Link>
                        {user && (
                            <Nav.Link as={Link} to="/dashboard">Dashboard</Nav.Link>
                        )}
                    </Nav>
                    <Nav>
                        {user ? (
                            <>
                                <Navbar.Text className="me-3">
                                    Signed in as: <span className="text-primary">{user.username}</span> ({user.role})
                                </Navbar.Text>
                                <Button variant="outline-light" onClick={handleLogout}>Logout</Button>
                            </>
                        ) : (
                            <>
                                <Nav.Link as={Link} to="/login">Login</Nav.Link>
                                <Nav.Link as={Link} to="/register">Register</Nav.Link>
                            </>
                        )}
                    </Nav>
                </Navbar.Collapse>
            </Container>
        </Navbar>
    );
};

export default Header;