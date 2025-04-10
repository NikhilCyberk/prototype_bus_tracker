import React from 'react';
import { Card, Container, Row, Col } from 'react-bootstrap';
import { Link, useNavigate } from 'react-router-dom';
import AuthForm from '../components/AuthForm';
import { useAuth } from '../context/AuthContext';

const Register = () => {
    const { register, login } = useAuth();
    const navigate = useNavigate();

    const handleRegister = async (username, password, role) => {
        const result = await register(username, password, role);

        if (result.success) {
            // Auto login after successful registration
            await login(username, password);
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
                        <h2 className="text-center mb-4">Create an Account</h2>
                        <AuthForm type="register" onSubmit={handleRegister} />
                        <div className="text-center mt-3">
                            <p>
                                Already have an account? <Link to="/login">Login</Link>
                            </p>
                        </div>
                    </div>
                </Col>
            </Row>
        </Container>
    );
};

export default Register;