import React from 'react';
import { Card, Container, Row, Col } from 'react-bootstrap';
import { Link, useNavigate } from 'react-router-dom';
import AuthForm from '../components/AuthForm';
import { useAuth } from '../context/AuthContext';

const Login = () => {
    const { login } = useAuth();
    const navigate = useNavigate();

    const handleLogin = async (username, password) => {
        const result = await login(username, password);

        if (result.success) {
            navigate('/dashboard');
        } else {
            throw new Error(result.error);
        }
    };

    return (
        <Container>
            <Row className="justify-content-center">
                <Col md={6}>
                    <div className="auth-container">
                        <h2 className="text-center mb-4">Login</h2>
                        <AuthForm type="login" onSubmit={handleLogin} />
                        <div className="text-center mt-3">
                            <p>
                                Don't have an account? <Link to="/register">Register</Link>
                            </p>
                        </div>
                    </div>
                </Col>
            </Row>
        </Container>
    );
};

export default Login;
