import React, { useState } from 'react';
import { Form, Button, Alert } from 'react-bootstrap';

const AuthForm = ({ type, onSubmit }) => {
    const [formData, setFormData] = useState({
        username: '',
        password: '',
        role: type === 'register' ? 'passenger' : ''
    });
    const [error, setError] = useState('');
    const [loading, setLoading] = useState(false);

    const handleChange = (e) => {
        setFormData({
            ...formData,
            [e.target.name]: e.target.value
        });
    };

    const handleSubmit = async (e) => {
        e.preventDefault();
        setError('');
        setLoading(true);

        try {
            await onSubmit(formData.username, formData.password, formData.role);
        } catch (err) {
            setError(err.message || 'An error occurred');
        } finally {
            setLoading(false);
        }
    };

    return (
        <Form onSubmit={handleSubmit}>
            {error && <Alert variant="danger">{error}</Alert>}

            <Form.Group className="mb-3">
                <Form.Label>Username</Form.Label>
                <Form.Control
                    type="text"
                    name="username"
                    value={formData.username}
                    onChange={handleChange}
                    placeholder="Enter username"
                    required
                />
            </Form.Group>

            <Form.Group className="mb-3">
                <Form.Label>Password</Form.Label>
                <Form.Control
                    type="password"
                    name="password"
                    value={formData.password}
                    onChange={handleChange}
                    placeholder="Enter password"
                    required
                />
            </Form.Group>

            {type === 'register' && (
                <Form.Group className="mb-3">
                    <Form.Label>Role</Form.Label>
                    <Form.Select
                        name="role"
                        value={formData.role}
                        onChange={handleChange}
                        required
                    >
                        <option value="passenger">Passenger</option>
                        <option value="driver">Driver</option>
                    </Form.Select>
                </Form.Group>
            )}

            <Button
                variant="primary"
                type="submit"
                className="w-100"
                disabled={loading}
            >
                {loading ? 'Processing...' : type === 'login' ? 'Login' : 'Register'}
            </Button>
        </Form>
    );
};

export default AuthForm;